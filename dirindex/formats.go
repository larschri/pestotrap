package dirindex

import (
	"strings"

	"github.com/itchyny/gojq"
)

const (
	Field_ID          = "xxid"
	Field_Type        = "xxtype"
	Field_Name        = "xxname"
	Field_Taxonomy    = "xxtaxonomy"
	Field_Filename    = "xxfilename"
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
