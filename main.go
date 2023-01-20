package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/larschri/pestotrap/dirindex"
	"github.com/larschri/pestotrap/documents"
	"github.com/larschri/pestotrap/hxwrapper"
	"github.com/larschri/pestotrap/searchpage"
)

func main() {
	dir := flag.String("dir", "testdata/jsons", "directory with json files")
	index := flag.String("index", "myindex", "blevesearch index")
	addr := flag.String("addr", "localhost:8090", "the address for listening")
	flag.Parse()

	idx, err := dirindex.OpenIndex(*index)
	if err != nil {
		panic(err)
	}

	runner := dirindex.NewWatcher(os.DirFS(*dir), idx)
	go func() {
		if err := runner.Watch(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	searchHandler := searchpage.New(idx)

	r := http.NewServeMux()
	r.Handle("/s/", http.StripPrefix("/s", hxwrapper.Handler(searchHandler)))
	r.Handle("/x/", http.StripPrefix("/x", hxwrapper.Handler(documents.Server{os.DirFS(*dir)})))
	r.Handle("/", http.RedirectHandler("/s", http.StatusFound))
	log.Println("Starting server on ", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}
