// Author: HyperparameterSI

package scene3d

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// MaxPixels bounds color/depth buffer allocation (32 MiB combined). A viewport
// can render at lower resolution than its display rectangle.
const MaxPixels = 4 * 1024 * 1024

// Target holds reusable CPU buffers. Use one per viewport/rendering goroutine.
// Image is borrowed until the next Render; copy it to retain a frame.
type Target struct {
	Image       *image.RGBA
	depth       []float32
	world, view []Vec3
}

// Release drops CPU buffers; the UI compositor owns its uploaded texture.
func (t *Target) Release() { *t = Target{} }

// BufferBytes reports retained color, depth and transformed-vertex capacity,
// excluding immutable meshes, Go headers and the UI compositor's GPU texture.
func (t *Target) BufferBytes() int {
	n := 0
	if t.Image != nil {
		n = cap(t.Image.Pix)
	}
	return n + 4*cap(t.depth) + 12*(cap(t.world)+cap(t.view))
}
func (t *Target) resize(w, h int) error {
	if w <= 0 || h <= 0 || w > MaxPixels || h > MaxPixels/w {
		return fmt.Errorf("scene3d: render size %dx%d exceeds pixel budget", w, h)
	}
	if t.Image == nil || t.Image.Rect.Dx() != w || t.Image.Rect.Dy() != h {
		t.Image = image.NewRGBA(image.Rect(0, 0, w, h))
		t.depth = make([]float32, w*h)
	}
	return nil
}

// Request describes one immutable frame. Scene and its slices must not change
// until Render returns. Width and Height are render pixels, not window DIPs.
type Request struct {
	Scene         *Scene
	Camera        Camera
	Width, Height int
}

// Stats counts actual work; it does not measure GPU time or promise FPS.
type Stats struct {
	Instances, CulledInstances, Triangles, RasterizedTriangles, ShadedPixels int
	BufferBytes                                                              int
}

// Renderer allows alternative CPU/image-producing renderers without changing
// scene ownership or the viewport. Direct GPU texture sharing is not yet part
// of this interface and requires a future compositor integration.
type Renderer interface {
	Render(*Target, Request) (Stats, error)
}

// Software is a stateless, flat-shaded reference renderer. Different targets
// may be rendered concurrently. It uses no worker pool or global resources.
type Software struct{}

type projection struct {
	pos, right, up, forward Vec3
	tx, ty, near, far       float32
}

func project(c Camera, aspect float32) projection {
	f := c.Target.Sub(c.Position).Unit()
	r := f.Cross(c.Up).Unit()
	ty := float32(math.Tan(float64(c.FOVY) / 2))
	return projection{c.Position, r, r.Cross(f), f, ty * aspect, ty, c.Near, c.Far}
}
func (p projection) view(v Vec3) Vec3 {
	d := v.Sub(p.pos)
	return Vec3{d.Dot(p.right), d.Dot(p.up), d.Dot(p.forward)}
}
func (p projection) distance(v Vec3, plane int) float32 {
	switch plane {
	case 0:
		return v.Z - p.near
	case 1:
		return p.far - v.Z
	case 2:
		return v.X + v.Z*p.tx
	case 3:
		return v.Z*p.tx - v.X
	case 4:
		return v.Y + v.Z*p.ty
	default:
		return v.Z*p.ty - v.Y
	}
}
func (p projection) culled(m *Mesh, transform Mat4) bool {
	mask := uint8(63)
	for i := 0; i < 8; i++ {
		v := m.min
		if i&1 != 0 {
			v.X = m.max.X
		}
		if i&2 != 0 {
			v.Y = m.max.Y
		}
		if i&4 != 0 {
			v.Z = m.max.Z
		}
		v = p.view(transform.point(v))
		var outside uint8
		for plane := 0; plane < 6; plane++ {
			if p.distance(v, plane) < 0 {
				outside |= 1 << plane
			}
		}
		mask &= outside
	}
	return mask != 0
}

// Render validates before allocating and reuses buffers at a stable size.
// Meshes are shared, invisible instance bounds are rejected before transforming
// vertices, and depth is tested before writing a pixel.
func (Software) Render(t *Target, r Request) (Stats, error) {
	var stats Stats
	if t == nil || r.Scene == nil {
		return stats, fmt.Errorf("scene3d: nil target or scene")
	}
	if err := r.Camera.validate(); err != nil {
		return stats, err
	}
	s := r.Scene
	if !finite(s.Ambient) || s.Ambient < 0 || s.Background.A != 255 {
		return stats, fmt.Errorf("scene3d: ambient must be nonnegative and background opaque")
	}
	for i, l := range s.Lights {
		if !l.Direction.valid() || !finite(l.Intensity) || l.Intensity < 0 || l.Direction.Dot(l.Direction) == 0 {
			return stats, fmt.Errorf("scene3d: invalid light %d", i)
		}
	}
	for i, in := range s.Instances {
		if in.Mesh == nil || len(in.Mesh.indices) == 0 || (in.Transform != (Mat4{}) && !in.Transform.affine()) || in.Material.Color.A != 255 {
			return stats, fmt.Errorf("scene3d: invalid mesh, transform or opaque material at instance %d", i)
		}
	}
	if err := t.resize(r.Width, r.Height); err != nil {
		return stats, err
	}
	for i := range t.depth {
		t.depth[i] = 0
		j := i * 4
		t.Image.Pix[j], t.Image.Pix[j+1], t.Image.Pix[j+2], t.Image.Pix[j+3] = s.Background.R, s.Background.G, s.Background.B, 255
	}
	p := project(r.Camera, float32(r.Width)/float32(r.Height))
	for _, in := range s.Instances {
		stats.Instances++
		m := in.Mesh
		transform := in.Transform
		if transform == (Mat4{}) {
			transform = Identity()
		}
		if p.culled(m, transform) {
			stats.CulledInstances++
			continue
		}
		n := len(m.vertices)
		if cap(t.world) < n {
			t.world = make([]Vec3, n)
			t.view = make([]Vec3, n)
		}
		t.world = t.world[:n]
		t.view = t.view[:n]
		for i, v := range m.vertices {
			t.world[i] = transform.point(v)
			t.view[i] = p.view(t.world[i])
		}
		for i := 0; i < len(m.indices); i += 3 {
			stats.Triangles++
			ia, ib, ic := m.indices[i], m.indices[i+1], m.indices[i+2]
			a, b, c := t.world[ia], t.world[ib], t.world[ic]
			normal := b.Sub(a).Cross(c.Sub(a))
			length := normal.Dot(normal)
			if !finite(length) || length < 1e-20 {
				continue
			}
			facing := normal.Dot(r.Camera.Position.Sub(a))
			if !in.Material.DoubleSided && facing <= 0 {
				continue
			}
			normal = normal.Unit()
			if facing < 0 {
				normal = normal.Scale(-1)
			}
			brightness := float32(1)
			if !in.Material.Unlit {
				brightness = s.Ambient
				for _, l := range s.Lights {
					brightness += max(float32(0), normal.Dot(l.Direction.Unit())) * l.Intensity
				}
			}
			col := shade(in.Material.Color, brightness)
			var poly, work [12]Vec3
			poly[0], poly[1], poly[2] = t.view[ia], t.view[ib], t.view[ic]
			count := 3
			for plane := 0; plane < 6 && count > 0; plane++ {
				next := 0
				prev := poly[count-1]
				pd := p.distance(prev, plane)
				for j := 0; j < count; j++ {
					cur := poly[j]
					cd := p.distance(cur, plane)
					if (pd >= 0) != (cd >= 0) {
						work[next] = prev.Add(cur.Sub(prev).Scale(pd / (pd - cd)))
						next++
					}
					if cd >= 0 {
						work[next] = cur
						next++
					}
					prev, pd = cur, cd
				}
				poly, work = work, poly
				count = next
			}
			for j := 1; j+1 < count; j++ {
				t.triangle(p, poly[0], poly[j], poly[j+1], col, &stats)
			}
		}
	}
	stats.BufferBytes = t.BufferBytes()
	return stats, nil
}
func shade(c color.RGBA, intensity float32) color.RGBA {
	// Approximate linear-light shading rather than multiplying sRGB bytes.
	f := float32(math.Sqrt(float64(max(float32(0), intensity))))
	return color.RGBA{uint8(min(float32(255), float32(c.R)*f)), uint8(min(float32(255), float32(c.G)*f)), uint8(min(float32(255), float32(c.B)*f)), 255}
}

type screenVertex struct{ x, y, iz float32 }

func edge(a, b screenVertex, x, y float32) float32 { return (b.x-a.x)*(y-a.y) - (b.y-a.y)*(x-a.x) }

// The top-left fill rule assigns a shared edge to exactly one triangle.
func inclusiveEdge(a, b screenVertex) bool {
	return b.y < a.y || (b.y == a.y && b.x > a.x)
}
func (t *Target) triangle(p projection, a, b, c Vec3, col color.RGBA, stats *Stats) {
	w, h := t.Image.Rect.Dx(), t.Image.Rect.Dy()
	toScreen := func(v Vec3) screenVertex {
		iz := 1 / v.Z
		return screenVertex{(v.X*iz/p.tx + 1) * float32(w) / 2, (1 - v.Y*iz/p.ty) * float32(h) / 2, iz}
	}
	v0, v1, v2 := toScreen(a), toScreen(b), toScreen(c)
	area := edge(v0, v1, v2.x, v2.y)
	if !finite(area) || float32(math.Abs(float64(area))) < 1e-8 {
		return
	}
	if area < 0 {
		v1, v2 = v2, v1
		area = -area
	}
	minX, maxX := max(0, int(math.Floor(float64(min(v0.x, v1.x, v2.x))))), min(w-1, int(math.Ceil(float64(max(v0.x, v1.x, v2.x)))))
	minY, maxY := max(0, int(math.Floor(float64(min(v0.y, v1.y, v2.y))))), min(h-1, int(math.Ceil(float64(max(v0.y, v1.y, v2.y)))))
	stats.RasterizedTriangles++
	inv := 1 / area
	in0, in1, in2 := inclusiveEdge(v1, v2), inclusiveEdge(v2, v0), inclusiveEdge(v0, v1)
	// Increment edge equations across rows; reciprocal depth is linear in screen space.
	e0, e1, e2 := edge(v1, v2, float32(minX)+.5, float32(minY)+.5), edge(v2, v0, float32(minX)+.5, float32(minY)+.5), edge(v0, v1, float32(minX)+.5, float32(minY)+.5)
	for y := minY; y <= maxY; y++ {
		q0, q1, q2 := e0, e1, e2
		for x := minX; x <= maxX; x++ {
			if (q0 > 0 || q0 == 0 && in0) && (q1 > 0 || q1 == 0 && in1) && (q2 > 0 || q2 == 0 && in2) {
				z := (q0*v0.iz + q1*v1.iz + q2*v2.iz) * inv
				at := y*w + x
				if z > t.depth[at] {
					t.depth[at] = z
					j := at * 4
					t.Image.Pix[j], t.Image.Pix[j+1], t.Image.Pix[j+2] = col.R, col.G, col.B
					stats.ShadedPixels++
				}
			}
			q0 -= v2.y - v1.y
			q1 -= v0.y - v2.y
			q2 -= v1.y - v0.y
		}
		e0 += v2.x - v1.x
		e1 += v0.x - v2.x
		e2 += v1.x - v0.x
	}
}
