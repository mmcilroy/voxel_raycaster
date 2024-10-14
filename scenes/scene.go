package scene

import (
	"fmt"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mmcilroy/voxel_raycaster/voxel"
)

const RESOLUTION_X, RESOLUTION_Y = 1600, 900

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

func DrawVoxel(x, y, z int32, s float32, c rl.Color) {
	hs := s / 2
	rl.DrawCube(rl.NewVector3(float32(x)*s+hs, float32(y)*s+hs, float32(z)*s+hs), s, s, s, c)
}

func DrawVoxelOutline(x, y, z int32, s float32, c rl.Color) {
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
	rl.SetTargetFPS(60)

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
		for z := int32(0); z < voxels.NumVoxelsZ-1; z++ {
			for y := int32(0); y < voxels.NumVoxelsY-1; y++ {
				for x := int32(0); x < voxels.NumVoxelsX-1; x++ {
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

// would this be better named VoxelColorFn?
// voxels are fixed color but pixel color is affected by lighting
type PixelColorFn func(hit int32, mapPos voxel.Vector3i) rl.Color

func raycastPixel(scene *RaycastingScene, x, y int32, pixelColorFn PixelColorFn) rl.Color {
	// get the ray direction
	plane := scene.Camera.Plane()
	_, rayDir := scene.Camera.RayDir(&plane, int32(x), int32(y))

	// decide what version of the voxel grid to use
	voxels := scene.Voxels
	if !scene.EnableRecursiveDDA {
		voxels = scene.UncompressedVoxels
	}

	// fire a ray into the scene and check what we hit
	hit, hitPos, mapPos := voxels.RaycastRecursive(scene.Camera.Body.Position, rayDir)

	// get the pixel color for the voxel and face
	color := pixelColorFn(hit, mapPos)

	// if lightning is enabled and something was hit apply shadows
	lightScale := float32(1)

	// check if the hit point is visible to the sun
	if scene.EnableLighting && hit != 0 && hit != 4 {

		// if we are not lighting per pixel do it per voxel face
		if !scene.EnablePerPixelLighting {
			hitPos = voxel.HitFaceCenter(hit, hitPos, mapPos, scene.UncompressedVoxels.VoxelSize)
		}
		sunHit, sunHitPos, _ := voxels.RaycastRecursive(scene.SunPos, voxel.Direction(hitPos, scene.SunPos))

		// if sun ray hits the same block and face as our initial ray calc lighting
		if sunHit == hit && voxel.Distance(hitPos, sunHitPos) < 0.5 /*sunMapPos.Equals(mapPos)*/ {
			lightScale = voxel.DiffuseLight(sunHit, voxel.Direction(scene.SunPos, sunHitPos))
		} else {
			lightScale = 0.5 // shadow
		}
	}

	return rl.NewColor(
		uint8(float32(color.R)*lightScale),
		uint8(float32(color.G)*lightScale),
		uint8(float32(color.B)*lightScale),
		255)
}

func raycastQuad(scene *RaycastingScene, xa, xb, ya, yb int32, pixelColorFn PixelColorFn, pixels *[]rl.Color, frameWait *sync.WaitGroup) {
	defer frameWait.Done()

	for y := ya; y < yb; y++ {
		for x := xa; x < xb; x++ {
			// get pixel color
			color := raycastPixel(scene, x, y, pixelColorFn)

			// write into pixel buffer
			(*pixels)[x+y*int32(scene.Camera.Resolution.X)] = color
		}
	}
}

func renderSoftware(scene *RaycastingScene, frame *rl.RenderTexture2D, pixelColorFn PixelColorFn, pixels *[]rl.Color, frameWait *sync.WaitGroup) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	// spread rays across threads
	frameWait.Add(NUM_THREADS_X * NUM_THREADS_Y)
	for ty := 0; ty < NUM_THREADS_Y; ty++ {
		for tx := 0; tx < NUM_THREADS_X; tx++ {
			go raycastQuad(
				scene,
				int32(int(scene.Camera.Resolution.X)/NUM_THREADS_X*tx),
				int32(int(scene.Camera.Resolution.X)/NUM_THREADS_X*(tx+1)),
				int32(int(scene.Camera.Resolution.Y)/NUM_THREADS_Y*ty),
				int32(int(scene.Camera.Resolution.Y)/NUM_THREADS_Y*(ty+1)),
				pixelColorFn,
				pixels,
				frameWait)
		}
	}
	frameWait.Wait()

	// center pixel is red for debugging
	cx, cy := int32(scene.Camera.Resolution.X/2), int32(scene.Camera.Resolution.Y/2)
	(*pixels)[cx+int32(cy*scene.Camera.Resolution.X)] = rl.Red

	// use output color to create frame
	rl.BeginTextureMode(*frame)
	for ry := 0; ry < int(scene.Camera.Resolution.Y); ry++ {
		for rx := 0; rx < int(scene.Camera.Resolution.X); rx++ {
			rl.DrawPixel(int32(rx), int32(ry), (*pixels)[rx+ry*int(scene.Camera.Resolution.X)])
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
}

type RaycastingScene struct {
	UncompressedVoxels     *voxel.VoxelGrid
	Voxels                 *voxel.VoxelGrid
	Camera                 voxel.Camera
	SunPos                 voxel.Vector3f
	EnableRecursiveDDA     bool
	EnableLighting         bool
	EnablePerPixelLighting bool
}

func RenderRaycastingScene(scene *RaycastingScene, pixelColorFn PixelColorFn, preFn func(), postFn func()) {
	// compress voxels
	for scene.Voxels.NumVoxelsY > 2 {
		scene.Voxels = scene.Voxels.Compress()
	}

	// also keep uncompressed handy
	scene.UncompressedVoxels = scene.Voxels
	for scene.UncompressedVoxels.Parent != nil {
		scene.UncompressedVoxels = scene.UncompressedVoxels.Parent
	}

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(RESOLUTION_X, RESOLUTION_Y, "")
	defer rl.CloseWindow()

	rl.DisableCursor()

	// indicates when frame is complete
	var frameWait sync.WaitGroup

	// the color for each pixel (cpu only)
	pixels := make([]rl.Color, int(scene.Camera.Resolution.X*scene.Camera.Resolution.Y))

	// the frame that will be displayed (cpu only)
	texture := rl.LoadRenderTexture(int32(scene.Camera.Resolution.X), int32(scene.Camera.Resolution.Y))

	for !rl.WindowShouldClose() {
		// character controls
		var moveForward, moveSide, moveUp float32

		speed := 5 * rl.GetFrameTime() * 1
		if rl.IsKeyDown(rl.KeyLeftShift) {
			speed *= 10
		}

		if rl.IsKeyDown('W') {
			moveForward += speed
		}

		if rl.IsKeyDown('S') {
			moveForward -= speed
		}

		if rl.IsKeyDown('A') {
			moveSide -= speed
		}

		if rl.IsKeyDown('D') {
			moveSide += speed
		}

		if rl.IsKeyDown(rl.KeySpace) {
			moveUp += speed
		}

		if rl.IsKeyDown(rl.KeyLeftControl) {
			moveUp -= speed
		}

		//  rendering controls
		if rl.IsKeyPressed('L') {
			scene.EnableLighting = !scene.EnableLighting
		}

		if rl.IsKeyPressed('R') {
			scene.EnableRecursiveDDA = !scene.EnableRecursiveDDA
		}

		if rl.IsKeyPressed('P') {
			scene.EnablePerPixelLighting = !scene.EnablePerPixelLighting
		}

		if rl.IsKeyDown(rl.KeyUp) {
			scene.SunPos.Y += speed
		}

		if rl.IsKeyDown(rl.KeyDown) {
			scene.SunPos.Y -= speed
		}

		if rl.IsKeyDown(rl.KeyLeft) {
			scene.SunPos.X += speed
		}

		if rl.IsKeyDown(rl.KeyRight) {
			scene.SunPos.X -= speed
		}

		// allows us to debug raycasting for the center pixel
		if rl.IsKeyPressed('T') {
			cx, cy := int32(scene.Camera.Resolution.X/2), int32(scene.Camera.Resolution.Y/2)
			raycastPixel(scene, cx, cy, pixelColorFn)
		}

		// move and rotate camera
		mouseDelta := rl.GetMouseDelta()
		scene.Camera.Body.Rotate(mouseDelta.X*-0.003, mouseDelta.Y*-0.003)

		//prevPos := scene.Camera.Body.Position
		scene.Camera.Body.Move(moveForward, moveSide, moveUp)

		// move to new position as long as there isn't a voxel there
		//if scene.UncompressedVoxels.RectangleIntersects(scene.Camera.Body.Position, 1, 2) {
		//	scene.Camera.Body.Position = prevPos
		//}

		preFn()

		renderSoftware(scene, &texture, pixelColorFn, &pixels, &frameWait)

		rl.DrawFPS(20, 20)
		rl.DrawText(fmt.Sprintf("%.02f, %.02f, %.02f, %.02f, %.02f", scene.Camera.Body.Position.X, scene.Camera.Body.Position.Y, scene.Camera.Body.Position.Z, scene.Camera.Body.Rotation.X, scene.Camera.Body.Rotation.Y), 20, 40, 20, rl.White)
		rl.DrawText(fmt.Sprintf("Lighting (L): %t, RecursiveDDA (R): %t, PerPixelLighting (P): %t", scene.EnableLighting, scene.EnableRecursiveDDA, scene.EnablePerPixelLighting), 20, 60, 20, rl.White)

		postFn()

		rl.EndDrawing()
	}
}
