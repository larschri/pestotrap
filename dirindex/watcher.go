package dirindex

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"gopkg.in/fsnotify.v1"
)

type Watcher struct {
	dir   fs.FS
	index bleve.Index
	state dirstate
	mu    sync.Mutex
}

func NewWatcher(dir fs.FS, index bleve.Index) *Watcher {
	return &Watcher{dir, index, nil, sync.Mutex{}}
}

// UpdateIfModified updates the index if one or more files have been updated.
// It is thread safe and will block if another update is already in progress.
func (r *Watcher) UpdateIfModified() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	newState, err := newDirstate(r.dir)
	if err != nil {
		return fmt.Errorf("failed to create dirState: %w", err)
	}

	if newState.equal(r.state) {
		return nil
	}

	if err := update(r.index, r.dir); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	r.state = newState
	return nil
}

// Watch the directory and update the directory until ctx is closed or another
// error occurs.
func (r *Watcher) Watch(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create new watcher: %w", err)
	}
	defer watcher.Close()

	// Convert to string - works if r.dir was created by os.DirFS
	dir := fmt.Sprintf("%v", r.dir)
	if err := watcher.Add(dir); err != nil {
		return fmt.Errorf("failed to add directory %v to watcher: %w", r.dir, err)
	}

	if err := r.UpdateIfModified(); err != nil {
		return err
	}

	for {
		select {
		case _, ok := <-watcher.Events:

			if !ok {
				return errors.New("failed to read fsnotify.Watcher.Events")
			}

			if err := r.UpdateIfModified(); err != nil {
				return err
			}

		case err, ok := <-watcher.Errors:

			if !ok {
				return errors.New("failed to read fsnotify.Watcher.Errors")
			}
			return fmt.Errorf("error while watching: %w", err)

		case <-ctx.Done():
			return nil
		}
	}

	return nil
}
