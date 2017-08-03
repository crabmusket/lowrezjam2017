package loadobj

import (
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
}

func (self Object) Render() {
	gl.BindVertexArray(self.Id)
	gl.DrawElements(gl.TRIANGLES, int32(len(self.Indices)), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}
