package main

import (
	"context"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/larschri/pestotrap/internal/demo"
	"github.com/larschri/pestotrap/internal/documents"
	"github.com/larschri/pestotrap/pkg/dirindex"
	"github.com/larschri/pestotrap/pkg/searchpage"
)

func main() {

	dir := flag.String("dir", "", "directory with json files")
	index := flag.String("index", "", "blevesearch index")
	addr := flag.String("addr", "localhost:8090", "the address for listening")
	flag.Parse()

	var dirFS = fs.FS(demo.Embed)

	if *dir != "" {
		dirFS = os.DirFS(*dir)
	}

	docserv, err := documents.NewServer(dirFS)
	if err != nil {
		panic(err)
	}

	idx, err := dirindex.OpenIndex(*index)
	if err != nil {
		panic(err)
	}

	runner := dirindex.NewWatcher(docserv.FS, idx)
	if *dir == "" {
		// dir is embedded, so don't start watcher
		runner.UpdateIfModified()
	} else {
		go func() {
			if err := runner.Watch(context.TODO(), *dir+"/data"); err != nil {
				log.Fatal(err)
			}
		}()
	}
	searchHandler := searchpage.New(idx)

	r := http.NewServeMux()
	r.Handle("/s/", http.StripPrefix("/s", searchHandler))
	r.Handle("/x/", http.StripPrefix("/x", docserv))
	r.Handle("/", http.RedirectHandler("/s/", http.StatusMovedPermanently))
	log.Println("Starting server on ", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}
