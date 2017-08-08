package main

import (
	"fmt"
	gfx "github.com/crabmusket/lowrezjam2017/graphics"
	game "github.com/crabmusket/lowrezjam2017/game"
	obj "github.com/crabmusket/lowrezjam2017/obj"
	tex "github.com/crabmusket/lowrezjam2017/tex"
	flag "github.com/ogier/pflag"
	"os"
	"runtime/pprof"
)

const (
	TITLE = "CASTLEROOK"
	VERSION = "#LOWREZJAM2017"
	width = 320
	height = 320
)

var (
	flagCpuProfile = flag.String("cpuprofile", "", "output CPU profile information to this file")
	flagWatch = flag.Bool("watch", false, "watch texture and model files for live-reloading")
)

func main() {
	flag.Parse()

	if *flagCpuProfile != "" {
		file, err := os.Create(*flagCpuProfile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	renderer, err := gfx.Init(width, height, TITLE + " - " + VERSION)
	if err != nil {
		panic(err)
	}
	defer gfx.Terminate()

	fmt.Println("OpenGL version", renderer.Version)

	scene, err := game.BuildScene(*flagWatch)
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
