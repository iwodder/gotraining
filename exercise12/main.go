package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	var (
		directory  = flag.String("dir", "", "The directory to search for files.")
		curPattern = flag.String("cur", "", "The current filename format.")
		newPattern = flag.String("new", "", "The current new filename format.")
	)
	flag.Parse()

	if *directory == "" || *curPattern == "" || *newPattern == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	absDir, err := filepath.Abs(*directory)
	if err != nil {
		panic(err)
	}

	rex, err := regexp.Compile(*curPattern)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Running rename in %s, searching for %s, renaming to %s", absDir, *curPattern, *newPattern)

	toRename := make([]string, 0, 10)
	err = filepath.WalkDir(absDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			baseName := filepath.Base(path)
			if rex.MatchString(baseName) {
				toRename = append(toRename, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	renameFiles(toRename, *newPattern)

	fmt.Printf("Done!")
}

func renameFiles(files []string, newPattern string) {
	cnt := len(files)
	for i, v := range files {
		_ = os.Rename(v, filepath.Join(filepath.Dir(v), fmt.Sprintf(newPattern, i+1, cnt)))
	}
}
