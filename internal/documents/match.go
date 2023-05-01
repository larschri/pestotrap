package documents

import (
	_ "embed"
	"fmt"

	"github.com/blevesearch/bleve/v2/search"
	"github.com/larschri/pestotrap/pkg/filetypes"
)

var (
	matchFields = []string{
		field_Title,
		field_Description,
		filetypes.Field_Filename,
		field_ID,
	}
)

type Match search.DocumentMatch

func (m *Match) Title() interface{} {
	return m.Fields[field_Title]
}

func (m *Match) Description() interface{} {
	d, ok := m.Fields[field_Description]
	if !ok {
		return fmt.Sprintf("Document %v from file %v",
			m.Fields[field_ID],
			m.Fields[filetypes.Field_Filename])
	}
	return d
}

func (m *Match) Url() interface{} {
	return fmt.Sprintf("/x/%v?id=%v",
		m.Fields[filetypes.Field_Filename],
		m.Fields[field_ID])
}
