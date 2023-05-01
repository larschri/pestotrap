package dirindex

import (
	"fmt"
	"io/fs"
)

type dirstate []string

func filestate(e fs.DirEntry) string {
	f, err := e.Info()
	if err != nil {
		return e.Name()
	}
	return fmt.Sprintf("%s, %v", e.Name(), f.ModTime())
}

func (d dirstate) equal(other dirstate) bool {
	if len(d) != len(other) {
		return false
	}

	for i := range d {
		if d[i] != other[i] {
			return false
		}
	}

	return true
}

func newDirstate(dir fs.FS) (dirstate, error) {
	es, err := fs.ReadDir(dir, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %w", dir, err)
	}

	d := make([]string, len(es))
	for _, e := range es {
		d = append(d, filestate(e))
	}

	return dirstate(d), nil
}
