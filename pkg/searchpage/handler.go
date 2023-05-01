package searchpage

import (
	_ "embed"
	"net/http"

	"github.com/blevesearch/bleve/v2"
	bhttp "github.com/blevesearch/bleve/v2/http"
	"github.com/larschri/pestotrap/web/templates"
)

type Handler struct {
	indices map[string]bleve.Index
	alias   bleve.IndexAlias
}

func (h *Handler) indexHTMLHandler(w http.ResponseWriter, r *http.Request) {
	templates.Templates.ExecuteTemplate(w, "index.tmpl", nil)
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
