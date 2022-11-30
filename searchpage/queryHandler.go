package searchpage

import (
	_ "embed"
	"fmt"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/larschri/pestotrap/documents"
)

const pageSize = 30

func request(r *http.Request) *bleve.SearchRequest {
	srch := ""
	if len(r.Form["search"]) > 0 {
		srch = r.Form["search"][0]
	}

	query := bleve.NewQueryStringQuery(srch)
	b := bleve.NewSearchRequest(query)

	b.Fields = documents.MatchFields

	b.Size = pageSize

	if len(r.Form["offset"]) > 0 {
		b.From, _ = strconv.Atoi(r.Form["offset"][0])
	}

	return b
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

	w.Header().Add("HX-Trigger", `{"setValidSearchInput":true}`)
	for _, m := range result.Hits {
		documents.RenderMatch(w, m.Fields)
	}

	nextOffset := uint64(request.From + pageSize)
	if nextOffset >= result.Total {
		return
	}

	q := r.URL.Query()
	q.Set("offset", strconv.Itoa(request.From+pageSize))
	fmt.Fprintf(w, `<div hx-get="q?%s" hx-trigger="revealed"/>`, q.Encode())

}
