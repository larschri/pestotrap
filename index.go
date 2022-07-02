package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func decodeAndBatch(r io.Reader, b *bleve.Batch, s string) error {
	decoder := json.NewDecoder(r)
	i := 0
	for {
		var doc map[string]interface{}
		if err := decoder.Decode(&doc); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode json: %w", err)
		}

		if _, ok := doc["_type"]; !ok {
			doc["_type"] = doc["kind"]
		}
		doc["Taxonomy"] = s

		b.Index(fmt.Sprintf("%s.%v", s, i), doc)
		i++
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
		f, err := os.Open(e)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", e, err)
		}
		defer f.Close()

		t := strings.TrimPrefix(e, dir+"/")

		if err := decodeAndBatch(f, batch, t); err != nil {
			return fmt.Errorf("failed to process %s: %w", e, err)
		}
	}

	if err := index.Batch(batch); err != nil {
		return fmt.Errorf("failed to index: %w", err)
	}

	return nil
}
