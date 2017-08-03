package obj

import (
	"github.com/fsnotify/fsnotify"
	"os"
)

type ObjectUpdate struct {
	Data *Object
	Warnings []*Warning
	Object *Object
}

var (
	updates chan *ObjectUpdate
)

func init() {
	updates = make(chan *ObjectUpdate, 10)
}

func ProcessUpdates(warn func([]*Warning)) {
	select {
	case update := <-updates:
		// TODO: release buffer
		update.Object.Indices = make([]uint32, len(update.Data.Indices))
		update.Object.Vertices = make([]float32, len(update.Data.Vertices))
		copy(update.Object.Indices, update.Data.Indices)
		copy(update.Object.Vertices, update.Data.Vertices)
		update.Object.Unbind()
		update.Object.Bind()
		if update.Warnings != nil && len(update.Warnings) > 0 && warn != nil {
			warn(update.Warnings)
		}

	default:
		// do nothing
	}
}

func (self *Object) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Add(self.Filename)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Name == self.Filename {
					file, err := os.Open(self.Filename)
					if err != nil {
						continue
					}
					defer file.Close()

					obj, warnings, err := Read(file)
					if err != nil {
						continue
					}

					update := &ObjectUpdate{
						Data: obj,
						Warnings: warnings,
						Object: self,
					}
					updates <- update
				}

			default:
				// do nothing
			}
		}
	}()

	return nil
}
