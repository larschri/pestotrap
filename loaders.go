package main

import (
	"io/ioutil"
	"strings"

	"github.com/itchyny/gojq"
)

type doc struct {
	Id       string
	Name     string
	Type     string
	Taxonomy string
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

func loadK8s(filename string) ([]doc, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var xx []doc
	if err := jq(k8sJq, bs, &xx); err != nil {
		return nil, err
	}

	return xx, nil
}

func load(fn string) ([]doc, error) {
	if strings.HasSuffix(fn, ".k8s") {
		return loadK8s(fn)
	}

	panic("not supported")
	return nil, nil
}
