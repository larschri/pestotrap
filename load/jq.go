package load

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

type doc map[string]any

func toDocs(q *gojq.Query, bs []byte) ([]doc, error) {
	var all any
	if err := json.Unmarshal(bs, &all); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	var result []doc
	iter := q.Run(all)
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
