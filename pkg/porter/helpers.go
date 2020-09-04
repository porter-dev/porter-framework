package porter

import (
	"os"
)

// TODO -- move these into an internal "utils" directory, together with
// other common helpers

// FileExists returns true if a file exists, false otherwise
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return !os.IsNotExist(err) && !info.IsDir()
}

// IsDirectory returns true if a path is a directory, false otherwise
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.IsDir()
}

// CLEANUP FUNCTION
