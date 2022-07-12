package load

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/itchyny/gojq"
)

type file struct {
	query    *gojq.Query
	filename string
	index    bleve.Index
}

type doc map[string]any

func (f *file) docs() ([]doc, error) {
	bs, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return nil, err
	}

	var all any
	if err := json.Unmarshal(bs, &all); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	var result []doc
	iter := f.query.Run(all)
	for {
		d, ok := iter.Next()
		if !ok {
			break
		}

		m, ok := d.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("illformed document %v", d)
		}

		result = append(result, m)
	}

	return result, nil
}

func (f *file) Key() string {
	s, _, _ := strings.Cut(path.Base(f.filename), ".")
	return s
}

func (f *file) indexDocs() error {
	docs, err := f.docs()
	if err != nil {
		return err
	}

	batch := f.index.NewBatch()
	for _, a := range docs {
		batch.Index(fmt.Sprintf("%v", a[Field_ID]), a)
	}

	if err := f.index.Batch(batch); err != nil {
		return fmt.Errorf("failed to index: %w", err)
	}

	return nil
}

func NewFile(fn string) (*file, error) {
	for _, p := range parsers {
		if strings.HasSuffix(fn, p.fileSuffix) {
			return &file{
				p.query,
				fn,
				nil,
			}, nil
		}
	}

	return nil, fmt.Errorf("unsupported file type")
}
