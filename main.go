package main

import (
	"fmt"
	gfx "github.com/crabmusket/lowrezjam2017/graphics"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
)

const (
	VERSION = "#LOWREZJAM2017"
	width = 600
	height = 600
)

func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := gfx.InitWindow(width, height, VERSION)
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
