// Author: HyperparameterSI

// Package scene3d supplies an opt-in 3D viewport for Golet native UI.
// The initial renderer is a pure-Go software rasterizer, not a native GPU
// engine. It supports indexed meshes, six-plane clipping, depth testing,
// opaque flat materials and directional lighting. No frontend is required.
package scene3d

import (
	"fmt"
	"image/color"
	"math"
)

// Mesh owns immutable geometry, shared by any number of instances.
// NewMesh copies its input. CCW triangles face outwards.
type Mesh struct {
	vertices []Vec3
	indices  []uint32
	min, max Vec3
}

// NewMesh validates and copies indexed triangles. Degenerate triangles are
// legal and skipped during rendering. Vertices must be finite.
func NewMesh(vertices []Vec3, indices []uint32) (*Mesh, error) {
	m, err := ownedMesh(vertices, indices)
	if err != nil {
		return nil, err
	}
	m.vertices = append([]Vec3(nil), vertices...)
	m.indices = append([]uint32(nil), indices...)
	return m, nil
}

// ownedMesh takes exclusively owned slices, avoiding a second full model
// allocation when a loader has already constructed fresh geometry.
func ownedMesh(vertices []Vec3, indices []uint32) (*Mesh, error) {
	if len(vertices) == 0 || len(indices) == 0 || len(indices)%3 != 0 {
		return nil, fmt.Errorf("scene3d: expected vertices and triangle indices")
	}
	m := &Mesh{min: vertices[0], max: vertices[0]}
	for _, v := range vertices {
		if !v.valid() {
			return nil, fmt.Errorf("scene3d: non-finite vertex")
		}
		m.min = Vec3{min(m.min.X, v.X), min(m.min.Y, v.Y), min(m.min.Z, v.Z)}
		m.max = Vec3{max(m.max.X, v.X), max(m.max.Y, v.Y), max(m.max.Z, v.Z)}
	}
	for _, i := range indices {
		if uint64(i) >= uint64(len(vertices)) {
			return nil, fmt.Errorf("scene3d: vertex index %d out of range", i)
		}
	}
	m.vertices = vertices
	m.indices = indices
	return m, nil
}
func (m *Mesh) VertexCount() int   { return len(m.vertices) }
func (m *Mesh) TriangleCount() int { return len(m.indices) / 3 }

// Geometry returns independent copies for exporters or alternate renderers.
func (m *Mesh) Geometry() ([]Vec3, []uint32) {
	return append([]Vec3(nil), m.vertices...), append([]uint32(nil), m.indices...)
}

// Cube returns a unit cube centered on the origin. Retain and share it.
func Cube() *Mesh {
	m, _ := NewMesh([]Vec3{{-.5, -.5, -.5}, {.5, -.5, -.5}, {.5, .5, -.5}, {-.5, .5, -.5}, {-.5, -.5, .5}, {.5, -.5, .5}, {.5, .5, .5}, {-.5, .5, .5}}, []uint32{0, 2, 1, 0, 3, 2, 4, 5, 6, 4, 6, 7, 0, 4, 7, 0, 7, 3, 1, 2, 6, 1, 6, 5, 3, 7, 6, 3, 6, 2, 0, 1, 5, 0, 5, 4})
	return m
}

// Material is opaque flat shading. Alpha must be 255; blending is not yet
// supported. Unlit bypasses lighting; DoubleSided disables back-face culling.
type Material struct {
	Color              color.RGBA
	Unlit, DoubleSided bool
}

// Instance references shared geometry with a per-instance affine transform.
type Instance struct {
	Mesh      *Mesh
	Transform Mat4
	Material  Material
}

// DirectionalLight.Direction points FROM the surface TOWARD the light.
// Intensity is nonnegative. Lighting is flat Lambert plus Scene.Ambient.
type DirectionalLight struct {
	Direction Vec3
	Intensity float32
}

// Scene is a caller-owned snapshot. Do not mutate while rendering. The caller
// controls animation, simulation and scheduling; there is no background loop.
type Scene struct {
	Instances  []Instance
	Lights     []DirectionalLight
	Ambient    float32
	Background color.RGBA
}

// Camera is a perspective camera. FOVY is vertical radians; Near must be
// positive and Far greater than Near. Defaults are explicit via DefaultCamera.
type Camera struct {
	Position, Target, Up Vec3
	FOVY, Near, Far      float32
}

func DefaultCamera() Camera {
	return Camera{Position: Vec3{3, 2, 4}, Up: Vec3{0, 1, 0}, FOVY: float32(math.Pi / 3), Near: .1, Far: 100}
}
func (c Camera) validate() error {
	direction := c.Target.Sub(c.Position)
	right := direction.Cross(c.Up)
	length, side := direction.Dot(direction), right.Dot(right)
	if !c.Position.valid() || !c.Target.valid() || !c.Up.valid() || !finite(c.FOVY) || c.FOVY <= 0 || c.FOVY >= math.Pi || !finite(c.Near) || !finite(c.Far) || c.Near <= 0 || c.Far <= c.Near || !finite(length) || !finite(side) || length < 1e-12 || side < 1e-12 {
		return fmt.Errorf("scene3d: invalid perspective camera")
	}
	return nil
}
