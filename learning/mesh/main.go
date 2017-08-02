package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"runtime"
	"strings"
)

const (
	width = 640
	height = 640
)

var (
	rectVerts = []float32{
		0.5, 0.5, 0, // position
		0, 1, 0, // normal
		1, 0, // texc

		0.5, -0.5, 0,
		0, 1, 0,
		1, 1,

		-0.5, -0.5, 0,
		0, 1, 0,
		0, 1,

		-0.5, 0.5, 0,
		0, 1, 0,
		0, 0,
	}

	rectIndices = []uint32{
		0, 2, 1,
		2, 0, 3,
	}
)

func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := initWindow(width, height, "LOWREZJAM")
	if err != nil {
		panic(err)
	}

	program, err := initOpenGL()
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	objs, warnings, err := loadObj("./suzanne.obj")
	if err != nil {
		panic(err)
	} else if len(objs) == 0 {
		panic(fmt.Errorf("no objects parsed from suzanne.obj"))
	} else if len(warnings) > 0 {
		for _, warning := range warnings {
			fmt.Println(warning)
		}
	}
	obj := objs[0]

	projection := mgl.Perspective(mgl.DegToRad(45), 1, 0.1, 100)
	view := mgl.Translate3D(0, 0, -2)
	rectTransform := mgl.Translate3D(0, -0.5, 0).
		Mul4(mgl.HomogRotate3D(mgl.DegToRad(-90), mgl.Vec3{1, 0, 0})).
		Mul4(mgl.Scale3D(2, 2, 2))

	textureUpdates := makeTextureUpdates()

	containerTex, err := makeTexture("texture-container.png", textureUpdates)
	if err != nil {
		panic(err)
	}

	gridTex, err := makeTexture("texture-tiles.png", textureUpdates)
	if err != nil {
		panic(err)
	}

	rectVAO := makeVAO(rectVerts, rectIndices)
	meshVAO := makeVAO(obj.Vertices, nil) //obj.Indices)
	for !window.ShouldClose() {
		processTextureUpdates(textureUpdates)

		meshTransform := mgl.Translate3D(0, 0, 0).
			Mul4(mgl.HomogRotate3D(mgl.DegToRad(-15 * float32(glfw.GetTime())), mgl.Vec3{0, 1, 0}.Normalize())).
			Mul4(mgl.Scale3D(0.5, 0.5, 0.5))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("view\x00")), 1, false, &view[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &projection[0])

		gl.BindVertexArray(rectVAO)
		gl.BindTexture(gl.TEXTURE_2D, containerTex.Id)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &rectTransform[0])
		gl.DrawElements(gl.TRIANGLES, int32(len(rectIndices)), gl.UNSIGNED_INT, nil)

		gl.BindVertexArray(meshVAO)
		gl.BindTexture(gl.TEXTURE_2D, gridTex.Id)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("model\x00")), 1, false, &meshTransform[0])
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(obj.Vertices) / 8))

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func initWindow(width int, height int, title string) (window *glfw.Window, err error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err == nil {
		window.MakeContextCurrent()
	}

	return
}

func initOpenGL() (program uint32, err error) {
	if err = gl.Init(); err != nil {
		return
	}

	vertexShader, err := compileShader("./vertex.glsl", gl.VERTEX_SHADER)
	if err != nil {
		return
	}

	fragmentShader, err := compileShader("./fragment.glsl", gl.FRAGMENT_SHADER)
	if err != nil {
		return
	}

	program = gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	gl.ClearColor(0.211, 0.098, 0.239, 1)
	gl.Enable(gl.DEPTH_TEST)

	return
}

func makeVAO(data []float32, indices []uint32) uint32 {
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &ebo)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4 * len(data), gl.Ptr(data), gl.STATIC_DRAW)

	if indices != nil {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4 * len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	}

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

	return vao
}

func compileShader(filename string, vertexOrFragment uint32) (shader uint32, err error) {
	sourceBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	shader = gl.CreateShader(vertexOrFragment)

	cSource, free := gl.Strs(string(sourceBytes) + "\x00")
	gl.ShaderSource(shader, 1, cSource, nil)
	free()

	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength + 1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		shader = 0
		err = fmt.Errorf("failed to compile %v: %v", filename, log)
	}

	return
}
