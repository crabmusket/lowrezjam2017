package game

import (
	"fmt"
	gfx "github.com/crabmusket/lowrezjam2017/graphics"
	obj "github.com/crabmusket/lowrezjam2017/obj"
	tex "github.com/crabmusket/lowrezjam2017/tex"
	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Scene struct{
	Camera *Camera
	Level *StaticRendered
}

type Camera struct{
	Position mgl.Vec3
	Pitch float32
	Yaw float32
	Transform mgl.Mat4
	Projection mgl.Mat4
}

type StaticRendered struct{
	Transform mgl.Mat4
	Texture *tex.Texture
	Geometry *obj.Object
	Shader uint32
}

func BuildScene() (*Scene, error) {
	levelTexture, err := tex.Load("resources/textures/level.png")
	if err != nil {
		return nil, err
	}

	if true {
		err := levelTexture.Watch()
		if err != nil {
			return nil, err
		}
	}

	level1, warnings, err := obj.Load("resources/meshes/floor1.obj")
	if err != nil {
		return nil, err
	}
	for _, warning := range warnings {
		fmt.Printf("%+v\n", warning)
	}

	if true {
		err := level1.Watch()
		if err != nil {
			return nil, err
		}
	}

	staticShader, err := gfx.MakeProgram("resources/shaders/static.vert.glsl", "resources/shaders/static.frag.glsl")
	if err != nil {
		return nil, err
	}

	scene := &Scene{
		Camera: &Camera{
			Position: mgl.Vec3{0, 0, 0},
			Pitch: 0,
			Yaw: 0,
			Transform: mgl.Translate3D(0, 0, 0),
			Projection: mgl.Perspective(mgl.DegToRad(45), 1, 0.1, 100),
		},

		Level: &StaticRendered{
			Transform: mgl.Translate3D(0, 0, 0),
			Texture: levelTexture,
			Geometry: level1,
			Shader: staticShader,
		},
	}

	return scene, nil
}

func (self Scene) Render() {
	program := self.Level.Shader
	gl.UseProgram(program)

	gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("ambient\x00")), 0.5)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &self.Camera.Projection[0])
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("view\x00")), 1, false, &self.Camera.Transform[0])
	l := gl.GetUniformLocation(program, gl.Str("cameraPos\x00"))
	p := self.Camera.Position;
	gl.Uniform3f(l, p[0], p[1], p[2])

	// Render the level
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &self.Level.Transform[0])
	gl.BindTexture(gl.TEXTURE_2D, self.Level.Texture.Id)
	self.Level.Geometry.Render()
}
