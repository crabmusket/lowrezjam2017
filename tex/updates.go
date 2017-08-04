package textures

import (
	"github.com/fsnotify/fsnotify"
	"image"
)

type TextureUpdate struct {
	Data []byte
	Size *image.Point
	Texture *Texture
}

var (
	updates chan *TextureUpdate
)

func init() {
	updates = make(chan *TextureUpdate, 10)
}

// You must call this function from the main thread which is running opengl.
func ProcessUpdates() {
	select {
	case update := <-updates:
		// TODO: release textures
		update.Texture.Bind(update.Data, update.Size)
	default:
		// do nothing
	}
}

func (self *Texture) Watch() error {
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
					data, size, err := loadImage(self.Filename)
					if err != nil {
						continue
					}
					update := &TextureUpdate{
						Data: data,
						Size: size,
						Texture: self,
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
