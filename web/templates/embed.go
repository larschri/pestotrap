package templates

import (
	"embed"
	"html/template"
)

var (
	//go:embed *.tmpl
	files     embed.FS
	Templates = template.Must(template.ParseFS(files, "*.tmpl"))
)
