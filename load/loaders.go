package load

import (
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

var k8sJq *gojq.Query = jqMust(`.items[]
	| {
		Id: .metadata.uid,
		Type: .kind,
		Name: .metadata.name,
		Taxonomy: .metadata.namespace,
		Doc: .
	}`)

func loadDocs(filename string, query *gojq.Query) ([]doc, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var docs []doc
	if err := jq(query, bs, &docs); err != nil {
		return nil, err
	}

	return docs, nil
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
			"taxonomy": s + "/" + a.Taxonomy,
			"type":     a.Type,
			"filename": fn,
		}
		r[a.Id] = a.Doc

	}
	return r, nil

}

func fileToDocs(fn string) ([]doc, error) {
	if strings.HasSuffix(fn, ".k8s") {
		return loadDocs(fn, k8sJq)
	}

	panic("not supported")
	return nil, nil
}
