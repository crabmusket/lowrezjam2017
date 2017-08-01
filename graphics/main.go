package graphics

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.2-core/gl"
)

func InitWindow(width int, height int, title string) (*glfw.Window, error) {
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

	_, err = initOpenGL()
	if err != nil {
		return nil, err
	}

	return window, nil
}

func initOpenGL() (uint32, error) {
	err := gl.Init()
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.LinkProgram(program)

	gl.ClearColor(0, 0, 0, 1)
	gl.Enable(gl.DEPTH_TEST)

	return program, nil
}
