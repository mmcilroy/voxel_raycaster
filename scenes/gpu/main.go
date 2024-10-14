package main

// go build -tags opengl43

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmcilroy/voxel_raycaster/voxel"
)

const RESOLUTION_X, RESOLUTION_Y = 1600, 900

const WORLD_SIZE = 256

func initWorld() *voxel.VoxelGrid {
	voxels := voxel.NewVoxelGrid(WORLD_SIZE, WORLD_SIZE, WORLD_SIZE, 1.0)
	perlinNoise := rl.GenImagePerlinNoise(WORLD_SIZE, WORLD_SIZE, 0, 0, 0.5)
	colors := rl.LoadImageColors(perlinNoise)

	for z := int32(0); z < WORLD_SIZE; z++ {
		for x := int32(0); x < WORLD_SIZE; x++ {
			color := colors[x+z*WORLD_SIZE]
			height := float32(color.R) / 255.0 * float32(WORLD_SIZE)
			for y := int32(0); y < int32(height)+1; y++ {
				voxels.SetVoxel(x, y, z, true)
			}
		}
	}

	return voxels
}

func main() {
	v := initWorld()
	s := voxel.Vector3f{X: float32(v.NumVoxelsX - 1), Y: float32(v.NumVoxelsY) - 1, Z: 0}
	c := voxel.NewCamera(RESOLUTION_X, RESOLUTION_Y, 0.66)
	c.Body.Position = voxel.Vector3f{X: 1, Y: WORLD_SIZE - 1, Z: 1}

	rl.InitWindow(RESOLUTION_X, RESOLUTION_Y, "")
	defer rl.CloseWindow()

	rl.DisableCursor()

	imBlank := rl.GenImageColor(RESOLUTION_X, RESOLUTION_Y, rl.Blank)
	texture := rl.LoadTextureFromImage(imBlank)
	rl.UnloadImage(imBlank)

	shader := rl.LoadShader("", "..\\..\\assets\\shaders\\frag.fs")

	resolutionLoc := rl.GetShaderLocation(shader, "resolution")
	cameraPosLoc := rl.GetShaderLocation(shader, "cameraPos")
	cameraPlaneLoc := rl.GetShaderLocation(shader, "cameraPlane")
	cameraUpLoc := rl.GetShaderLocation(shader, "cameraUp")
	cameraRightLoc := rl.GetShaderLocation(shader, "cameraRight")
	numVoxelsLoc := rl.GetShaderLocation(shader, "numVoxels")
	sunPosLoc := rl.GetShaderLocation(shader, "sunPos")

	rl.SetShaderValue(shader, resolutionLoc, []float32{RESOLUTION_X, RESOLUTION_Y}, rl.ShaderUniformVec2)
	rl.SetShaderValue(shader, numVoxelsLoc, []float32{float32(v.NumVoxelsX), float32(v.NumVoxelsY), float32(v.NumVoxelsZ)}, rl.ShaderUniformVec3)

	ssbo := rl.LoadShaderBuffer(uint32(len(v.Voxels)), unsafe.Pointer(unsafe.SliceData(v.Voxels)), rl.DynamicCopy)
	rl.BindShaderBuffer(ssbo, 13)

	for !rl.WindowShouldClose() {

		mouseDelta := rl.GetMouseDelta()
		rotateUp, rotateSide := mouseDelta.Y*-0.003, mouseDelta.X*-0.003

		speed := 10 * rl.GetFrameTime()
		moveForward, moveRight, moveUp := float32(0), float32(0), float32(0)

		if rl.IsKeyDown('W') {
			moveForward += speed
		}

		if rl.IsKeyDown('S') {
			moveForward -= speed
		}

		if rl.IsKeyDown('A') {
			moveRight -= speed
		}

		if rl.IsKeyDown('D') {
			moveRight += speed
		}

		if rl.IsKeyDown(rl.KeySpace) {
			moveUp += speed
		}

		if rl.IsKeyDown(rl.KeyLeftControl) {
			moveUp -= speed
		}

		c.Body.Rotate(rotateSide, rotateUp)
		c.Body.Move(moveForward, moveRight, moveUp)
		plane := c.Plane()

		rl.SetShaderValue(shader, cameraPosLoc, []float32{c.Body.Position.X, c.Body.Position.Y, c.Body.Position.Z}, rl.ShaderUniformVec3)
		rl.SetShaderValue(shader, cameraPlaneLoc, []float32{plane.CenterPos.X, plane.CenterPos.Y, plane.CenterPos.Z}, rl.ShaderUniformVec3)
		rl.SetShaderValue(shader, cameraUpLoc, []float32{plane.UpDir.X, plane.UpDir.Y, plane.UpDir.Z}, rl.ShaderUniformVec3)
		rl.SetShaderValue(shader, cameraRightLoc, []float32{plane.RightDir.X, plane.RightDir.Y, plane.RightDir.Z}, rl.ShaderUniformVec3)
		rl.SetShaderValue(shader, sunPosLoc, []float32{float32(s.X), float32(s.Y), float32(s.Z)}, rl.ShaderUniformVec3)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginShaderMode(shader)
		rl.DrawTexture(texture, 0, 0, rl.White)
		rl.EndShaderMode()
		rl.DrawFPS(20, 20)
		rl.EndDrawing()
	}

	rl.UnloadShader(shader)
}
