package documents

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/larschri/pestotrap/pkg/filetypes"
)

type Server struct {
	fs.FS
	*template.Template
}

var errNotFound = errors.New("Resource was not found")

// NewServer creates a new Server object. The given fs.FS must have a
// particular structure.
//   - config - a directory with one .jq file for each filetype
//   - templates - a directory with one .tmpl file for each filetype
//   - data - a directory of data files. The file extension defines the
//     filetype.
//
// See ../demo for example
func NewServer(cfg fs.FS) (*Server, error) {
	entries, err := fs.ReadDir(cfg, "config")
	if err != nil {
		return nil, fmt.Errorf("failed to read config dir: %w", err)
	}

	tpls, err := fs.Sub(cfg, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to open templates dir: %w", err)
	}

	dataFS, err := fs.Sub(cfg, "data")
	if err != nil {
		return nil, fmt.Errorf("failed to open data dir: %w", err)
	}

	tpl, err := template.ParseFS(tpls, "*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	for _, e := range entries {
		f, err := cfg.Open("config/" + e.Name())
		if err != nil {
			return nil, err
		}

		bs, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		base, ok := strings.CutSuffix(e.Name(), ".jq")
		if !ok {
			continue
		}

		tmpl := tpl.Lookup(base + ".tmpl")
		if tmpl == nil {
			return nil, fmt.Errorf("No template for %s", base)
		}

		if err := registerFiletype("."+base, string(bs)); err != nil {
			return nil, err
		}
	}

	return &Server{dataFS, tpl}, nil
}

// document return the document with the given id from the give fname. The
// entire file is read and parsed every time this function is invoked. Consider
// something more efficient when becomes too slow.
func (s Server) document(fname string, id string) (filetypes.Doc, error) {
	docs, err := filetypes.Documents(s.FS, fname)
	if err != nil {
		return nil, err
	}

	for _, d := range docs {
		if d.ID() == id {
			return d, nil
		}
	}

	return nil, errNotFound
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/")
	id := r.URL.Query().Get("id")

	doc, err := s.document(filename, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := s.Template.Lookup(strings.TrimPrefix(path.Ext(filename), ".") + ".tmpl")
	if tmpl != nil {
		tmpl.Execute(w, doc)
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
