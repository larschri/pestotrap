package load

import (
	"fmt"
	"io/ioutil"

	"github.com/blevesearch/bleve/v2"
)

// indexDirectory indexes all files in the given directory
func IndexDirectory(dir string, index bleve.Index) error {

	fls, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read dir %s: %w", dir, err)
	}

	flsMap := make(map[string]*File)
	for _, fn := range fls {
		if fn.IsDir() {
			continue
		}

		f2, err := NewFile(dir + "/" + fn.Name())
		if err != nil {
			return err
		}

		if _, ok := flsMap[f2.Key()]; ok {
			return fmt.Errorf("duplicate file name/key: %v", f2.Key())
		}

		flsMap[f2.Key()] = f2
	}

	for _, fl := range flsMap {
		docs, err := fl.Docs()
		if err != nil {
			return err
		}

		batch := index.NewBatch()
		for _, a := range docs {
			batch.Index(fmt.Sprintf("%v.%v", fl.Key(), a[Field_ID]), a)
		}

		if err := index.Batch(batch); err != nil {
			return fmt.Errorf("failed to index: %w", err)
		}

	}

	return nil
}
