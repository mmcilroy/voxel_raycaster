package scene

import (
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmcilroy/structure_go/voxel"
)

const NUM_THREADS_X = 4
const NUM_THREADS_Y = 4

func ToRlVector(v voxel.Vector3f) rl.Vector3 {
	return rl.NewVector3(v.X, v.Y, v.Z)
}

func DrawRay(pos voxel.Vector3f, dir voxel.Vector3f, color rl.Color) {
	rl.DrawRay(rl.NewRay(ToRlVector(pos), ToRlVector(dir)), rl.SkyBlue)
}

func DrawSphere(pos voxel.Vector3f, radius float32, color rl.Color) {
	rl.DrawSphere(ToRlVector(pos), radius, color)
}

func DrawVoxel(x, y, z int, s float32, c rl.Color) {
	hs := s / 2
	rl.DrawCube(rl.NewVector3(float32(x)*s+hs, float32(y)*s+hs, float32(z)*s+hs), s, s, s, c)
}

func DrawVoxelOutline(x, y, z int, s float32, c rl.Color) {
	hs := s / 2
	rl.DrawCubeWires(rl.NewVector3(float32(x)*s+hs, float32(y)*s+hs, float32(z)*s+hs), s, s, s, c)
}

func RotatingPosition(origin voxel.Vector3f, radius, angleX, angleY float32) voxel.Vector3f {
	up := voxel.Vector3f{X: 0, Y: 1, Z: 0}
	pos := voxel.Vector3f{X: 0, Y: 0, Z: radius}
	pos = pos.RotateByAxisAngle(up, angleX)
	forward := pos.RotateByAxisAngle(up, angleX).Normalize()
	right := forward.CrossProduct(up)
	pos = pos.RotateByAxisAngle(right, angleY)
	pos = pos.Plus(origin)
	return pos
}

func RenderScene(handleInput func(), render3D func(), render2D func()) {
	rl.InitWindow(1600, 900, "")
	defer rl.CloseWindow()

	rl.DisableCursor()

	camera := rl.Camera{
		Position:   rl.Vector3{X: 0, Y: 10, Z: -10},
		Target:     rl.Vector3{X: 4, Y: 4, Z: 4},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       60,
		Projection: rl.CameraPerspective,
	}

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeySpace) {
			camera.Position.Y += 10 * rl.GetFrameTime()
		}
		if rl.IsKeyDown(rl.KeyLeftControl) {
			camera.Position.Y -= 10 * rl.GetFrameTime()
		}
		handleInput()
		rl.UpdateCamera(&camera, rl.CameraThirdPerson)
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		rl.BeginMode3D(camera)
		render3D()
		rl.EndMode3D()
		render2D()
		rl.EndDrawing()
	}
}

func RenderVoxelScene(voxels *voxel.VoxelGrid, handleInput func(), render3D func(), render2D func()) {
	halfSize := voxels.VoxelSize / 2
	RenderScene(handleInput, func() {
		for z := 0; z < voxels.NumVoxelsZ-1; z++ {
			for y := 0; y < voxels.NumVoxelsY-1; y++ {
				for x := 0; x < voxels.NumVoxelsX-1; x++ {
					if voxels.GetVoxel(x, y, z) {
						rl.DrawCube(rl.NewVector3(
							voxels.VoxelSize*float32(x)+halfSize, voxels.VoxelSize*float32(y)+halfSize, voxels.VoxelSize*float32(z)+halfSize),
							voxels.VoxelSize, voxels.VoxelSize, voxels.VoxelSize,
							rl.NewColor(255, 0, 0, 127))
					}
				}
			}
		}
		rl.DrawGrid(128, 1)
		render3D()
	}, render2D)
}

type PixelColorFn func(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, rayDir voxel.Vector3f) rl.Color

func raycast(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, xa, xb, ya, yb int32, pixelColorFn PixelColorFn, pixels *[]rl.Color, frameWait *sync.WaitGroup) {
	defer frameWait.Done()

	for y := ya; y < yb; y++ {
		for x := xa; x < xb; x++ {
			// get the ray direction
			_, rayDir := camera.GetRayForPixel(int32(x), int32(y))

			// get the pixel color
			(*pixels)[x+y*int32(camera.Resolution.X)] = pixelColorFn(camera, voxels, rayDir)
		}
	}
}

func RenderRaycastingScene(camera *voxel.RaycastingCamera, voxels *voxel.VoxelGrid, pixelColorFn PixelColorFn, preFn func(), postFn func()) {
	rl.InitWindow(1600, 900, "")
	defer rl.CloseWindow()

	rl.DisableCursor()

	// indicates when frame is complete
	var frameWait sync.WaitGroup

	// the frame that will be displayed
	frame := rl.LoadRenderTexture(int32(camera.Resolution.X), int32(camera.Resolution.Y))

	// the color for each pixel
	pixels := make([]rl.Color, int(camera.Resolution.X*camera.Resolution.Y))

	// the direction the camera is facing
	rotationX, rotationY := float32(-5.5), float32(-0.5)

	for !rl.WindowShouldClose() {
		distanceMoved := 1.3 * rl.GetFrameTime()

		if rl.IsKeyDown(rl.KeyLeftShift) {
			distanceMoved *= 5
		}

		if rl.IsKeyDown('W') {
			camera.Position = camera.Position.Plus(camera.Forward.MulScalar(distanceMoved))
		}

		if rl.IsKeyDown('A') {
			camera.Position = camera.Position.Sub(camera.Right.MulScalar(distanceMoved))
		}

		if rl.IsKeyDown('S') {
			camera.Position = camera.Position.Sub(camera.Forward.MulScalar(distanceMoved))
		}

		if rl.IsKeyDown('D') {
			camera.Position = camera.Position.Plus(camera.Right.MulScalar(distanceMoved))
		}

		if rl.IsKeyDown(rl.KeySpace) {
			camera.Position.Y += distanceMoved
		}

		if rl.IsKeyDown(rl.KeyLeftControl) {
			camera.Position.Y -= distanceMoved
		}

		mouseDelta := rl.GetMouseDelta()
		rotationX += mouseDelta.X * -0.003
		rotationY += mouseDelta.Y * -0.003
		camera.Rotate(rotationX, rotationY)

		preFn()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// spread rays across threads
		frameWait.Add(NUM_THREADS_X * NUM_THREADS_Y)
		for ty := 0; ty < NUM_THREADS_Y; ty++ {
			for tx := 0; tx < NUM_THREADS_X; tx++ {
				go raycast(
					camera,
					voxels,
					int32(int(camera.Resolution.X)/NUM_THREADS_X*tx),
					int32(int(camera.Resolution.X)/NUM_THREADS_X*(tx+1)),
					int32(int(camera.Resolution.Y)/NUM_THREADS_Y*ty),
					int32(int(camera.Resolution.Y)/NUM_THREADS_Y*(ty+1)),
					pixelColorFn,
					&pixels,
					&frameWait)
			}
		}
		frameWait.Wait()

		// use output color to create frame
		rl.BeginTextureMode(frame)
		for ry := 0; ry < int(camera.Resolution.Y); ry++ {
			for rx := 0; rx < int(camera.Resolution.X); rx++ {
				rl.DrawPixel(int32(camera.Resolution.X)-int32(rx), int32(camera.Resolution.Y)-int32(ry), pixels[rx+ry*int(camera.Resolution.X)])
			}
		}
		rl.EndTextureMode()

		// scale frame to window and draw it
		rl.DrawTexturePro(frame.Texture,
			rl.NewRectangle(0, 0, float32(frame.Texture.Width), -float32(frame.Texture.Height)),
			rl.NewRectangle(0, 0, 1600, 900),
			rl.NewVector2(0, 0),
			0,
			rl.White)

		rl.DrawFPS(20, 20)
		rl.DrawText(fmt.Sprintf("%.02f, %.02f, %.02f", camera.Position.X, camera.Position.Y, camera.Position.Z), 20, 40, 20, rl.White)
		rl.DrawText(fmt.Sprintf("%.02f, %.02f", rotationX, rotationY), 20, 60, 20, rl.White)

		postFn()

		rl.EndDrawing()
	}
}
