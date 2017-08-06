package main

import (
	"fmt"
	gfx "github.com/crabmusket/lowrezjam2017/graphics"
	game "github.com/crabmusket/lowrezjam2017/game"
	obj "github.com/crabmusket/lowrezjam2017/obj"
	tex "github.com/crabmusket/lowrezjam2017/tex"
)

const (
	TITLE = "CASTLEROOK"
	VERSION = "#LOWREZJAM2017"
	width = 320
	height = 320
)

func main() {
	renderer, err := gfx.Init(width, height, TITLE + " - " + VERSION)
	if err != nil {
		panic(err)
	}
	defer gfx.Terminate()

	fmt.Println("OpenGL version", renderer.Version)

	scene, err := game.BuildScene()
	if err != nil {
		panic(err)
	}

	gfx.CheckAndPrintErrors()

	game.InitInput(renderer.Window)

	for renderer.Run() {
		tex.ProcessUpdates()
		obj.ProcessUpdates(nil)

		renderer.Render(func() {
			scene.Render()
		})

		gfx.CheckAndPrintErrors()

		if game.ProcessInput(renderer.Window, scene) == false {
			break
		}
	}
}
