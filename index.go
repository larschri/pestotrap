package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/larschri/pestotrap/load"
)

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

func fileToBatch(fn string, b *bleve.Batch) error {
	as, err := load.File(fn)
	if err != nil {
		return err
	}

	for _, a := range as {
		a.Doc["render"] = map[string]string{
			"name":     a.Name,
			"taxonomy": path.Base(fn) + " / " + a.Taxonomy,
			"type":     a.Type,
		}
		b.Index(fmt.Sprintf("%s.%v", path.Base(fn), a.Id), a.Doc)
	}
	return nil
}

func indexDirectory(dir string, index bleve.Index) error {

	fls, err := findFiles(dir)
	if err != nil {
		return fmt.Errorf("failed to travese %s: %w", dir, err)
	}

	batch := index.NewBatch()

	for _, e := range []string(fls) {
		if err := fileToBatch(e, batch); err != nil {
			return err
		}
	}

	if err := index.Batch(batch); err != nil {
		return fmt.Errorf("failed to index: %w", err)
	}

	return nil
}
