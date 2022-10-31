package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/larschri/pestotrap/dirindex"
	"github.com/larschri/pestotrap/documents"
	"github.com/larschri/searchpage"
)

var config = searchpage.Config{

	Request: func(r *http.Request) *bleve.SearchRequest {
		b := searchpage.DefaultConfig.Request(r)
		b.Fields = []string{
			documents.Field_Name,
			documents.Field_Type,
			documents.Field_Taxonomy,
			documents.Field_Filename,
			documents.Field_ID,
		}
		return b
	},

	RenderMatches: func(w io.Writer, matches []*search.DocumentMatch) {
		for _, m := range matches {
			searchpage.DefaultMatch.Execute(w, map[string]interface{}{
				"Name": m.Fields[documents.Field_Name],
				"Type": m.Fields[documents.Field_Type],
				"Taxonomy": fmt.Sprintf("%v / %v",
					m.Fields[documents.Field_Filename],
					m.Fields[documents.Field_Taxonomy]),
				"Url": fmt.Sprintf("/x/%v?id=%v",
					m.Fields[documents.Field_Filename],
					m.Fields[documents.Field_ID]),
			})
		}
	},
}

func main() {
	dir := flag.String("dir", "", "directory with json files")
	index := flag.String("index", ".", "blevesearch index")
	flag.Parse()

	idx, err := dirindex.OpenIndex(*index)
	if err != nil {
		panic(err)
	}

	if err := dirindex.Start(os.DirFS(*dir), idx, time.NewTicker(10*time.Second).C); err != nil {
		panic(err)
	}

	searchHandler := searchpage.New(&config, idx)

	r := http.NewServeMux()
	r.Handle("/s/", http.StripPrefix("/s", searchHandler))
	r.Handle("/x/", http.StripPrefix("/x", documents.Server{os.DirFS(*dir)}))
	hostport := "localhost:8090"
	log.Println("Starting server on ", hostport)
	log.Fatal(http.ListenAndServe(hostport, r))
}
