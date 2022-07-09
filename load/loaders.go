package load

import (
	"fmt"
	"io/ioutil"
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

func File(fn string) (map[string]any, error) {
	parser := findParser(fn)
	if parser == nil {
		return nil, fmt.Errorf("unsupported file type")
	}

	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	docs, err := toDocs(parser, bs)
	if err != nil {
		return nil, err
	}

	s, _, _ := strings.Cut(path.Base(fn), ".")

	r := make(map[string]any)
	for _, a := range docs {
		r[fmt.Sprintf("%v.%v", s, a[Field_ID])] = a

	}
	return r, nil

}

func findParser(fn string) *gojq.Query {
	for _, p := range parsers {
		if strings.HasSuffix(fn, p.fileSuffix) {
			return p.query
		}
	}
	return nil
}
