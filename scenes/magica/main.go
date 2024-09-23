package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mattkimber/gandalf/magica"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

const WORLD_WIDTH, WORLD_HEIGHT = 128, 128

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

var world *voxel.VoxelGrid

var palette = make([]rl.Color, 256)

var voxelColors = map[voxel.Vector3i]byte{}

//var voxelColors2 = make([]byte, WORLD_WIDTH*WORLD_HEIGHT*WORLD_WIDTH)

func initWorld() {
	world = voxel.NewVoxelGrid(WORLD_WIDTH, WORLD_HEIGHT, WORLD_WIDTH, 1.0)

	object, _ := magica.FromFile("..\\..\\assets\\monu3.vox")

	for i := 0; i < len(object.PaletteData); i += 4 {
		palette[i/4] = rl.NewColor(object.PaletteData[i], object.PaletteData[i+1], object.PaletteData[i+2], 255)
	}

	for z := 0; z < object.Size.Z; z++ {
		for y := 0; y < object.Size.Y; y++ {
			for x := 0; x < object.Size.X; x++ {
				v := object.Voxels[x][y][z]
				if v != 0 {
					world.SetVoxel(int32(x+10), int32(z), int32(y+10), true)
					voxelColors[voxel.Vector3i{X: int32(x + 10), Y: int32(z), Z: int32(y + 10)}] = v

					//colorIndex := (x + 10) + z*WORLD_WIDTH + (y+10)*WORLD_WIDTH*WORLD_HEIGHT
					//voxelColors2[colorIndex] = v
				}
			}
		}
	}

	for world.NumVoxelsY > 2 {
		world = world.Compress()
	}
}

func pixelColorFn(hit int32, mapPos voxel.Vector3i) rl.Color {
	color := rl.SkyBlue

	if hit != 0 {
		//colorIndex := (mapPos.X) + mapPos.Z*WORLD_WIDTH + (mapPos.Y)*WORLD_WIDTH*WORLD_HEIGHT
		//paletteIndex := int(voxelColors2[colorIndex])

		paletteIndex := int(voxelColors[voxel.Vector3i{X: mapPos.X, Y: mapPos.Y, Z: mapPos.Z}])
		if paletteIndex > 0 {
			paletteIndex -= 1
		}
		color = palette[paletteIndex]
	}

	return color
}

func main() {
	initWorld()

	raycastingScene := scene.RaycastingScene{
		Voxels:                 world,
		Camera:                 voxel.NewCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66),
		SunPos:                 voxel.Vector3f{X: WORLD_WIDTH - 1, Y: WORLD_HEIGHT - 1, Z: 0},
		EnableRecursiveDDA:     true,
		EnableLighting:         true,
		EnablePerPixelLighting: true,
	}
	raycastingScene.Camera.Body.Position.X = 1
	raycastingScene.Camera.Body.Position.Y = 5
	raycastingScene.Camera.Body.Position.Z = 1

	scene.RenderRaycastingScene(&raycastingScene, pixelColorFn, func() {}, func() {})
}
