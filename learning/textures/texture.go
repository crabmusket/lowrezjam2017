package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-gl/gl/v3.2-core/gl"
	"image"
	"image/draw"
	_ "image/png"
	_ "image/jpeg"
	"os"
)

type Texture struct {
	Id uint32
	Filename string
}

type TextureUpdate struct {
	Data []byte
	Size *image.Point
	Texture *Texture
}

func makeTexture(filename string, updates chan<- *TextureUpdate) (*Texture, error) {
	data, size, err := loadImage(filename)
	if err != nil {
		return nil, err
	}

	texture := &Texture{
		Filename: filename,
	}
	texture.Bind(data, size)

	err = texture.Watch(updates)
	if err != nil {
		return nil, err
	}

	return texture, nil
}

func (self *Texture) Bind(data []byte, size *image.Point) {
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	self.Id = tex
}

func loadImage(filename string) ([]byte, *image.Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)
	size := rgba.Rect.Size()

	return rgba.Pix, &size, nil
}

func (self *Texture) Watch(updates chan<- *TextureUpdate) error {
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
				fmt.Println(event.Name)
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

func (self *TextureUpdate) Process() {
	self.Texture.Bind(self.Data, self.Size)
}
