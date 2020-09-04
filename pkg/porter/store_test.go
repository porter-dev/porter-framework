package porter

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/porterdev/ego/pkg/json"

	v "github.com/porterdev/ego/internal/value"
)

func TestStoreSingleWrite(t *testing.T) {
	logger := NewLogger(2)

	val := v.Object{
		v.String("hello"):   v.String("there"),
		v.String("general"): v.String("kenobi"),
	}

	store, _ := NewLocalStore("12345", logger, "./", "./")

	store.WriteState(val)
	defer os.RemoveAll("./state_12345.json")

	// read the contents of the file and convert to object
	dat, _ := ioutil.ReadFile("./state_12345.json")
	res, _ := json.Inject(string(dat))

	if !v.IsEqual(res, val) {
		res1, _ := json.ToJSON(res)
		res2, _ := json.ToJSON(val)
		t.Errorf("Failed on simple write, %s, %s", res1, res2)
	}
}

func TestStoreMultiWrite(t *testing.T) {
	logger := NewLogger(2)

	val1 := v.Object{
		v.String("hello"):   v.String("there"),
		v.String("general"): v.String("kenobi"),
	}

	val2 := v.Object{
		v.String("hello2"):   v.String("there2"),
		v.String("general2"): v.String("kenobi2"),
	}

	store, _ := NewLocalStore("12345", logger, "./", "./")

	store.WriteState(val1)
	store.WriteState(val2)
	defer os.RemoveAll("./state_12345.json")

	// read the backups
	backups, _ := store.GetAllBackups()

	fmt.Println(backups)

	// store.WriteState(val2)
	// defer os.RemoveAll("./state_12345.json")

	// read the contents of the file and convert to object
	// dat, _ := ioutil.ReadFile("./state_12345.json")
	// res, _ := json.Inject(string(dat))

	// if !v.IsEqual(res, val) {
	// 	res1, _ := json.ToJSON(res)
	// 	res2, _ := json.ToJSON(val)
	// 	t.Errorf("Failed on simple write, %s, %s", res1, res2)
	// }
}

// func TestStoreMultipleWrite(t *testing.T) {
// 	for _, c := range writeSingleTests {
// 		logger := NewLogger(2)

// 		store := LocalStore{
// 			ID:     c.id,
// 			Path:   c.dirpath,
// 			Logger: logger,
// 		}

// 		store.WriteState(c.value)
// 		defer os.RemoveAll(c.want.filepath)

// 		// read the contents of the file and convert to object
// 		dat, _ := ioutil.ReadFile(c.want.filepath)
// 		res, _ := json.Inject(string(dat))

// 		if !v.IsEqual(res, c.want.res) {
// 			res1, _ := json.ToJSON(res)
// 			res2, _ := json.ToJSON(c.want.res)
// 			t.Errorf("Failed on: %s, %s, %s", c.name, res1, res2)
// 		}
// 	}
// }
