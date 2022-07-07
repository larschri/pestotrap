package main

import (
	"fmt"
	"os"
	"path"
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

	batch := index.NewBatch()

	for _, fn := range fls {
		as, err := load.File(fn)
		if err != nil {
			return err
		}

		for k, v := range as {
			batch.Index(fmt.Sprintf("%s.%v", path.Base(fn), k), v)
		}
	}

	if err := index.Batch(batch); err != nil {
		return fmt.Errorf("failed to index: %w", err)
	}

	return nil
}
