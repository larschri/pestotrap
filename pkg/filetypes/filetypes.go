package filetypes

import (
	"fmt"
	"io"
	"io/fs"
	"path"

	"github.com/blevesearch/bleve/v2/search"
)

const (
	Field_Prefix = "xx"

	// Field_Filename is the basename of the source file. Used to construct
	// the document id
	Field_Filename = Field_Prefix + "filename"

	// Field_FileVersion is a keyword field that identifies versioned
	// source file. It must change when the source file changes.
	Field_FileVersion = Field_Prefix + "fileversion"
)

var (
	filetypes = make(map[string]FileType)

	MatchFields []string
)

// FileType enables client code to process a type of files. FileTypes should be
// registered in the Register function.
type FileType interface {
	Documents(bs []byte) ([]Doc, error)
	MatchFields() []string
	Render(io.Writer, *search.DocumentMatch) error
}

// Doc is a document object returned from FileType
type Doc interface {
	// ID identifies the object
	ID() string

	// Set a key/value in the indexed document. Used to set source filename
	// and version.
	Set(key, val string)
}

func Register(suffix string, t FileType) {
	filetypes[suffix] = t

outer:
	for _, f1 := range t.MatchFields() {
		for _, f0 := range MatchFields {
			if f1 == f0 {
				continue outer
			}
		}
		MatchFields = append(MatchFields, f1)
	}
}

// Documents returns the documents from the given file
func Documents(dir fs.FS, fname string) ([]Doc, error) {
	t, ok := filetypes[path.Ext(fname)]
	if !ok {
		return nil, fmt.Errorf("no type for %s", fname)
	}

	bs, err := fs.ReadFile(dir, fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", fname, err)
	}

	return t.Documents(bs)
}

func Render(w io.Writer, d *search.DocumentMatch) error {
	fname, ok := d.Fields[Field_Filename].(string)
	if !ok {
		return fmt.Errorf("missing filename for document %s", d.ID)
	}

	t, ok := filetypes[path.Ext(fname)]
	if !ok {
		return fmt.Errorf("no type for %s", fname)
	}

	return t.Render(w, d)
}
