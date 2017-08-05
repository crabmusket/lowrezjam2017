package game

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
	"math"
)

const (
	moveSpeed float32 = 3
	pitchSpeed float32 = 1
	yawSpeed float32 = 2
)

var (
	time float64
	exit bool
)

func InitInput(window *glfw.Window) {
	time = glfw.GetTime()
	exit = false
	window.SetKeyCallback(ProcessKey)
}

func ProcessKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		exit = true
	}
}

func ProcessInput(window *glfw.Window, scene *Scene) bool {
	t := glfw.GetTime()
	dt := float32(t - time)
	time = t

	glfw.PollEvents()

	// Deal with camera rotation changes first
	du := keyToInt(window.GetKey(glfw.KeyUp))
	dd := keyToInt(window.GetKey(glfw.KeyDown))
	dl := keyToInt(window.GetKey(glfw.KeyLeft))
	dr := keyToInt(window.GetKey(glfw.KeyRight))

	yaw := dl - dr
	pitch := dd - du
	scene.Camera.Yaw += yawSpeed * yaw * dt
	scene.Camera.Pitch += pitchSpeed * pitch * dt
	threshold := float32(math.Pi)/2 - 0.01
	if scene.Camera.Pitch > threshold {
		scene.Camera.Pitch = threshold
	}
	if scene.Camera.Pitch < -threshold {
		scene.Camera.Pitch = -threshold
	}

	// Now we can calculate the camera's direction to use below
	yawMat := mgl.HomogRotate3DY(scene.Camera.Yaw)
	cameraRightV := mgl.TransformNormal(mgl.Vec3{1, 0, 0}, yawMat)
	cameraFrontV := mgl.TransformNormal(mgl.TransformNormal(mgl.Vec3{0, 0, -1}, yawMat), mgl.HomogRotate3D(scene.Camera.Pitch, cameraRightV))

	// Movement relative to camera facing
	w := keyToInt(window.GetKey(glfw.KeyW))
	a := keyToInt(window.GetKey(glfw.KeyA))
	s := keyToInt(window.GetKey(glfw.KeyS))
	d := keyToInt(window.GetKey(glfw.KeyD))

	ahead := w - s
	right := d - a
	scene.Camera.Position = scene.Camera.Position.
		Add(cameraRightV.Mul(moveSpeed * right * dt)).
		Add(cameraFrontV.Mul(moveSpeed * ahead * dt))

	// Construct final transform
	up := mgl.Vec3{0, 1, 0}
	lookAt := scene.Camera.Position.Add(cameraFrontV)
	scene.Camera.Transform = mgl.LookAtV(scene.Camera.Position, lookAt, up)

	return !exit
}

func keyToInt(value glfw.Action) float32 {
	if value == glfw.Press {
		return 1
	} else {
		return 0
	}
}
