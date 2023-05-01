package documents

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"

	"github.com/blevesearch/bleve/v2/search"
	"github.com/itchyny/gojq"
	"github.com/larschri/pestotrap/pkg/filetypes"
)

const (
	// Field_ID identifies an object inside a file, combined with
	// Field_Filename to construct the document id
	field_ID = "id"

	// field_Title is the title for the document. Used for
	// display only
	field_Title = "title"

	// field_Description is a description of the document. Used for
	// display only
	field_Description = "description"
)

var (
	//go:embed match.tmpl
	m         string
	matchTmpl = template.Must(template.New("t").Parse(m))
)

type (
	// doc implements the filetypes.doc interface
	doc map[string]any

	// filetype implements the filetype.FileType interface
	filetype struct {
		query *gojq.Query
	}
)

func registerFiletype(t, s string) error {
	q, err := gojq.Parse(s)
	if err != nil {
		return err
	}

	filetypes.Register(t, filetype{q})
	return nil
}

func (d doc) ID() string {
	return d[field_ID].(string)
}

func (d doc) Set(key, val string) {
	d[key] = val
}

func (ft filetype) MatchFields() []string {
	return matchFields
}

func (ft filetype) Render(w io.Writer, d *search.DocumentMatch) error {
	return matchTmpl.Execute(w, (*Match)(d))
}

func (ft filetype) Documents(bs []byte) ([]filetypes.Doc, error) {

	var all any
	if err := json.Unmarshal(bs, &all); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	var docs []filetypes.Doc

	iter := ft.query.Run(all)

	for {
		d, ok := iter.Next()
		if !ok {
			break
		}

		m, ok := d.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("illformed document %v", d)
		}

		docs = append(docs, doc(m))
	}

	return docs, nil
}
