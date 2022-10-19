package dirindex

import (
	"strings"

	"github.com/itchyny/gojq"
)

const (
	// Field_ID identifies an object inside a file, combined with
	// Field_Filename to construct the document id
	Field_ID = "xxid"

	// Field_Type is the type/kind for the document. Used for display only
	Field_Type = "xxtype"

	// Field_Name is a human readlable name for the document. Used for
	// display only
	Field_Name = "xxname"

	// Field_Taxonomy is a human readable taxonomy. Used for display only.
	Field_Taxonomy = "xxtaxonomy"

	// Field_Filename is the basename of the source file. Used to construct
	// the document id
	Field_Filename = "xxfilename"

	// Field_FileVersion is a keyword field that identifies versioned
	// source file. It must change when the source file changes.
	Field_FileVersion = "xxfileversion"
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
	{
		".raw",
		jqMust(".[]"),
	},
}

// jqMust parses s into a gojq.Query or panics on failure
func jqMust(s string) *gojq.Query {
	q, err := gojq.Parse(s)
	if err != nil {
		panic(err)
	}
	return q
}

func parser(name string) *gojq.Query {
	for _, p := range parsers {
		if strings.HasSuffix(name, p.fileSuffix) {
			return p.query
		}
	}

	return nil
}
