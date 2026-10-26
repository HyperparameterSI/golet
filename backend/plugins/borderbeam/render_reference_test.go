// Author: HyperparameterSI
// Pre-optimization renderer retained as a pixel-equivalence oracle.
package borderbeam

import "math"

func referenceRenderLayers(layers *[3]rasterLayer, w, h, radius float64, f colorField, o Options, strength float64) {
	padding := 0.
	if f.mode == PulseOutside {
		padding = 60
	}
	// Two samples per DIP keep the one-DIP rim sharp on Retina displays. Cap
	// texture area for unusually large windows rather than allocating without bound.
	scale := min(2., math.Sqrt(2_000_000/((w+2*padding)*(h+2*padding))))
	pw, ph := max(1, int(math.Ceil((w+2*padding)*scale))), max(1, int(math.Ceil((h+2*padding)*scale)))
	for i := range layers {
		layers[i].resize(pw, ph)
		clear(layers[i].pixels.Pix)
	}
	phase := fract(f.seconds / f.duration)
	co, si := math.Cos(f.hue*math.Pi/180), math.Sin(f.hue*math.Pi/180)
	opacity := [3]float64{.48, .7, .8}
	if f.mode == Line {
		opacity = [3]float64{.72, .7, .8}
	}
	if f.mode == PulseInner {
		opacity = [3]float64{.32, .4, .8}
		if f.dark {
			opacity = [3]float64{1.54, .44, .66}
		}
	}
	if f.mode == PulseOutside {
		opacity = [3]float64{1.96, 1.04, .42}
		if f.dark {
			opacity = [3]float64{.94, .34, .3}
		}
	}
	if o.Palette == Mono {
		strength *= .5
	}
	rim := float64(dimension(o.Width, 1, 8))
	depth := float64(dimension(o.Glow, 28, 48))
	for py := 0; py < ph; py++ {
		y := (float64(py)+.5)/scale - padding
		for px := 0; px < pw; px++ {
			x := (float64(px)+.5)/scale - padding
			distance := roundedDistance(x, y, w, h, radius)
			inside := clamp(.5 - distance*scale)
			stroke := inside * clamp(.5+(distance+rim)*scale)
			if f.mode != PulseOutside && inside == 0 {
				continue
			}
			vertical := clamp(1 - min(y, h-y)/depth)
			horizontal := clamp(1 - min(x, w-x)/depth)
			edge := vertical + horizontal - vertical*horizontal
			masks := [3]float64{stroke, edge * inside, stroke}
			if (f.mode == Medium || f.mode == PulseInner) && edge == 0 && stroke == 0 {
				continue
			}
			if f.mode < Line {
				spin := spinMask(x, y, w, h, phase, false)
				masks[0] *= spin
				if f.mode == Medium {
					masks[1] *= spin
				}
			}
			if f.mode == Small {
				masks[1] = inside * spinMask(x, y, w, h, phase, true)
			}
			if f.mode == Line {
				location := lineTravel(phase) * w
				bw := gradient(phase, [2]float64{0, .5}, [2]float64{.1, .8}, [2]float64{.2, 1.1}, [2]float64{.3, 1.3}, [2]float64{.4, 1.45}, [2]float64{.5, 1.5}, [2]float64{.6, 1.45}, [2]float64{.7, 1.3}, [2]float64{.8, 1.1}, [2]float64{.9, .8}, [2]float64{1, .5})
				bh := gradient(fract(f.seconds/(f.duration*1.3)), [2]float64{0, .8}, [2]float64{.25, 1.25}, [2]float64{.55, .85}, [2]float64{.8, 1.3}, [2]float64{1, .8})
				radial := math.Hypot((x-location)/(78*bw), (y-h)/(60*bh))
				mask := gradient(radial, [2]float64{0, 1}, [2]float64{.45, .5}, [2]float64{1, 0}) * lineEdge(phase)
				masks[0] *= mask
				masks[1] *= mask
				masks[2] = inside * gradient(math.Hypot((x-location)/(84*bw), (y-h)/(110*bh)), [2]float64{0, 1}, [2]float64{.35, .5}, [2]float64{1, 0}) * lineEdge(phase)
			}
			if f.mode == PulseOutside {
				masks[1], masks[2] = 1-inside, 1-inside
			}
			for layer, mask := range masks {
				if mask <= 0 {
					continue
				}
				sx, sy := x, y
				if f.mode == PulseOutside && layer > 0 {
					inset := 10.
					if layer == 2 {
						inset = 30
					}
					// Inverse of the reference inset box and scale(.95,.9).
					sx = ((x-w/2)/.95 + w/2 + inset) * w / (w + 2*inset)
					sy = ((y-h/2)/.9 + h/2 + inset) * h / (h + 2*inset)
				}
				col := f.sample(layer, sx, sy)
				if f.mode < Line && layer == 2 {
					angle := fract(math.Atan2(y-h/2, x-w/2)/(2*math.Pi) + .25 - phase)
					highlight := gradient(angle, [2]float64{0, 0}, [2]float64{.58, 0}, [2]float64{.65, .08}, [2]float64{.70, .85}, [2]float64{.705, .85}, [2]float64{.735, .2}, [2]float64{.82, 0}, [2]float64{1, 0})
					v := 0.
					if f.dark {
						v = 1
					}
					col = vector{v, v, v, highlight}
				}
				if col[3] <= 0 {
					continue
				}
				// Field accumulation is premultiplied; apply CSS color filters to straight RGB.
				for k := 0; k < 3; k++ {
					col[k] /= col[3]
				}
				col = hueRotate(col, co, si, f.brightness, f.saturation)
				alpha := clamp(col[3] * mask * opacity[layer] * strength)
				at := py*layers[layer].pixels.Stride + px*4
				pixel := layers[layer].pixels.Pix[at : at+4]
				for k := 0; k < 3; k++ {
					pixel[k] = uint8(clamp(col[k])*alpha*255 + .5)
				}
				pixel[3] = uint8(alpha*255 + .5)
			}
		}
	}
	for i := range layers {
		pixels := layers[i].pixels
		sigma := 0.
		if i == 2 {
			sigma = 8
			if f.mode == Line {
				sigma = 6
			}
			if f.mode == PulseOutside {
				sigma = 15
				if f.dark {
					sigma = 22.5
				}
			}
		}
		if i == 1 && f.mode == PulseOutside {
			sigma = 6
			if f.dark {
				sigma = 3
			}
		}
		if sigma > 0 {
			pixels = layers[i].blurred(sigma, scale)
		}
		layers[i].bitmap.Update(pixels)
	}
}
