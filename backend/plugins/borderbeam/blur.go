// Author: HyperparameterSI
package borderbeam

import (
	"image"
	"math"
)

// blurred downsamples only the soft light, then applies a separable Gaussian.
// Its radius never affects the crisp border texture. Buffers survive repaint.
func (l *rasterLayer) blurred(sigma, scale float64) *image.RGBA {
	step := max(1, int(sigma*scale/2))
	src := l.pixels
	w, h := (src.Rect.Dx()+step-1)/step, (src.Rect.Dy()+step-1)/step
	for i := range l.scratch {
		if l.scratch[i] == nil || l.scratch[i].Rect.Dx() != w || l.scratch[i].Rect.Dy() != h {
			l.scratch[i] = image.NewRGBA(image.Rect(0, 0, w, h))
		}
	}
	a, b := l.scratch[0], l.scratch[1]
	// Area sampling preserves narrow rim energy before convolution.
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sum [4]int
			count := 0
			for sy := y * step; sy < min((y+1)*step, src.Rect.Dy()); sy++ {
				for sx := x * step; sx < min((x+1)*step, src.Rect.Dx()); sx++ {
					at := sy*src.Stride + sx*4
					for k := 0; k < 4; k++ {
						sum[k] += int(src.Pix[at+k])
					}
					count++
				}
			}
			at := y*a.Stride + x*4
			for k := 0; k < 4; k++ {
				a.Pix[at+k] = uint8(sum[k] / count)
			}
		}
	}
	s := sigma * scale / float64(step)
	radius := int(math.Ceil(3 * s))
	weights := make([]float64, 2*radius+1)
	total := 0.
	for i := range weights {
		d := float64(i - radius)
		weights[i] = math.Exp(-d * d / (2 * s * s))
		total += weights[i]
	}
	for i := range weights {
		weights[i] /= total
	}
	for pass := 0; pass < 2; pass++ {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				var sum [4]float64
				for i, weight := range weights {
					sx, sy := x, y
					if pass == 0 {
						sx += i - radius
					} else {
						sy += i - radius
					}
					if sx < 0 || sx >= w || sy < 0 || sy >= h {
						continue
					}
					at := sy*a.Stride + sx*4
					for k := 0; k < 4; k++ {
						sum[k] += float64(a.Pix[at+k]) * weight
					}
				}
				at := y*b.Stride + x*4
				for k := 0; k < 4; k++ {
					b.Pix[at+k] = uint8(sum[k] + .5)
				}
			}
		}
		a, b = b, a
	}
	return a
}
