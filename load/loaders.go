package load

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/itchyny/gojq"
)

type doc struct {
	Id       string
	Name     string
	Type     string
	Taxonomy string
	Filename string
	Doc      map[string]interface{}
}

var parsers = []struct {
	fileSuffix string
	query      *gojq.Query
}{
	{
		".k8s",
		jqMust(`.items|map({
				Id: .metadata.uid,
				Type: .kind,
				Name: .metadata.name,
				Taxonomy: .metadata.namespace,
				Doc: .
		})`),
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
		a.Doc["render"] = map[string]string{
			"name":     a.Name,
			"taxonomy": s + " / " + a.Taxonomy,
			"type":     a.Type,
		}
		r[a.Id] = a.Doc

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
