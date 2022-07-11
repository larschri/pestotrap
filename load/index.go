package load

import (
	"fmt"
	"io/ioutil"

	"github.com/blevesearch/bleve/v2"
)

// indexDirectory indexes all files in the given directory
func IndexDirectory(dir string, idxDir string) ([]bleve.Index, error) {

	fls, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %w", dir, err)
	}

	flsMap := make(map[string]*file)
	for _, fn := range fls {
		if fn.IsDir() {
			continue
		}

		f2, err := NewFile(dir + "/" + fn.Name())
		if err != nil {
			return nil, err
		}

		if _, ok := flsMap[f2.Key()]; ok {
			return nil, fmt.Errorf("duplicate file name/key: %v", f2.Key())
		}

		flsMap[f2.Key()] = f2
	}

	mapping := bleve.NewIndexMapping()
	var indices []bleve.Index
	for _, fl := range flsMap {
		index, err := bleve.New(idxDir+"/"+fl.Key(), mapping)
		if err != nil {
			index, err = bleve.Open(idxDir + "/" + fl.Key())
			if err != nil {
				return nil, err
			}
		}

		docs, err := fl.docs()
		if err != nil {
			return nil, err
		}

		batch := index.NewBatch()
		for _, a := range docs {
			batch.Index(fmt.Sprintf("%v", a[Field_ID]), a)
		}

		if err := index.Batch(batch); err != nil {
			return nil, fmt.Errorf("failed to index: %w", err)
		}

		indices = append(indices, index)
	}

	return indices, nil
}
