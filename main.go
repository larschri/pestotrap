package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/gorilla/mux"
	"github.com/larschri/pestotrap/load"
	"github.com/larschri/searchpage"
)

var config = searchpage.Config{

	Request: func(r *http.Request) *bleve.SearchRequest {
		b := searchpage.DefaultConfig.Request(r)
		b.Fields = []string{
			load.Field_Name,
			load.Field_Type,
			load.Field_Taxonomy,
		}
		return b
	},

	RenderMatches: func(w io.Writer, matches []*search.DocumentMatch) {
		for _, m := range matches {
			f, _, _ := strings.Cut(m.ID, ".")
			searchpage.DefaultMatch.Execute(w, map[string]interface{}{
				"Name": m.Fields[load.Field_Name],
				"Type": m.Fields[load.Field_Type],
				"Taxonomy": fmt.Sprintf("%v / %v",
					f,
					m.Fields[load.Field_Taxonomy]),
				"Url": "d/" + m.Index + "/" + m.ID,
			})
		}
	},
}

func main() {
	dir := flag.String("dir", "", "directory with json files")
	index := flag.String("index", ".", "blevesearch index")
	flag.Parse()

	mapping := bleve.NewIndexMapping()
	idx, err := bleve.New(*index, mapping)
	if err != nil {
		idx, err = bleve.Open(*index)
		if err != nil {
			panic(err)
		}
	}

	if *dir != "" {
		if err := load.IndexDirectory(*dir, idx); err != nil {
			panic(err)
		}
	}

	searchHandler := searchpage.New(&config, idx)

	r := mux.NewRouter()
	r.PathPrefix("/s").Handler(http.StripPrefix("/s", searchHandler))
	hostport := "localhost:8090"
	log.Println("Starting server on ", hostport)
	log.Fatal(http.ListenAndServe(hostport, r))
}
