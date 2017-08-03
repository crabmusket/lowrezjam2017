package main

import (
	"fmt"
	gfx "github.com/crabmusket/lowrezjam2017/graphics"
	game "github.com/crabmusket/lowrezjam2017/game"
	obj "github.com/crabmusket/lowrezjam2017/obj"
	tex "github.com/crabmusket/lowrezjam2017/tex"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
)

const (
	VERSION = "#LOWREZJAM2017"
	width = 320
	height = 320
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

	gfx.CheckAndPrintErrors()

	for !renderer.Window.ShouldClose() {
		tex.ProcessUpdates()
		obj.ProcessUpdates(nil)

		renderer.Render(func() {
			scene.Render()
		})

		gfx.CheckAndPrintErrors()

		glfw.PollEvents()
	}
}
