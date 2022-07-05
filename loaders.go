package main

import (
	"io/ioutil"
	"strings"

	"github.com/itchyny/gojq"
)

var k8sJq *gojq.Query = jqMust(`.items[]
	| . += {
		id: .metadata.uid,
		render:{
			type: .kind,
			name: .metadata.name,
			taxonomy: .metadata.namespace
		}
	}`)

func loadK8s(filename string) ([]map[string]any, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var xx []map[string]any
	if err := jq(k8sJq, bs, &xx); err != nil {
		return nil, err
	}

	return xx, nil
}

func load(fn string) ([]map[string]any, error) {
	if strings.HasSuffix(fn, ".k8s") {
		return loadK8s(fn)
	}

	panic("not supported")
	return nil, nil
}
