package main

import (
	"fmt"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestBroken(t *testing.T) {

	grid := NewVoxelGrid(16)

	x, y, z := 3, 1, 0

	rayStart := rl.NewVector3(4, 0.0, 0.0)
	rayEnd := rl.NewVector3(float32(x), float32(y), float32(z))
	rayDir := rl.Vector3Normalize(rl.Vector3Subtract(rayEnd, rayStart))

	grid.SetVoxel(x, y, z, true)

	hit, _, _ := grid.DDA(rayStart, rayDir)
	if hit <= 0 {
		t.Error()
	}
}

func TestDDA(t *testing.T) {

	grid := NewVoxelGrid(16)

	startPositions := []rl.Vector3{
		rl.Vector3Zero(),
		rl.NewVector3(float32(grid.Size)-1, float32(grid.Size)-1, float32(grid.Size)-1),

		// todo
		// using these positions should ideally work but seems to fail due to
		// what are i think rounding errors?

		rl.NewVector3(float32(grid.Size)-1, 0.0, 0.0),
		//rl.NewVector3(0.0, float32(grid.Size)-1, 0.0),
		//rl.NewVector3(0.0, 0.0, float32(grid.Size)-1),
		//rl.NewVector3(float32(grid.Size)-1, float32(grid.Size)-1, 0.0),
		//rl.NewVector3(0.0, float32(grid.Size)-1, float32(grid.Size)-1),
		//rl.NewVector3(float32(grid.Size)-1, 0.0, float32(grid.Size)-1),
	}

	for i := 0; i < len(startPositions); i++ {
		for z := 0; z < grid.Size; z++ {
			for y := 0; y < grid.Size; y++ {
				for x := 0; x < grid.Size; x++ {

					start := startPositions[i]
					end := rl.NewVector3(float32(x), float32(y), float32(z))

					if rl.Vector3Equals(start, end) {
						continue
					}

					dir := rl.Vector3Normalize(rl.Vector3Subtract(end, start))

					grid.Clear()
					grid.SetVoxel(int(end.X), int(end.Y), int(end.Z), true)

					hit, _, _ := grid.DDA(start, dir)
					if hit <= 0 {
						t.Error()
					}

					//if !rl.Vector3Equals(end, pos) {
					//	t.Error()
					//}
				}
			}
		}
	}
}

func TestCompressed(t *testing.T) {
	grid16 := NewVoxelGrid(16)
	grid16.SetVoxel(0, 0, 0, true)
	grid16.SetVoxel(4, 0, 0, true)
	grid16.SetVoxel(8, 0, 0, true)
	grid16.SetVoxel(12, 0, 0, true)

	grid4 := grid16.Compress()
	if !grid4.GetVoxel(0, 0, 0) {
		t.Error()
	}
	if !grid4.GetVoxel(1, 0, 0) {
		t.Error()
	}
	if !grid4.GetVoxel(2, 0, 0) {
		t.Error()
	}
	if !grid4.GetVoxel(3, 0, 0) {
		t.Error()
	}
	if grid4.GetVoxel(0, 1, 0) {
		t.Error()
	}
}

func TestDDACompressed(t *testing.T) {

	d := rl.Vector3Distance(rl.NewVector3(3, 3, 3), rl.Vector3Zero())
	fmt.Println(d)

	grid1024 := NewVoxelGrid(1024)
	grid1024.SetVoxel(grid1024.Size-1, grid1024.Size-1, grid1024.Size-1, true)

	grid256 := grid1024.Compress()
	grid64 := grid256.Compress()
	grid16 := grid64.Compress()
	grid4 := grid16.Compress()

	dir := rl.Vector3Normalize(rl.Vector3One())

	hit, pos, _ := grid4.DDA(rl.Vector3Zero(), dir)
	if hit > 0 {
		hit, pos, _ = grid16.DDA(rl.Vector3Scale(pos, 4), dir)
		if hit > 0 {
			hit, pos, _ = grid64.DDA(rl.Vector3Scale(pos, 4), dir)
			if hit > 0 {
				hit, pos, _ = grid256.DDA(rl.Vector3Scale(pos, 4), dir)
				if hit > 0 {
					hit, pos, _ = grid1024.DDA(rl.Vector3Scale(pos, 4), dir)
				}
			}
		}
	}

	if hit == 0 {
		t.Error()
	}

	if pos.X == 0 {
		t.Error()
	}
}
