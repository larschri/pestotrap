package searchpage

import (
	_ "embed"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/larschri/pestotrap/pkg/filetypes"
	"github.com/larschri/pestotrap/web/templates"
)

// pageSize is the number of hits to return in one search request. It is used
// in combination with the "offset" form value.
const pageSize = 20

// request converts a *http.Request into a *bleve.SearchRequest
func request(r *http.Request) *bleve.SearchRequest {
	srch := ""
	if len(r.Form["search"]) > 0 {
		srch = r.Form["search"][0]
	}

	query := bleve.NewQueryStringQuery(srch)
	b := bleve.NewSearchRequest(query)

	b.Fields = filetypes.MatchFields

	b.Size = pageSize

	if len(r.Form["offset"]) > 0 {
		b.From, _ = strconv.Atoi(r.Form["offset"][0])
	}

	return b
}

type RenderContext struct {
	Matches template.HTML
	Next    string
}

func (h *Handler) searchQueryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Add("HX-Trigger", `{"setValidSearchInput":false}`)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request := request(r)

	result, err := h.alias.Search(request)
	if err != nil {
		w.Header().Add("HX-Trigger", `{"setValidSearchInput":false}`)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var buf strings.Builder
	for _, r := range result.Hits {
		if err := filetypes.Render(&buf, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("HX-Trigger", `{"setValidSearchInput":true}`)

	var next string
	if uint64(request.From+pageSize) < result.Total {
		q := r.URL.Query()
		q.Set("offset", strconv.Itoa(request.From+pageSize))
		next = q.Encode()
	}

	templates.Templates.ExecuteTemplate(w, "matches.tmpl", &RenderContext{
		template.HTML(buf.String()),
		next,
	})
}
