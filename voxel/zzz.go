package voxel

type Voxels interface {
	Set(x, y, z int32, b bool)
	Get(x, y, z int32) bool
	Size() float32
	Count() Vector3i
	Compress() *Voxels
}

type TestVoxels struct {
	voxels []uint8
	size   float32
	count  Vector3i
}

func NewTestVoxels(nx, ny, nz int32, sz float32) Voxels {
	return &TestVoxels{
		voxels: make([]uint8, nx*ny*nz),
		size:   sz,
		count:  Vector3i{X: nx, Y: ny, Z: nz},
	}
}

func (voxels *TestVoxels) index(x, y, z int32) int32 {
	return x + z*voxels.count.X + y*voxels.count.X*voxels.count.Z
}

func (voxels *TestVoxels) Set(x, y, z int32, b bool) {
	var v uint8
	if b {
		v = 1
	}
	i := voxels.index(x, y, z)
	voxels.voxels[i] = v
}

func (voxels *TestVoxels) Get(x, y, z int32) bool {
	i := voxels.index(x, y, z)
	return voxels.voxels[i] != 0
}

func (voxels *TestVoxels) Size() float32 {
	return voxels.size
}

func (voxels *TestVoxels) Count() Vector3i {
	return voxels.count
}

func (voxels *TestVoxels) Compress() *Voxels {
	newVoxels := NewTestVoxels(voxels.count.X/2, voxels.count.Y/2, voxels.count.Z/2, voxels.size*2)

	for y := int32(0); y < newVoxels.Count().Y; y++ {
		for z := int32(0); z < newVoxels.Count().Z; z++ {
			for x := int32(0); x < newVoxels.Count().X; x++ {

				px, py, pz := x*2, y*2, z*2

				if voxels.Get(px, py, pz) ||
					voxels.Get(px+1, py, pz) ||
					voxels.Get(px, py+1, pz) ||
					voxels.Get(px+1, py+1, pz) ||
					voxels.Get(px, py, pz+1) ||
					voxels.Get(px+1, py, pz+1) ||
					voxels.Get(px, py+1, pz+1) ||
					voxels.Get(px+1, py+1, pz+1) {
					newVoxels.Set(x, y, z, true)
				}
			}
		}
	}

	return &newVoxels
}

type TraceCallback func(voxels *Voxels, mapPos Vector3i)

type TraceParams struct {
	RayStart Vector3f
	RayDir   Vector3f
	MaxSteps int32
	Callback TraceCallback
}

type TraceResult struct {
	Hit      bool
	OOB      bool
	Side     int32
	HitPos   Vector3f
	MapPos   Vector3i
	NumSteps int32
}

type Tracer interface {
	Trace(TraceParams) TraceResult
}

type TracerImpl struct {
	voxels *Voxels
}

func (tracer *TracerImpl) Trace(params TraceParams) TraceResult {
	return Trace(tracer.voxels, params)
}

type MipmapTracerImpl struct {
	Voxels []*Voxels
}

func (tracer *MipmapTracerImpl) Trace(params TraceParams) TraceResult {

	resolution := len(tracer.Voxels) - 1
	params.MaxSteps = 4
	numSteps := int32(0)

	for {
		// check for hit at current resolution
		result := Trace(tracer.Voxels[resolution], params)
		params.RayStart = result.HitPos
		numSteps += result.NumSteps

		// hit at highest resolution / oob so stop
		if result.OOB || (result.Hit && resolution == 0) {
			result.NumSteps = numSteps
			return result
		}

		// hit at low resolution so switch to higher res and continue
		if result.Hit && resolution > 0 {
			resolution--
			continue
		}

		// missed so switch to lower res and continue
		//if resolution < len(tracer.Voxels)-1 {
		//	resolution++
		//	continue
		//}
	}
}

func isOutside(voxels *Voxels, mapPos Vector3i) bool {
	// returns true if mapPos is outside the grid
	numVoxels := (*voxels).Count()
	return mapPos.X < 0 ||
		mapPos.Y < 0 ||
		mapPos.Z < 0 ||
		mapPos.X >= numVoxels.X ||
		mapPos.Y >= numVoxels.Y ||
		mapPos.Z >= numVoxels.Z
}

func isOutOfBounds(voxels *Voxels, mapPos Vector3i, rayDir Vector3f) bool {
	// returns true if mapPos is outside grid and we are heading away from it
	numVoxels := (*voxels).Count()
	return (mapPos.X < 0 && rayDir.X <= 0) ||
		(mapPos.Y < 0 && rayDir.Y <= 0) ||
		(mapPos.Z < 0 && rayDir.Z <= 0) ||
		(mapPos.X >= numVoxels.X && rayDir.X >= 0) ||
		(mapPos.Y >= numVoxels.Y && rayDir.Y >= 0) ||
		(mapPos.Z >= numVoxels.Z && rayDir.Z >= 0)
}

func checkVoxel(voxels *Voxels, mapPos Vector3i, rayDir Vector3f) (bool, bool) {
	// check if the voxel is outside the grid
	if isOutside(voxels, mapPos) {
		if isOutOfBounds(voxels, mapPos, rayDir) {
			// no hit, oob
			return false, true
		} else {
			// no hit, not oob
			return false, false
		}
	}

	// check if the current voxel is empty
	present := (*voxels).Get(mapPos.X, mapPos.Y, mapPos.Z)

	// possibly hit, not oob
	return present, false
}

func Trace(voxels *Voxels, params TraceParams) TraceResult {
	var result TraceResult

	// convert rayPos to voxel space
	rayPos := params.RayStart.DivScalar((*voxels).Size())

	// which box of the map we're in
	mapPos := rayPos.ToVector3i()

	// length of ray from one xyz side to next
	deltaDist := params.RayDir.Inverse().Abs()

	// what direction to step in x or y-direction (either +1 or -1)
	step := params.RayDir.Sign().ToVector3i()

	// length of ray from current position to next x or y-side
	sideDist := calcSideDist(rayPos, params.RayDir, deltaDist, mapPos)

	// how far the ray has travelled
	dist := float32(0.0)

	// loop until we hit something or we go oob
	for {

		// make callback if one is provided
		if params.Callback != nil {
			params.Callback(voxels, mapPos)
		}

		// check current voxel
		result.Hit, result.OOB = checkVoxel(voxels, mapPos, params.RayDir)
		result.NumSteps++
		if result.Hit || result.OOB {
			break
		}

		// jump to next map square, either in x, y or z direction
		if sideDist.X <= sideDist.Y && sideDist.X <= sideDist.Z {
			dist = sideDist.X
			sideDist.X += deltaDist.X
			mapPos.X += step.X
			result.Side = 1 * step.X
		} else if sideDist.Y <= sideDist.X && sideDist.Y <= sideDist.Z {
			dist = sideDist.Y
			sideDist.Y += deltaDist.Y
			mapPos.Y += step.Y
			result.Side = 2 * step.X
		} else {
			dist = sideDist.Z
			sideDist.Z += deltaDist.Z
			mapPos.Z += step.Z
			result.Side = 3 * step.X
		}

		// stop if we hit max steps
		if params.MaxSteps > 0 && result.NumSteps >= params.MaxSteps {
			break
		}
	}

	// calculate the hit point
	result.MapPos = mapPos
	result.HitPos = rayPos.Plus(params.RayDir.MulScalar(dist)).MulScalar((*voxels).Size())

	// snap to grid to prevent rounding errors
	if result.Side == 1 || result.Side == -1 {
		result.HitPos = result.HitPos.RoundX()
	} else if result.Side == 2 || result.Side == -2 {
		result.HitPos = result.HitPos.RoundY()
	} else if result.Side == 3 || result.Side == -3 {
		result.HitPos = result.HitPos.RoundZ()
	}

	return result
}
