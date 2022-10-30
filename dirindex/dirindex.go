package dirindex

import (
	"fmt"
	"io/fs"
	"log"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/larschri/pestotrap/documents"
)

type doc map[string]any

// OpenIndex opens or creates the index
func OpenIndex(dir string) (bleve.Index, error) {
	idx, err := bleve.Open(dir)
	if err == nil {
		return idx, nil
	}

	return bleve.New(dir, Mapping())
}

func Mapping() mapping.IndexMapping {
	m := bleve.NewIndexMapping()
	m.DefaultMapping.AddFieldMappingsAt(documents.Field_FileVersion,
		bleve.NewKeywordFieldMapping())
	return m
}

// fileModTimes helper function to extract versions from index
func fileModTimes(index bleve.Index) (map[string]bool, error) {
	q := bleve.NewMatchAllQuery()
	fr := bleve.NewFacetRequest(documents.Field_FileVersion, 1000000)
	r := bleve.NewSearchRequest(q)
	r.AddFacet("f", fr)
	r.Size = 0
	rs, err := index.Search(r)
	if err != nil {
		return nil, err
	}

	m := make(map[string]bool)
	for _, f := range rs.Facets["f"].Terms.Terms() {
		m[f.Term] = true
	}
	return m, nil
}

// indexFile helper function to iterate over the documents in a file
func indexFile(index bleve.Index, dir fs.FS, e fs.DirEntry) error {
	docs, err := documents.Documents(dir, e.Name())
	if err != nil {
		return err
	}

	batch := index.NewBatch()
	for _, doc := range docs {
		doc[documents.Field_FileVersion] = filestate(e)
		doc[documents.Field_Filename] = e.Name()
		batch.Index(fmt.Sprintf("%s/%s", doc[documents.Field_Filename], doc[documents.Field_ID]), doc)
	}
	return index.Batch(batch)
}

// deleteDocsByFileModTime deletes old docs
func deleteDocsByFileModTime(idx bleve.Index, mtm string) error {
	q := bleve.NewTermQuery(mtm)
	q.FieldVal = documents.Field_FileVersion

	res, err := idx.Search(bleve.NewSearchRequest(q))
	if err != nil {
		return err
	}

	batch := idx.NewBatch()
	for _, r := range res.Hits {
		batch.Delete(r.ID)
	}

	return idx.Batch(batch)
}

// update updates the index with the contents of the given dir
func update(index bleve.Index, dir fs.FS) error {
	if index.Mapping().AnalyzerNameForPath(documents.Field_FileVersion) != "keyword" {
		return fmt.Errorf("incompatible index")
	}

	modTimes, err := fileModTimes(index)
	if err != nil {
		return err
	}

	files, err := fs.ReadDir(dir, ".")
	if err != nil {
		return err
	}

	for _, f := range files {
		if _, ok := modTimes[filestate(f)]; ok {
			delete(modTimes, filestate(f))
			continue
		}

		if err := indexFile(index, dir, f); err != nil {
			log.Printf("skipping %s: %s", f.Name(), err.Error())
		}
	}

	for k, _ := range modTimes {
		if err := deleteDocsByFileModTime(index, k); err != nil {
			log.Printf("delete failed %s: %s", k, err.Error())
			return err
		}
	}

	return nil
}

func updateIfModified(index bleve.Index, dir fs.FS, state dirstate) dirstate {
	newState, err := newDirstate(dir)
	if err != nil {
		return state
	}

	if state.equal(newState) {
		return state
	}

	if err := update(index, dir); err != nil {
		return state
	}

	return newState
}

// Start the indexing goroutine that checks for updates
func Start(dir fs.FS, index bleve.Index, c <-chan time.Time) error {
	state, err := newDirstate(dir)
	if err != nil {
		return err
	}

	if err := update(index, dir); err != nil {
		return err
	}

	go func() {
		for {
			_, ok := <-c
			if !ok {
				return
			}
			state = updateIfModified(index, dir, state)
		}
	}()
	return nil
}
