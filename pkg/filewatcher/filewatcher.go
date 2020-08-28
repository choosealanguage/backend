package filewatcher

import (
	"io/ioutil"
	"path"

	"gopkg.in/fsnotify.v1"
)

type Handler func(event fsnotify.Event)
type ErrorHandler func(err error)

type handlerReg struct {
	handler Handler
	op      fsnotify.Op
}

type FileWatcher struct {
	callbacks             []handlerReg
	errorCallbacks        []ErrorHandler
	watcher               *fsnotify.Watcher
	paths                 []string
	emitExistentAsCreated bool
	done                  chan struct{}
}

func New(emitExistentAsCreated bool) (fw *FileWatcher, err error) {
	fw = &FileWatcher{
		callbacks:             make([]handlerReg, 0),
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

func (fw *FileWatcher) AddPath(path string) error {
	if fw.emitExistentAsCreated {
		fw.paths = append(fw.paths, path)
	}

	return fw.watcher.Add(path)
}

func (fw *FileWatcher) Handle(op fsnotify.Op, handler Handler) {
	fw.callbacks = append(fw.callbacks, handlerReg{
		handler, op,
	})
}

func (fw *FileWatcher) HandleError(handler ErrorHandler) {
	fw.errorCallbacks = append(fw.errorCallbacks, handler)
}

func (fw *FileWatcher) Start() {
	if fw.emitExistentAsCreated {
		fw.emitExistingFiles()
	}

	go fw.eventLoop()
}

func (fw *FileWatcher) Close() error {
	return fw.watcher.Close()
}

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
