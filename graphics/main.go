package graphics

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.2-core/gl"
	"io/ioutil"
	"strings"
	"runtime"
)

const (
	realWidth = 64
	realHeight = 64
)

type Renderer struct{
	Window *glfw.Window
	Version string
	Shader uint32
	Framebuffer uint32
	Texture uint32
	Plane uint32
}

func Init(width int, height int, title string) (*Renderer, error) {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()

	err = initOpenGL()
	if err != nil {
		return nil, err
	}

	screenShader, err := MakeProgram("resources/shaders/screen.vert.glsl", "resources/shaders/screen.frag.glsl")
	if err != nil {
		return nil, err
	}

	framebuffer, texture, err := makeFramebuffer()
	if err != nil {
		return nil, err
	}

	vao := bindGeometry()

	renderer := &Renderer{
		Window: window,
		Version: gl.GoStr(gl.GetString(gl.VERSION)),
		Shader: screenShader,
		Framebuffer: framebuffer,
		Texture: texture,
		Plane: vao,
	}

	return renderer, nil
}

func Terminate() {
	defer glfw.Terminate()
}

func MakeProgram(vert string, frag string) (uint32, error) {
	vertexShader, err := compileShader(vert, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(frag, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	return program, nil
}

func initOpenGL() error {
	err := gl.Init()
	if err != nil {
		return err
	}

	gl.ClearColor(0, 0, 0, 1)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.FRAMEBUFFER_SRGB);
	gl.Enable(gl.CULL_FACE);

	return nil
}

func (self Renderer) Run() bool {
	return !self.Window.ShouldClose()
}

func (self Renderer) Render(renderScene func()) {
	// 1. render scene to framebuffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, self.Framebuffer)
	gl.Viewport(0, 0, realWidth, realHeight)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)

	renderScene()

	// 2. render framebuffer to quad as texture
	width, height := self.Window.GetSize()
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Disable(gl.DEPTH_TEST)
	gl.UseProgram(self.Shader)
	gl.BindTexture(gl.TEXTURE_2D, self.Texture)
	gl.BindVertexArray(self.Plane)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	self.Window.SwapBuffers()
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

func makeFramebuffer() (uint32, uint32, error) {
	var fb uint32
	gl.GenFramebuffers(1, &fb)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb)

	var colour uint32
	gl.GenTextures(1,  &colour)
	gl.BindTexture(gl.TEXTURE_2D, colour)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, realWidth, realHeight, 0,  gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, colour, 0)

	var rbo uint32
	gl.GenRenderbuffers(1, &rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, realWidth, realHeight)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, rbo)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return 0, 0, fmt.Errorf("framebuffer is not complete")
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return fb, colour, nil
}

var (
	planeVerts = []float32{
		// pos  tex
		-1, -1, 0, 0,
		1, -1, 1, 0,
		1, 1, 1, 1,
		-1, 1, 0, 1,
	}

	planeIndices = []uint32{
		0, 1, 2,
		0, 2, 3,
	}
)

func bindGeometry() uint32 {
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &ebo)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4 * len(planeVerts), gl.Ptr(planeVerts), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4 * len(planeIndices), gl.Ptr(planeIndices), gl.STATIC_DRAW)

	// positions
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, nil)
	gl.EnableVertexAttribArray(0)

	// tex coords
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return vao
}

func CheckAndPrintErrors() {
	for {
		errorCode := gl.GetError()
		if errorCode == gl.NO_ERROR {
			break
		}

		err := fmt.Errorf("OpenGL error:", errorCode)
		fmt.Println(err.Error())
	}
}
