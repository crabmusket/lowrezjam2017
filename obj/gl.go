package obj

import (
	tex "github.com/crabmusket/lowrezjam2017/tex"
	"github.com/go-gl/gl/v3.2-core/gl"
)

func (self *Object) Bind() {
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &ebo)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4 * len(self.Vertices), gl.Ptr(self.Vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4 * len(self.Indices), gl.Ptr(self.Indices), gl.STATIC_DRAW)

	// positions
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, nil)
	gl.EnableVertexAttribArray(0)

	// normals
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// tex coords
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	self.Id = vao
	self.Vbo = vbo
	self.Ebo = ebo
}

func (self *Object) Unbind() {
	gl.DeleteVertexArrays(1, &self.Id)
	gl.DeleteBuffers(1, &self.Vbo)
	gl.DeleteBuffers(1, &self.Ebo)
}

func (self Object) Render(textures tex.Library) {
	gl.BindVertexArray(self.Id)

	for _, material := range(self.Materials) {
		texture := textures[material.Name]
		if texture == nil {
			continue
		}
		gl.BindTexture(gl.TEXTURE_2D, texture.Id)

		span := int32(material.End - material.Start)
		begin := gl.PtrOffset(4 * int(material.Start))
		gl.DrawElements(gl.TRIANGLES, span, gl.UNSIGNED_INT, begin)
	}

	gl.BindVertexArray(0)
}
