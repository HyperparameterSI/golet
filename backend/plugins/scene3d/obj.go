// Author: HyperparameterSI

package scene3d

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// OBJLimits bounds untrusted input. Zero fields select defaults: 16 MiB,
// 250,000 vertices and 500,000 triangles. Lines are limited to 64 KiB.
type OBJLimits struct {
	Bytes               int64
	Vertices, Triangles int
}

// LoadOBJ loads triangulated Wavefront OBJ positions and faces. Positive and
// relative negative indices and v/vt/vn syntax are accepted. UVs, normals and
// materials are not imported: the initial renderer uses flat solid materials.
// Polygons must be triangulated by the exporter. No external files are opened.
func LoadOBJ(r io.Reader, limits OBJLimits) (*Mesh, error) {
	if r == nil {
		return nil, fmt.Errorf("scene3d: nil OBJ reader")
	}
	if limits.Bytes == 0 {
		limits.Bytes = 16 << 20
	}
	if limits.Vertices == 0 {
		limits.Vertices = 250000
	}
	if limits.Triangles == 0 {
		limits.Triangles = 500000
	}
	if limits.Bytes < 0 || limits.Vertices < 0 || limits.Triangles < 0 {
		return nil, fmt.Errorf("scene3d: invalid OBJ limits")
	}
	limited := &io.LimitedReader{R: r, N: limits.Bytes}
	scan := bufio.NewScanner(limited)
	var vertices []Vec3
	var indices []uint32
	line := 0
	fail := func(msg string) (*Mesh, error) { return nil, fmt.Errorf("scene3d: OBJ line %d: %s", line, msg) }
	for scan.Scan() {
		line++
		text, _, _ := strings.Cut(scan.Text(), "#")
		fields := strings.Fields(text)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "v":
			if len(fields) != 4 {
				return fail("expected three position coordinates")
			}
			if len(vertices) >= limits.Vertices {
				return fail("vertex limit exceeded")
			}
			var xyz [3]float32
			for i := range xyz {
				f, err := strconv.ParseFloat(fields[i+1], 32)
				if err != nil || !finite(float32(f)) {
					return fail("invalid position")
				}
				xyz[i] = float32(f)
			}
			vertices = append(vertices, Vec3{xyz[0], xyz[1], xyz[2]})
		case "f":
			if len(fields) != 4 {
				return fail("export triangulated faces (three vertices per face)")
			}
			if len(indices)/3 >= limits.Triangles {
				return fail("triangle limit exceeded")
			}
			for _, field := range fields[1:] {
				position, _, _ := strings.Cut(field, "/")
				i, err := strconv.ParseInt(position, 10, 64)
				if err != nil || i == 0 {
					return fail("invalid face index")
				}
				if i < 0 {
					i = int64(len(vertices)) + i
				} else {
					i--
				}
				if i < 0 || i >= int64(len(vertices)) {
					return fail("face index out of range")
				}
				indices = append(indices, uint32(i))
			}
		case "vt", "vn", "o", "g", "s", "usemtl", "mtllib": // Metadata is not required for flat rendering.
		default:
			return fail("unsupported record " + fields[0])
		}
	}
	if err := scan.Err(); err != nil {
		return fail(err.Error())
	}
	if limited.N == 0 {
		var b [1]byte
		n, err := io.ReadFull(r, b[:])
		if n != 0 {
			return fail("byte limit exceeded")
		}
		if err != nil && err != io.EOF {
			return fail(err.Error())
		}
	}
	return ownedMesh(vertices, indices)
}
