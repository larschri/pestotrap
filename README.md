# pestotrap

A search application for searching objects inside json files.

Run `go run ./cmd/pestotrap` and search for jokes on http<span>://</span>localhost:8090

## Configuration

[/internal/demo](/internal/demo) contains some example data and configuration. The data type _joke_ is defined by how the files are named.
* [/internal/demo/data/foo.joke](/internal/demo/data/foo.joke) is a file containing many _joke_ objects stored as json. There is another file named _bar.joke_.
* [/internal/demo/config/joke.jq](/internal/demo/config/joke.jq) is a [jq](https://stedolan.github.io/jq/manual/) script that converts _joke_ data files into searchable objects.
* [/internal/demo/templates/joke.tmpl](/internal/demo/templates/joke.tmpl) is a html template that renders one _joke_ object.

## Library

[/pkg](/pkg) should be reusable outside this project.
* [/pkg/dirindex](/pkg/dirindex) can watch a directory for changes and re-index when needed.
* [/pkg/filetypes](/pkg/filetypes) defines interfaces for setting up processing of custom file types.
* [/pkg/searchpage](/pkg/searchpage) serves the search page.

See example usage in [/cmd/pestotrap/main.go](/cmd/pestotrap.main.go) an [/internal/documents](/internal/documents).
