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
	"github.com/larschri/searchpage"
)

var config = searchpage.Config{

	Request: func(r *http.Request) *bleve.SearchRequest {
		b := searchpage.DefaultConfig.Request(r)
		b.Fields = []string{
			dirindex.Field_Name,
			dirindex.Field_Type,
			dirindex.Field_Taxonomy,
			dirindex.Field_Filename,
		}
		return b
	},

	RenderMatches: func(w io.Writer, matches []*search.DocumentMatch) {
		for _, m := range matches {
			searchpage.DefaultMatch.Execute(w, map[string]interface{}{
				"Name": m.Fields[dirindex.Field_Name],
				"Type": m.Fields[dirindex.Field_Type],
				"Taxonomy": fmt.Sprintf("%v / %v",
					m.Fields[dirindex.Field_Filename],
					m.Fields[dirindex.Field_Taxonomy]),
				"Url": "d?index=" + m.Index + "&doc=" + m.ID,
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
	hostport := "localhost:8090"
	log.Println("Starting server on ", hostport)
	log.Fatal(http.ListenAndServe(hostport, r))
}
