package load

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/itchyny/gojq"
)

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

// jqMust parses s into a gojq.Query or panics on failure
func jqMust(s string) *gojq.Query {
	q, err := gojq.Parse(s)
	if err != nil {
		panic(err)
	}
	return q
}
