package main

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

func jqOne(q *gojq.Query, bs []byte, target any) error {
	a, err := jqraw(q, bs)
	if err != nil {
		return err
	}

	if len(a) != 1 {
		return fmt.Errorf("expected on json object got %d", len(a))
	}

	return unmarshal(a[0], target)
}

// jq unmarshals the given bytes into target
func jq(q *gojq.Query, bs []byte, target any) error {
	a, err := jqraw(q, bs)
	if err != nil {
		return err
	}

	return unmarshal(a, target)
}

// jqraw unmarshals the given bytes into an untyped slice
func jqraw(q *gojq.Query, bs []byte) ([]any, error) {
	var all any
	if err := json.Unmarshal(bs, &all); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	var result []any

	iter := q.Run(all)
	for {
		x, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := x.(error); ok {
			return nil, fmt.Errorf("jq error: %w", err)
		}

		result = append(result, x)
	}

	return result, nil
}

// unmarshal is a helper function to translate a generic all object into a type
func unmarshal(src, target any) error {
	bs, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	if err := json.Unmarshal(bs, target); err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	return nil
}

// jqMust parses s into a gojq.Query or panics on failure
func jqMust(s string) *gojq.Query {
	q, err := gojq.Parse(s)
	if err != nil {
		panic(err)
	}
	return q
}
