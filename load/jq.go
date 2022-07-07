package load

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

func toDocs(q *gojq.Query, bs []byte) ([]doc, error) {
	var all any
	if err := json.Unmarshal(bs, &all); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	iter := q.Run(all)
	d, ok := iter.Next()
	if !ok {
		return nil, fmt.Errorf("failed get output from jq")
	}

	if _, ok := iter.Next(); ok {
		return nil, fmt.Errorf("unexpected output from jq")
	}

	bs2, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("failed to re-marshal object from jq: %w", err)
	}

	var docs []doc
	if err := json.Unmarshal(bs2, &docs); err != nil {
		return nil, fmt.Errorf("failed to re-unmarshal to docs: %w", err)
	}

	return docs, nil
}

// jqMust parses s into a gojq.Query or panics on failure
func jqMust(s string) *gojq.Query {
	q, err := gojq.Parse(s)
	if err != nil {
		panic(err)
	}
	return q
}
