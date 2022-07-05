package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
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

func fileToBatch(fn string, b *bleve.Batch, s string) error {
	as, err := load(fn)
	if err != nil {
		return err
	}

	for _, a := range as {
		b.Index(fmt.Sprintf("%s.%v", s, a["id"]), a)
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
		t := strings.TrimPrefix(e, dir+"/")
		if err := fileToBatch(e, batch, t); err != nil {
			return err
		}
	}

	if err := index.Batch(batch); err != nil {
		return fmt.Errorf("failed to index: %w", err)
	}

	return nil
}
