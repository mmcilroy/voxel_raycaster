package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mattkimber/gandalf/magica"
	scene "github.com/mmcilroy/structure_go/scenes"
	"github.com/mmcilroy/structure_go/voxel"
)

const WORLD_WIDTH, WORLD_HEIGHT = 256, 256

const NUM_RAYS_X, NUM_RAYS_Y = 320, 180

var world *voxel.VoxelGrid

var palette = make([]rl.Color, 256)

var voxelColors = map[voxel.Vector3i]byte{}

func initWorld() {
	world = voxel.NewVoxelGrid(WORLD_WIDTH, WORLD_HEIGHT, WORLD_WIDTH, 1)

	object, _ := magica.FromFile("monu3.vox")
	//object, _ := magica.FromFile("monu1.vox")
	//object, _ := magica.FromFile("chr_man.vox")
	//object, _ := magica.FromFile("chr_knight.vox")

	for i := 0; i < len(object.PaletteData); i += 4 {
		palette[i/4] = rl.NewColor(object.PaletteData[i], object.PaletteData[i+1], object.PaletteData[i+2], 255)
	}

	for z := 0; z < object.Size.Z; z++ {
		for y := 0; y < object.Size.Y; y++ {
			for x := 0; x < object.Size.X; x++ {
				v := object.Voxels[x][y][z]
				if v != 0 {
					world.SetVoxel(x+10, z, y+10, true)
					voxelColors[voxel.Vector3i{X: x + 10, Y: z, Z: y + 10}] = v
				}
			}
		}
	}

	for world.NumVoxelsY > 2 {
		world = world.Compress()
	}
}

func pixelColorFn(hit int, mapPos voxel.Vector3i) rl.Color {
	color := rl.SkyBlue

	if hit != 0 {
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
		Camera:                 voxel.NewRaycastingCamera(NUM_RAYS_X, NUM_RAYS_Y, 0.66),
		SunPos:                 voxel.Vector3f{X: WORLD_WIDTH - 1, Y: WORLD_HEIGHT - 1, Z: 0},
		EnableRecursiveDDA:     true,
		EnableLighting:         true,
		EnablePerPixelLighting: true,
	}

	scene.RenderRaycastingScene(&raycastingScene, pixelColorFn, func() {}, func() {})
}
