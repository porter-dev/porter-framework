package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/porterdev/ego/pkg/porter"
	t "github.com/porterdev/ego/pkg/translator"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a path to a file\n ")
		}

		err := porter.FileExists(args[0])

		if err != nil {
			return err
		}

		return nil
	},
	Use:   "apply [filename]",
	Short: "Applies a .gop file -- generates configuration and runs it.",
	Long:  `TBD -- not sure what this will accept yet.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		apply(filename)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func apply(filename string) {
	// TODO -- from the entry point, find all .gop files in the same directory or below

	// read the entry point
	dat, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Error while reading file:", err)
	}

	// translate the file to Golang
	trans := t.NewTranslator(dat)

	res := trans.TranslateToJSON()

	fmt.Println(string(res))

	dir, err := ioutil.TempDir("./", "build")

	if err != nil {
		fmt.Println(err)
		panic("ouch")
	}

	defer os.RemoveAll(dir)

	// TODO -- SUPPORT OTHER FILES THAN MAIN
	file, err := ioutil.TempFile("./"+dir, "build_*.go")

	if err != nil {
		fmt.Println(err)
		panic("ouch")
	}

	file.Write(res)

	fmt.Println(file.Name())

	// compile the file
	cmd := exec.Command("go", "build", "./"+file.Name())

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
