package documents

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
)

var (
	MatchFields = []string{
		Field_Name,
		Field_Type,
		Field_Taxonomy,
		Field_Filename,
		Field_ID,
	}

	//go:embed match.htmx
	matchTemplateContent string

	matchTempl = template.Must(template.New("").Parse(matchTemplateContent))
)

func RenderMatch(w io.Writer, doc map[string]interface{}) {
	matchTempl.Execute(w, map[string]interface{}{
		"Name": doc[Field_Name],
		"Type": doc[Field_Type],
		"Taxonomy": fmt.Sprintf("%v / %v",
			doc[Field_Filename],
			doc[Field_Taxonomy]),
		"Url": fmt.Sprintf("/x/%v?id=%v",
			doc[Field_Filename],
			doc[Field_ID]),
	})
}
