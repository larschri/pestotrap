package load

import (
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

// jqMust parses s into a gojq.Query or panics on failure
func jqMust(s string) *gojq.Query {
	q, err := gojq.Parse(s)
	if err != nil {
		panic(err)
	}
	return q
}
