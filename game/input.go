package game

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

func ProcessInput() {
	glfw.PollEvents()
}
