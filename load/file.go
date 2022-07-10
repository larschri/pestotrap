package load

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/itchyny/gojq"
)

type File struct {
	query    *gojq.Query
	filename string
}

type doc map[string]any

func (f *File) Docs() ([]doc, error) {
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

func (f *File) Key() string {
	s, _, _ := strings.Cut(path.Base(f.filename), ".")
	return s
}

func NewFile(fn string) (*File, error) {
	for _, p := range parsers {
		if strings.HasSuffix(fn, p.fileSuffix) {
			return &File{
				p.query,
				fn,
			}, nil
		}
	}

	return nil, fmt.Errorf("unsupported file type")
}
