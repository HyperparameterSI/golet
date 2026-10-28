// Author: HyperparameterSI

package scene3d

import (
	"bytes"
	"image/color"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/HyperparameterSI/golet/backend/ui"
)

func testCamera() Camera {
	return Camera{Position: Vec3{0, 0, 3}, Target: Vec3{}, Up: Vec3{0, 1, 0}, FOVY: math.Pi / 2, Near: .1, Far: 10}
}
func triangle(t testing.TB, z float32) *Mesh {
	t.Helper()
	m, err := NewMesh([]Vec3{{-1, -1, z}, {1, -1, z}, {0, 1, z}}, []uint32{0, 1, 2})
	if err != nil {
		t.Fatal(err)
	}
	return m
}
func testScene(m *Mesh) *Scene {
	return &Scene{Instances: []Instance{{Mesh: m, Material: Material{Color: color.RGBA{220, 40, 20, 255}, Unlit: true}}}, Background: color.RGBA{10, 12, 16, 255}}
}
func render(t testing.TB, s *Scene, c Camera, target *Target) Stats {
	t.Helper()
	stats, err := (Software{}).Render(target, Request{s, c, 96, 96})
	if err != nil {
		t.Fatal(err)
	}
	return stats
}

func TestDepthIndependentOfSubmissionOrder(t *testing.T) {
	near := testScene(triangle(t, 1))
	far := testScene(triangle(t, 0))
	far.Instances[0].Material.Color = color.RGBA{20, 80, 240, 255}
	s := *near
	s.Instances = append(s.Instances, far.Instances[0])
	var target Target
	render(t, &s, testCamera(), &target)
	first := append([]byte(nil), target.Image.Pix...)
	if got := target.Image.RGBAAt(48, 48); got != near.Instances[0].Material.Color {
		t.Fatalf("near triangle did not win: %v", got)
	}
	s.Instances[0], s.Instances[1] = s.Instances[1], s.Instances[0]
	render(t, &s, testCamera(), &target)
	if !bytes.Equal(first, target.Image.Pix) {
		t.Fatal("depth depends on submission order")
	}
}
func TestClippingAndCulling(t *testing.T) {
	for _, z := range []float32{4, -20} {
		s := testScene(triangle(t, z))
		var target Target
		if stats := render(t, s, testCamera(), &target); stats.CulledInstances != 1 || stats.ShadedPixels != 0 {
			t.Fatalf("outside bounds: %+v", stats)
		}
	}
	// One vertex behind the camera, two in front. Dropping the entire triangle
	// would lose its visible region; clipping must retain it without huge bounds.
	m, err := NewMesh([]Vec3{{-.5, -.5, 0}, {.5, -.5, 0}, {0, .5, 3.2}}, []uint32{0, 1, 2})
	if err != nil {
		t.Fatal(err)
	}
	s := testScene(m)
	s.Instances[0].Material.DoubleSided = true
	var target Target
	if stats := render(t, s, testCamera(), &target); stats.ShadedPixels == 0 {
		t.Fatal("near-plane crossing disappeared")
	}
	for _, z := range target.depth {
		if !finite(z) || z < 0 || z > 10.001 {
			t.Fatalf("invalid clipped depth: %g", z)
		}
	}
}
func TestBackfaceAndLighting(t *testing.T) {
	m := triangle(t, 0)
	s := testScene(m)
	s.Instances[0].Transform = Scaling(Vec3{-1, 1, 1})
	var target Target
	if stats := render(t, s, testCamera(), &target); stats.ShadedPixels != 0 {
		t.Fatal("back face not culled")
	}
	s.Instances[0].Material.DoubleSided = true
	if stats := render(t, s, testCamera(), &target); stats.ShadedPixels == 0 {
		t.Fatal("double-sided face disappeared")
	}
	s.Instances[0].Transform = Identity()
	s.Instances[0].Material.Unlit = false
	render(t, s, testCamera(), &target)
	dark := target.Image.RGBAAt(48, 48)
	s.Lights = []DirectionalLight{{Vec3{0, 0, 1}, 1}}
	render(t, s, testCamera(), &target)
	if target.Image.RGBAAt(48, 48).R <= dark.R {
		t.Fatal("light did not illuminate surface")
	}
}
func TestValidationAndBufferReuse(t *testing.T) {
	s := testScene(Cube())
	var target Target
	r := Request{s, testCamera(), 96, 96}
	if _, err := (Software{}).Render(&target, r); err != nil {
		t.Fatal(err)
	}
	pixels := &target.Image.Pix[0]
	depth := &target.depth[0]
	allocs := testing.AllocsPerRun(20, func() {
		if _, err := (Software{}).Render(&target, r); err != nil {
			panic(err)
		}
	})
	if allocs != 0 {
		t.Fatalf("warm render allocates: %g", allocs)
	}
	if pixels != &target.Image.Pix[0] || depth != &target.depth[0] {
		t.Fatal("buffers replaced")
	}
	r.Width = MaxPixels
	r.Height = 2
	if _, err := (Software{}).Render(&target, r); err == nil {
		t.Fatal("oversize accepted")
	}
	r.Width = 96
	r.Height = 96
	r.Camera.Near = 0
	if _, err := (Software{}).Render(&target, r); err == nil {
		t.Fatal("invalid camera accepted")
	}
	target.Release()
	if target.BufferBytes() != 0 {
		t.Fatal("buffers retained after release")
	}
	if _, err := NewMesh([]Vec3{{float32(math.NaN()), 0, 0}}, []uint32{0, 0, 0}); err == nil {
		t.Fatal("NaN accepted")
	}
}
func TestOBJ(t *testing.T) {
	valid := "v -1 -1 0\nv 1 -1 0\nv 0 1 0\nf -3/1/1 -2/2/1 -1/3/1\n"
	m, err := LoadOBJ(strings.NewReader(valid), OBJLimits{})
	if err != nil || m.TriangleCount() != 1 {
		t.Fatalf("load: %v", err)
	}
	for _, input := range []string{"v NaN 0 0", "v 0 0 0\nf 0 1 1", "v 0 0 0\nf 1 2 3", "v 0 0 0\nf 1 1 1 1"} {
		if _, err := LoadOBJ(strings.NewReader(input), OBJLimits{}); err == nil {
			t.Fatalf("invalid model accepted: %q", input)
		}
	}
	if _, err := LoadOBJ(strings.NewReader(valid), OBJLimits{Bytes: 10}); err == nil {
		t.Fatal("byte budget ignored")
	}
	if _, err := LoadOBJ(strings.NewReader(valid), OBJLimits{Vertices: 2}); err == nil {
		t.Fatal("vertex budget ignored")
	}
	if _, err := LoadOBJ(strings.NewReader(valid), OBJLimits{Bytes: int64(len(valid))}); err != nil {
		t.Fatal("exact byte budget rejected", err)
	}
}
func TestViewCacheAndCameraInvalidation(t *testing.T) {
	o := ViewOptions{Scene: testScene(Cube()), Camera: testCamera()}
	renders := 0
	o.OnFrame = func(_ Stats, err error) {
		if err != nil {
			t.Fatal(err)
		}
		renders++
	}
	tt := ui.NewTester(func(c *ui.Context) { View(c, o).Fill() }, 120, 100)
	tt.Frame()
	tt.Frame()
	if renders != 1 {
		t.Fatalf("static renders: %d", renders)
	}
	if tt.FramePending() {
		t.Fatal("static viewport requests frames")
	}
	o.Camera.Position.X = 1
	tt.Frame()
	if renders != 2 {
		t.Fatal("camera change ignored")
	}
	o.Scene.Instances[0].Material.Color = color.RGBA{20, 180, 70, 255}
	o.Revision++
	tt.Frame()
	if renders != 3 {
		t.Fatal("revision ignored")
	}
}
func FuzzOBJ(f *testing.F) {
	f.Add("v 0 0 0\nv 1 0 0\nv 0 1 0\nf 1 2 3\n")
	f.Fuzz(func(t *testing.T, input string) {
		_, _ = LoadOBJ(strings.NewReader(input), OBJLimits{Bytes: 4096, Vertices: 100, Triangles: 100})
	})
}
func BenchmarkSoftware(b *testing.B) {
	for _, count := range []int{1, 100} {
		b.Run(stringName(count), func(b *testing.B) {
			mesh := Cube()
			s := testScene(mesh)
			s.Instances = nil
			for i := 0; i < count; i++ {
				s.Instances = append(s.Instances, Instance{Mesh: mesh, Transform: Translation(Vec3{float32(i%10) - 4.5, 0, -float32(i/10) * 2}), Material: Material{Color: color.RGBA{80, 160, 240, 255}}})
			}
			camera := DefaultCamera()
			camera.Position = Vec3{4, 6, 12}
			s.Ambient = .3
			s.Lights = []DirectionalLight{{Vec3{2, 4, 3}, .7}}
			var target Target
			r := Request{s, camera, 640, 360}
			if _, err := (Software{}).Render(&target, r); err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := (Software{}).Render(&target, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
func stringName(n int) string {
	if n == 1 {
		return "1-cube-640x360"
	}
	return "100-cubes-640x360"
}

func TestSharedEdgeHasSingleOwner(t *testing.T) {
	vertices := []Vec3{{-1, -1, 0}, {1, -1, 0}, {1, 1, 0}, {-1, 1, 0}}
	left, err := NewMesh(vertices, []uint32{0, 1, 2})
	if err != nil {
		t.Fatal(err)
	}
	right, err := NewMesh(vertices, []uint32{0, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	s := testScene(left)
	s.Instances = append(s.Instances, Instance{Mesh: right, Material: Material{Color: color.RGBA{0, 180, 220, 255}, Unlit: true}})
	var target Target
	render(t, s, testCamera(), &target)
	first := append([]byte(nil), target.Image.Pix...)
	s.Instances[0], s.Instances[1] = s.Instances[1], s.Instances[0]
	render(t, s, testCamera(), &target)
	if !bytes.Equal(first, target.Image.Pix) {
		t.Fatal("shared diagonal depends on submission order")
	}
	for y := 34; y < 62; y++ {
		for x := 34; x < 62; x++ {
			if target.Image.RGBAAt(x, y) == s.Background {
				t.Fatal("crack in joined triangles")
			}
		}
	}
}

func TestMeshOwnsInputs(t *testing.T) {
	vertices := []Vec3{{-1, -1, 0}, {1, -1, 0}, {0, 1, 0}}
	indices := []uint32{0, 1, 2}
	m, err := NewMesh(vertices, indices)
	if err != nil {
		t.Fatal(err)
	}
	vertices[0].X = 99
	indices[0] = 99
	v, i := m.Geometry()
	if v[0].X != -1 || i[0] != 0 {
		t.Fatal("mesh aliases caller data")
	}
	v[0].X = 99
	if m.vertices[0].X != -1 {
		t.Fatal("export aliases mesh")
	}
}

func TestOffscreenViewDoesNotAnimate(t *testing.T) {
	updates := 0
	o := ViewOptions{Scene: testScene(Cube()), Camera: testCamera(), Animate: func(_ time.Time) { updates++ }}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Scroll(c).Fill().Children(func() {
			ui.Box(c).Height(300).Shrink(0)
			View(c, o).Size(100, 100).Shrink(0)
		})
	}, 120, 100)
	tt.Frame()
	if updates != 0 || tt.FramePending() {
		t.Fatal("offscreen view continues animating")
	}
}

func TestInvalidAnimatedSceneStops(t *testing.T) {
	failures := 0
	tt := ui.NewTester(func(c *ui.Context) {
		View(c, ViewOptions{Camera: testCamera(), Animate: func(_ time.Time) {}, OnFrame: func(_ Stats, err error) {
			if err == nil {
				t.Fatal("nil scene succeeded")
			}
			failures++
		}}).Fill()
	}, 100, 100)
	tt.Frame()
	if failures != 1 || tt.FramePending() {
		t.Fatal("invalid scene was hidden or keeps requesting frames")
	}
}
