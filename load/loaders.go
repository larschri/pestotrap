package load

import (
	"fmt"
	"path"
	"strings"

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
