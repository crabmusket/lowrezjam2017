package main

import (
	"fmt"
	gfx "github.com/crabmusket/lowrezjam2017/graphics"
	game "github.com/crabmusket/lowrezjam2017/game"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
)

const (
	VERSION = "#LOWREZJAM2017"
	width = 640
	height = 640
)

func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	renderer, err := gfx.Init(width, height, VERSION)
	if err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	scene, err := game.BuildScene()
	if err != nil {
		panic(err)
	}

	for !renderer.Window.ShouldClose() {
		renderer.Render(func() {
			scene.Render()
		})

		glfw.PollEvents()
	}
}
