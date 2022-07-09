package load

import (
	"fmt"
	"path"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/itchyny/gojq"
)

const (
	Field_ID       = "xxid"
	Field_Type     = "xxtype"
	Field_Name     = "xxname"
	Field_Taxonomy = "xxtaxonomy"
)

var parsers = []struct {
	fileSuffix string
	query      *gojq.Query
}{
	{
		".k8s",
		jqMust(`.items[] | . +
		{
			xxid: .metadata.uid,
			xxtype: .kind,
			xxname: .metadata.name,
			xxtaxonomy: .metadata.namespace
		}`),
	},
}

type File struct {
	query    *gojq.Query
	filename string
}

func (f *File) Key() string {
	s, _, _ := strings.Cut(path.Base(f.filename), ".")
	return s
}

func (f *File) Index(index bleve.Index) error {

	docs, err := f.Docs()
	if err != nil {
		return err
	}

	batch := index.NewBatch()
	for _, a := range docs {
		batch.Index(fmt.Sprintf("%v.%v", f.Key(), a[Field_ID]), a)
	}

	if err := index.Batch(batch); err != nil {
		return fmt.Errorf("failed to index: %w", err)
	}

	return nil
}

func NewFile(fn string) (*File, error) {
	for _, p := range parsers {
		if strings.HasSuffix(fn, p.fileSuffix) {
			return &File{
				p.query,
				fn,
			}, nil
		}
	}

	return nil, fmt.Errorf("unsupported file type")
}
