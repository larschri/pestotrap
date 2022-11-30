package searchpage

import (
	_ "embed"
	"net/http"

	"github.com/blevesearch/bleve/v2"
	bhttp "github.com/blevesearch/bleve/v2/http"
)

//go:embed form.htmx
var searchForm []byte

type Handler struct {
	indices map[string]bleve.Index
	alias   bleve.IndexAlias
}

func (h *Handler) indexHTMLHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(searchForm)
}

func New(indices ...bleve.Index) http.Handler {
	m := make(map[string]bleve.Index)

	for _, ix := range indices {
		m[ix.Name()] = ix
		bhttp.RegisterIndexName(ix.Name(), ix)
	}

	h := Handler{
		m,
		bleve.NewIndexAlias(indices...),
	}

	r := http.NewServeMux()
	r.HandleFunc("/", h.indexHTMLHandler)
	r.HandleFunc("/q", h.searchQueryHandler)
	return r
}
