# voxel_raycaster

Adventures in rendering voxels with raycasting

```
go run scenes\perlin-world\main.go
```

![Tower](gifs/tower.gif)

![Perlin](gifs/perlin.gif)


Scene {
    Voxels[]        // multiple resolutions
    AmbientLight
    DynamicLights
}

type VoxelGrid interface {
	Set(x, y, z int32, b bool)
	Get(x, y, z int32) bool
	Count() Vector3i
	Size() float32
	Compress() *VoxelGrid
}

type TraceOptions {
	maxSteps uint32
	callback func(VoxelGrid*, Vector3i)
}

type TraceResult struct {	// renamed as voxel intersection perhaps
	hit      int8			// miss / oob / side
	hitPos   Vector3f
	mapPos   Vector3i		// should be voxel not mapPos
	numSteps int32			// should not be part of this
}

type Tracer interface {
	Trace(ray Ray, options TraceOptions) TraceResult
}
