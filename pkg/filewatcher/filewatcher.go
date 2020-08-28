// Package filewatcher wraps fsnotify
// into a simple, callback based API.
package filewatcher

import (
	"io/ioutil"
	"path"

	"gopkg.in/fsnotify.v1"
)

// Handler is the callback function for an
// fsnotify.Event.
type Handler func(event fsnotify.Event)

// ErrorHandler is the handler function for
// occuring errors.
type ErrorHandler func(err error)

// handlerWrap wraps a callback handler
// together with its designated op code.
type handlerWrap struct {
	handler Handler
	op      fsnotify.Op
}

// FileWatcher wraps a fsnotify.Watcher to
// register event callbacks on.
type FileWatcher struct {
	callbacks             []handlerWrap
	errorCallbacks        []ErrorHandler
	watcher               *fsnotify.Watcher
	paths                 []string
	emitExistentAsCreated bool
	done                  chan struct{}
}

// New initializes a new instance of FileWatcher.
//
// When emitExistentAsCreated is passed as true,
// existing files in the specified directories
// are treated as 'created' on calling `Run`
// which empits an fsnotify.Create event emit
// for each existing file.
func New(emitExistentAsCreated bool) (fw *FileWatcher, err error) {
	fw = &FileWatcher{
		callbacks:             make([]handlerWrap, 0),
		errorCallbacks:        make([]ErrorHandler, 0),
		paths:                 make([]string, 0),
		emitExistentAsCreated: emitExistentAsCreated,
		done:                  make(chan struct{}),
	}

	if fw.watcher, err = fsnotify.NewWatcher(); err != nil {
		return
	}

	return
}

// AddPath adds a path to be watched for fs events.
// If the path bind fails, an error is returned.
func (fw *FileWatcher) AddPath(path string) error {
	if fw.emitExistentAsCreated {
		fw.paths = append(fw.paths, path)
	}

	return fw.watcher.Add(path)
}

// Handle registers the passed event handler on the
// specified op code.
//
// The op code can also be a bit-wise combination of
// multiple op codes using bitwise-or.
//
// Example:
// 	op := fsnotify.Create | fsnotify.Write
func (fw *FileWatcher) Handle(op fsnotify.Op, handler Handler) {
	fw.callbacks = append(fw.callbacks, handlerWrap{
		handler, op,
	})
}

// HandleError registers an event handler callback
// for risen events.
func (fw *FileWatcher) HandleError(handler ErrorHandler) {
	fw.errorCallbacks = append(fw.errorCallbacks, handler)
}

// Run initializes the event loop goroutine.
func (fw *FileWatcher) Run() {
	if fw.emitExistentAsCreated {
		fw.emitExistingFiles()
	}

	go fw.eventLoop()
}

// Close closes the watchers event loop.
func (fw *FileWatcher) Close() error {
	return fw.watcher.Close()
}

// Done returns a channel which receives an empty
// struct instance when the event loop was closed.
func (fw *FileWatcher) Done() <-chan struct{} {
	return fw.done
}

func (fw *FileWatcher) eventLoop() {
loop:
	for {
		select {

		case event, ok := <-fw.watcher.Events:
			if !ok {
				break loop
			}
			fw.emitEvent(event)

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				break loop
			}
			fw.emitErrorEvent(err)
		}
	}

	fw.done <- struct{}{}
}

func (fw *FileWatcher) emitEvent(event fsnotify.Event) {
	for _, h := range fw.callbacks {
		if h.op&event.Op == event.Op {
			h.handler(event)
		}
	}
}

func (fw *FileWatcher) emitErrorEvent(err error) {
	for _, cb := range fw.errorCallbacks {
		cb(err)
	}
}

func (fw *FileWatcher) emitExistingFiles() {
	for _, p := range fw.paths {
		dir, err := ioutil.ReadDir(p)
		if err != nil {
			fw.emitErrorEvent(err)
			continue
		}
		for _, d := range dir {
			fw.emitEvent(fsnotify.Event{
				Name: path.Join(p, d.Name()),
				Op:   fsnotify.Create,
			})
		}
	}
}
