package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"io/ioutil"
	"runtime"
	"strings"
)

const (
	width = 640
	height = 640
)

var (
	// http://www.colourlovers.com/palette/3501633/HV
	triangle = []float32{
		0, -0.5, 0, /* */ 0.125, 0.431, 0.549,
		-0.5, 0.5, 0, /* */ 1, 0.282, 0.27,
		0.5, 0.5, 0, /* */ 1, 0.859, 0.078,
	}
)

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	if window, err := initWindow(width, height, "LOWREZJAM"); err != nil {
		panic(err)
	} else if program, err := initOpenGL(); err != nil {
		panic(err)
	} else {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Println("OpenGL version", version)

		vao := makeVAO(triangle)
		for !window.ShouldClose() {
			draw(window, program, vao)
		}
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

	return
}

func draw(window *glfw.Window, program uint32, vao uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle) / 3))

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeVAO(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4 * len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6 * 4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6 * 4, gl.PtrOffset(3 * 4))
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

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
