package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/larschri/pestotrap/load"
)

// findFiles list the files in the given directory
func findFiles(dir string) ([]string, error) {
	var r []string
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			r = append([]string(r), path)

			return nil
		})

	return r, err
}

// indexDirectory indexes all files in the given directory
func indexDirectory(dir string, index bleve.Index) error {

	fls, err := findFiles(dir)
	if err != nil {
		return fmt.Errorf("failed to travese %s: %w", dir, err)
	}

	flsMap := make(map[string]*load.File)
	for _, fn := range fls {
		f2, err := load.NewFile(fn)
		if err != nil {
			return err
		}

		if _, ok := flsMap[f2.Key()]; ok {
			return fmt.Errorf("duplicate file name/key: %v", f2.Key())
		}

		flsMap[f2.Key()] = f2
	}

	for _, fl := range flsMap {
		if err := fl.Index(index); err != nil {
			return fmt.Errorf("failed to load file %v: %w", fl, err)
		}

	}

	return nil
}
