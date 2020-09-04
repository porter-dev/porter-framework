package porter

// TODO -- place a lock on the file for safety -- all Stores should implement
// the locking interface somehow. This needs some work and guidance from similarly
// motivated locking mechanisms (Terraform, .kube/config, etc)

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/porterdev/ego/pkg/json"
)

// Store implements methods to read and write from a state store. See
// LocalStore for an implementation.
type Store interface {
	GetState() (Object, error)
	GetAllBackups() ([]string, error)
	GetBackup(filename string) (Object, error)
	WriteState(v Object) error
	WriteBackup(filename string) error
}

// LocalStore is an implementation of a store that uses the local filesystem
// to store state and backups.
type LocalStore struct {
	ID string

	Logger *Logger

	StateDir  string
	BackupDir string

	Lock Lock
}

// NewLocalStore initializes a local store
func NewLocalStore(id string, logger *Logger, stateDir, backupDir string) (*LocalStore, error) {
	// verify that path is accessible and open for reading
	if !IsDirectory(stateDir) {
		return nil, errors.New("State directory is not a directory")
	}

	if !IsDirectory(backupDir) {
		return nil, errors.New("Backup directory is not a directory")
	}

	return &LocalStore{
		ID:        id,
		Logger:    logger,
		StateDir:  stateDir,
		BackupDir: backupDir,
	}, nil
}

// GetState returns the current state as a Porter object.
func (s *LocalStore) GetState() (Object, error) {
	filename := filepath.Join(s.StateDir, "state_"+s.ID+".json")

	dat, err := ioutil.ReadFile(filename)
	s.Logger.Check(err, s.ID, "error reading state file", filename)

	res, err := json.Inject(string(dat))
	s.Logger.Check(err, s.ID, "error converting state file to json", filename)

	return res, nil
}

// GetAllBackups returns an array of filenames, sorted from most recent to
// least recent, of all backup timestamps in a directory and all subdirectories.
//
// We query subdirectories so that backups stored in folders such as /archive and
// /backups can be added. This places some constraints on the recommended folder
// organization.
func (s *LocalStore) GetAllBackups() ([]string, error) {
	// add complete paths when encountered
	files := []string{}
	pattern := "backup_" + s.ID + `_\d*\.json`

	fileinfos, err := ioutil.ReadDir(s.BackupDir)
	s.Logger.Check(err, s.ID, "error reading directory", s.BackupDir)

	for _, fi := range fileinfos {
		fn := fi.Name()
		match, err := regexp.Match(pattern, []byte(fn))

		if err == nil && match && !fi.IsDir() {
			files = append(files, fn)
		}
	}

	// TODO -- SORT THE BACKUPS

	return files, nil
}

// GetBackup returns the a backup based on a timestamp in the form of a
// string.
func (s *LocalStore) GetBackup(filename string) (Object, error) {
	dat, err := ioutil.ReadFile(filename)
	s.Logger.Check(err, s.ID, "error reading backup file", filename)

	res, err := json.Inject(string(dat))
	s.Logger.Check(err, s.ID, "error converting backup file to json", filename)

	return res, nil
}

// WriteState saves a Porter object to the filesystem as JSON.
func (s *LocalStore) WriteState(v Object) error {
	// convert value to JSON string
	str, err := json.ToJSON(v)

	s.Logger.Check(err, s.ID, "error converting to JSON")

	// store in file
	filename := filepath.Join(s.StateDir, "state_"+s.ID+".json")

	// check if file already exists -- if it does, store as backup
	exists := FileExists(filename)

	if exists {
		err = s.WriteBackup(filename)
		s.Logger.Check(err, s.ID, "error writing backup file")
	}

	err = ioutil.WriteFile(filename, []byte(str), 0644)
	s.Logger.Check(err, s.ID, "error saving state file")

	return nil
}

// WriteBackup takes in a state file and writes a backup, using the ID
// stored in the LocalStore object. It then removes the existing state
// file.
func (s *LocalStore) WriteBackup(filename string) error {
	src, err := os.Open(filename)
	s.Logger.Check(err, s.ID, "error reading file")
	defer src.Close()

	dir := filepath.Dir(filename)

	// if s.StateDir is not same directory as filename, throw an error
	if rel, err := filepath.Rel(dir, s.StateDir); err != nil || rel != "." {
		return errors.New("LocalStore Path does not match filename directory")
	}

	ts := strconv.Itoa(int(time.Now().Unix()))
	backup := filepath.Join(s.BackupDir, "backup_"+s.ID+"_"+ts+".json")

	dest, err := os.OpenFile(backup, os.O_RDWR|os.O_CREATE, 0666)
	s.Logger.Check(err, s.ID, "error opening backup file")
	defer dest.Close()

	_, err = io.Copy(dest, src)
	s.Logger.Check(err, s.ID, "error copying from state file to backup file")

	// if we've reached here, remove state file
	err = os.Remove(filename)
	s.Logger.Check(err, s.ID, "error removing existing state file")

	return nil
}
