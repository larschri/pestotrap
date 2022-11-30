package dirindex

import "github.com/blevesearch/bleve/v2"

func TestIndex() (bleve.Index, error) {
	return bleve.NewMemOnly(indexMapping)
}
