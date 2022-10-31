package documents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

type Server struct {
	fs.FS
}

func (s Server) lookup(r *http.Request) (any, error) {
	docs, err := Documents(s, strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		return nil, fmt.Errorf("%v: not found", r.URL.Path)
	}

	id := r.URL.Query().Get("id")
	if id != "" {
		for _, d := range docs {
			if d[Field_ID] == id {
				return d, nil
			}
		}

		return nil, fmt.Errorf("%v: not found", id)
	}

	return docs, nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o, err := s.lookup(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "    ")

	if _, err = b.Write([]byte("<pre>")); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if err := enc.Encode(o); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if _, err = b.Write([]byte("</pre>")); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	_, _ = w.Write(b.Bytes())
}
