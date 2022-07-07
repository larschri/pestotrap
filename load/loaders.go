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
	docs, err := fileToDocs(fn)
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

func fileToDocs(fn string) ([]doc, error) {
	for _, p := range parsers {
		if strings.HasSuffix(fn, p.fileSuffix) {
			bs, err := ioutil.ReadFile(fn)
			if err != nil {
				return nil, err
			}

			return toDocs(p.query, bs)
		}
	}

	return nil, fmt.Errorf("unsupported file type")
}
