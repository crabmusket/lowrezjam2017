package textures

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"image"
	"image/draw"
	_ "image/png"
	_ "image/jpeg"
	"os"
	"path/filepath"
)

type Texture struct {
	Id uint32
	Filename string
}

func Load(filename string, library Library) (*Texture, error) {
	data, size, err := loadImage(filename)
	if err != nil {
		return nil, err
	}

	texture := &Texture{
		Filename: filename,
	}
	texture.Bind(data, size)

	if library != nil {
		textureName := filepath.Base(filename)
		extension := filepath.Ext(filename)
		textureKey := textureName[0:len(textureName)-len(extension)]
		library[textureKey] = texture
	}

	return texture, nil
}

func (self *Texture) Bind(data []byte, size *image.Point) {
	tex := self.Id
	if tex == 0 {
		gl.GenTextures(1, &tex)
	}

	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT);
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT);
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR);
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);

	gl.BindTexture(gl.TEXTURE_2D, 0)
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
