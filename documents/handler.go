package documents

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
)

type Server struct {
	fs.FS
}

func (s Server) Document(file string, id string) Doc {
	docs, err := Documents(s, file)
	if err != nil {
		return nil
	}

	for _, d := range docs {
		if d[Field_ID] == id {
			return d
		}
	}

	return nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/")
	id := r.URL.Query().Get("id")

	doc := s.Document(filename, id)
	if doc == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "    ")
	enc.SetEscapeHTML(true)

	if _, err := b.Write([]byte("<pre>")); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := enc.Encode(doc); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := b.Write([]byte("</pre>")); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(b.Bytes())
}
