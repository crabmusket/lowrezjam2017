package game

import (
	"fmt"
	gfx "github.com/crabmusket/lowrezjam2017/graphics"
	obj "github.com/crabmusket/lowrezjam2017/obj"
	tex "github.com/crabmusket/lowrezjam2017/tex"
	"github.com/go-gl/gl/v3.2-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"strconv"
)

type Scene struct{
	Camera *Camera
	Level *StaticRendered
	Textures tex.Library
	Lights []*Light
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
	Geometry *obj.Object
	Shader uint32
}

type Light struct{
	Position mgl.Vec3
	Colour mgl.Vec3
	Radius float32
}

var (
	textures = []string{
		"resources/textures/wall_stone.png",
		"resources/textures/wall_plain.png",
		"resources/textures/roof_wood.png",
		"resources/textures/floor_tiled.png",
	}
)

func BuildScene(watch bool) (*Scene, error) {
	library := tex.MakeLibrary()

	for _, file := range(textures) {
		texture, err := tex.Load(file, library)
		if err != nil {
			return nil, err
		}

		if watch {
			err := texture.Watch()
			if err != nil {
				return nil, err
			}
		}
	}

	level1, warnings, err := obj.Load("resources/meshes/floor1.obj")
	if err != nil {
		return nil, err
	}
	for _, warning := range warnings {
		fmt.Printf("%+v\n", warning)
	}

	if watch {
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
			Projection: mgl.Perspective(mgl.DegToRad(60), 1, 0.1, 100),
		},

		Level: &StaticRendered{
			Transform: mgl.Translate3D(0, 0, 0),
			Geometry: level1,
			Shader: staticShader,
		},

		Lights: []*Light{
			&Light{
				Position: mgl.Vec3{1, 0, -4},
				Colour: mgl.Vec3{1, 0.85, 0.5},
				Radius: 3,
			},
			&Light{
				Position: mgl.Vec3{-5.7, -0.6, -6.1},
				Colour: mgl.Vec3{0.6, 0.88, 1},
				Radius: 5,
			},
		},

		Textures: library,
	}

	return scene, nil
}

func (self Scene) Render() {
	program := self.Level.Shader
	gl.UseProgram(program)

	// Lighting
	gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("ambient\x00")), 0.05)
	for i, light := range(self.Lights) {
		if i > 5 {
			break
		}
		is := strconv.Itoa(i)
		gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("pointLights[" + is + "].position\x00")), light.Position[0], light.Position[1], light.Position[2])
		gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("pointLights[" + is + "].diffuseColour\x00")), light.Colour[0], light.Colour[1], light.Colour[2])
		gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("pointLights[" + is + "].radius\x00")), light.Radius)
	}

	// Transforms
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &self.Camera.Projection[0])
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("view\x00")), 1, false, &self.Camera.Transform[0])
	p := self.Camera.Position;
	gl.Uniform3f(gl.GetUniformLocation(program, gl.Str("cameraPos\x00")), p[0], p[1], p[2])

	// Render the level
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &self.Level.Transform[0])
	self.Level.Geometry.Render(self.Textures)
}
