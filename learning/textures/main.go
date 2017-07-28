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
	triangleVerts = []float32{
		0, -0.5, -1,
		-0.5, 0.5, -1,
		0.5, 0.5, -1,
	}

	// http://www.colourlovers.com/palette/3501633/HV
	triangleColours = []float32{
		0.125, 0.431, 0.549,
		1, 0.282, 0.27,
		1, 0.859, 0.078,
	}

	triangleTextureCoords = []float32{
		0.5, 1,
		0, 0,
		1, 0,
	}

	triangle = [][]float32{
		triangleVerts,
		triangleColours,
		triangleTextureCoords,
	}

	rectVerts = []float32{
		-1, 1, -0.5,
		1, 1,  -0.5,
		1, -1, -0.5,
		-1, -1, -0.5,
	}

	rectColours = []float32{
		0.125, 0.431, 0.549,
		1, 0.282, 0.27,
		1, 0.859, 0.078,
		0.211, 0.098, 0.239,
	}

	rectTextureCoords = []float32{
		0, 0,
		0, 1,
		1, 1,
		1, 0,
	}

	rect = [][]float32{
		rectVerts,
		rectColours,
		rectTextureCoords,
	}

	rectIndices = []uint32{
		0, 1, 3,
		1, 2, 3,
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

	textureUpdates := make(chan *TextureUpdate, 10)
	
	triangleTex, err := makeTexture("texture-tiles.png", textureUpdates)
	if err != nil {
		panic(err)
	}

	gridTex, err := makeTexture("texture-grid.png", textureUpdates)
	if err != nil {
		panic(err)
	}

	triangleVAO := makeVAO(3, triangle, nil)
	rectVAO := makeVAO(4, rect, rectIndices)
	for !window.ShouldClose() {
		select {
		case update := <-textureUpdates:
			update.Process()
		default:
			// do nothing
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		gl.BindTexture(gl.TEXTURE_2D, triangleTex.Id)
		gl.BindVertexArray(triangleVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangleVerts) / 3))

		gl.BindTexture(gl.TEXTURE_2D, gridTex.Id)
		gl.BindVertexArray(rectVAO)
		gl.DrawElements(gl.TRIANGLES, int32(len(rectIndices)), gl.UNSIGNED_INT, nil)

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

func makeVAO(size int, arrays [][]float32, indices []uint32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var ebo uint32
	gl.GenBuffers(1, &ebo)

	vbos := make([]uint32, len(arrays))
	gl.GenBuffers(int32(len(arrays)), &vbos[0])
	for i, points := range arrays {
		gl.BindBuffer(gl.ARRAY_BUFFER, vbos[i])
		gl.BufferData(gl.ARRAY_BUFFER, 4 * len(points), gl.Ptr(&arrays[i][0]), gl.STATIC_DRAW)

		if indices != nil {
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
			gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4 * len(indices), gl.Ptr(&indices[0]), gl.STATIC_DRAW)
		}

		gl.VertexAttribPointer(uint32(i), int32(len(points) / size), gl.FLOAT, false, 0, nil)
		gl.EnableVertexAttribArray(uint32(i))
	}

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
