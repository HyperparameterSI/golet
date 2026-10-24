// Author: HyperparameterSI
// Port of Frad's WGSL color fields. See reference/LICENSE for upstream credits.
package borderbeam

import (
	"github.com/HyperparameterSI/golet/backend/ui"
	"math"
)

type vector [4]float64
type spot struct{ shape, color vector }
type colorField struct {
	layers                                                        [3][]spot
	width, height, seconds, duration, hue, brightness, saturation float64
	mode                                                          Preset
	dark                                                          bool
}

func fract(x float64) float64 { return x - math.Floor(x) }
func clamp(x float64) float64 { return max(0, min(1, x)) }
func oscillate(a, b, period, delay, t float64) float64 {
	return a + (b-a)*(1-math.Cos(2*math.Pi*(t-delay)/max(period, .001)))/2
}
func mix(a, b, t float64) float64 { return a + (b-a)*t }
func table(name string) []vector  { return referenceTables[name] }

func makeField(w, h, t, duration float64, p Preset, o Options, dark bool, colors []ui.Color) colorField {
	f := colorField{width: w, height: h, seconds: t, duration: duration, mode: p, dark: dark, brightness: 1.3, saturation: 1.2}
	if !dark {
		f.saturation = 1.5
		if p == Small {
			f.saturation = 1.8
		}
		if p == Line {
			f.saturation = 1.95
		}
	}
	if o.Brightness > 0 {
		f.brightness = float64(o.Brightness)
	}
	if o.Saturation > 0 {
		f.saturation = float64(o.Saturation)
	}
	hueRange := float64(30)
	if o.HueRange != 0 {
		hueRange = float64(o.HueRange)
	}
	if p == Line {
		hueRange = min(hueRange, 13)
	}
	f.hue = -math.Cos(2*math.Pi*t/12) * hueRange
	if p >= PulseInner {
		period := 14.
		if p == PulseInner {
			period = 16
		}
		f.hue = fract(t/period) * 360
	}
	if o.StaticColors || o.Palette == Mono {
		f.hue = 0
	}
	variant := "COLORFUL"
	switch o.Palette {
	case Mono:
		variant = "MONO"
	case Ocean:
		variant = "OCEAN"
	case Sunset:
		variant = "SUNSET"
	}
	add := func(layer int, shapeName, colorName string, kind int) {
		shapes, palette := table(shapeName), table(colorName)
		for i, s := range shapes {
			if i >= len(palette) {
				break
			}
			col := palette[i]
			if len(colors) > 0 {
				c := colors[i%len(colors)]
				col = vector{float64(c.R) / 255, float64(c.G) / 255, float64(c.B) / 255, col[3] * float64(c.A) / 255}
			}
			switch {
			case p == Line:
				if layer < 2 {
					s[0] += lineTravel(fract(t/duration)) * w
					s[1] += h
					col[3] *= lineEdge(fract(t / duration))
				} else {
					s[0] = s[0] / 100 * w
					s[1] = s[1] / 100 * h
				}
			case p >= PulseInner:
				s[0] *= w
				s[1] *= h
				if kind != 2 {
					s, col = f.pulseSpot(s, col, i, kind)
				} else if p == PulseOutside {
					s[2] *= max(.35, min(4, w/350))
					s[3] *= max(.35, min(4, h/140))
				}
			default:
				s[0] *= w
				s[1] *= h
				if layer == 1 {
					s[2] *= .9
					s[3] *= .9
				}
			}
			f.layers[layer] = append(f.layers[layer], spot{s, col})
		}
	}
	switch p {
	case Medium, Small:
		prefix := "MD"
		if p == Small {
			prefix = "SM"
		}
		add(0, prefix+"_SHAPES", prefix+"_COLORS_"+variant, 0)
		add(1, prefix+"_SHAPES", prefix+"_INNER_COLORS_"+variant, 0)
	case Line:
		shade := "LIGHT"
		if dark {
			shade = "DARK"
		}
		add(0, "LINE_SHAPES_"+shade, "LINE_COLORS_"+variant+"_"+shade, 0)
		add(1, "LINE_INNER_SHAPES", "LINE_INNER_COLORS_"+variant, 0)
		// The first two bloom lights are followed by five narrow spikes.
		palette := append(append([]vector(nil), table("LINE_BLOOM_HEAD_"+variant+"_"+shade)...), table("LINE_BLOOM_SPIKES_"+variant+"_"+shade)...)
		for i, col := range palette {
			if len(colors) > 0 {
				c := colors[i%len(colors)]
				col = vector{float64(c.R) / 255, float64(c.G) / 255, float64(c.B) / 255, col[3] * float64(c.A) / 255}
			}
			s := table("LINE_BLOOM_SHAPES")[i]
			s[0] = s[0] / 100 * w
			s[1] = s[1] / 100 * h
			f.layers[2] = append(f.layers[2], spot{s, col})
		}
	case PulseInner:
		add(0, "PULSE_RING_SHAPES", "PULSE_RING_COLORS_"+variant, 0)
		add(1, "PULSE_INNER_SHAPES", "PULSE_INNER_COLORS_"+variant, 0)
		add(2, "PULSE_INNER_BLOOM_SHAPES", "PULSE_INNER_BLOOM_COLORS_"+variant, 2)
	case PulseOutside:
		add(0, "PULSE_OUTER_CORE_SHAPES", "PULSE_OUTER_CORE_COLORS_"+variant, 1)
		add(1, "PULSE_OUTER_CORE_SHAPES", "PULSE_OUTER_CORE_COLORS_"+variant, 1)
		add(2, "PULSE_OUTER_BLOOM_SHAPES", "PULSE_OUTER_BLOOM_COLORS_"+variant, 2)
	}
	return f
}

func (f colorField) pulseSpot(s, c vector, i, kind int) (vector, vector) {
	inner := f.mode == PulseInner
	drift, op, gh := 19., 0., .58
	if f.dark {
		drift, op, gh = 14, .46, .16
	}
	if inner {
		drift, op, gh = 40, .45, .22
		if f.dark {
			drift, op, gh = 33, .48, .34
		}
	}
	bs, ss, ghs := 2.3, 6.4, 3.8
	if inner {
		bs, ss, ghs = 1.9, 2.6, 2.4
	}
	scale := f.duration / 2.3
	bs *= scale
	ss *= scale
	ghs *= scale
	osc := func(a, b, p, d float64) float64 { return oscillate(a, b, p, d, f.seconds) }
	sp := .28
	regions := []int{1, 2, 3, 1, 2, 3, 1, 2, 3}
	quads := []int{0, 0, 2, 2, 3, 3, 1, 1, 1}
	if kind == 1 {
		regions = []int{1, 2, 3, 1, 2, 3, 1, 2}
		quads = []int{0, 1, 1, 3, 3, 2, 2, 0}
	}
	region, quad := regions[i%len(regions)], quads[i%len(quads)]
	wx, hy, dx, dy := 0., 0., 0., 0.
	switch region {
	case 1:
		wx = osc(1-sp, 1+sp*1.1, ss*.9, 0)
		hy = osc(1+sp*.9, 1-sp*.85, ss*1.26, 0)
		dx = osc(-drift, drift*.9, bs*1.6, 0)
		dy = osc(drift*.55, -drift*.7, bs*1.6, 0)
	case 2:
		wx = osc(1+sp, 1-sp*.85, ss*1.1, 0)
		hy = osc(1-sp*.8, 1+sp*1.05, ss*.81, 0)
		dx = osc(drift*.8, -drift*.9, bs*1.88, 0)
		dy = osc(-drift, drift*.65, bs*1.88, 0)
	case 3:
		wx = osc(1-sp*.6, 1+sp*1.15, ss*.98, 0)
		hy = osc(1+sp*.75, 1-sp, ss*1.4, 0)
		dx = osc(-drift*.6, drift, bs*1.45, 0)
		dy = osc(-drift*.85, drift*.45, bs*1.45, 0)
	}
	opacity := 1.
	switch quad {
	case 0:
		opacity = 1
	case 1:
		opacity = osc(1-op, 1, bs*1.32, bs*.28)
	case 2:
		opacity = osc(1-op, 1, bs*.84, bs*.55)
	case 3:
		opacity = osc(1-op, 1, bs*1.58, bs*.83)
	}
	sx, sy := 1., 1.
	if !inner {
		sx = max(.35, min(4, f.width/350))
		sy = max(.35, min(4, f.height/140))
	}
	s[0] += dx
	s[1] += dy
	s[2] *= wx * sx
	s[3] *= hy * osc(1-gh, 1+gh, ghs, 0) * sy
	c[3] *= opacity
	return s, c
}

func lineTravel(t float64) float64 {
	points := [...]float64{.06, .15, .25, .35, .44, .5, .56, .65, .75, .85, .94}
	u := clamp(t) * 10
	i := min(9, int(u))
	return mix(points[i], points[i+1], u-float64(i))
}
func lineEdge(t float64) float64 { return min(clamp((t-.125)/.2), clamp((.875-t)/.2)) }

// sample evaluates the unblurred spatial field using the WGSL over operator.
func (f colorField) sample(layer int, x, y float64) vector {
	var a vector
	for _, s := range f.layers[layer] {
		dx, dy := (x-s.shape[0])/max(.0001, s.shape[2]), (y-s.shape[1])/max(.0001, s.shape[3])
		if dx*dx+dy*dy >= 1 {
			continue
		}
		weight := (1 - math.Sqrt(dx*dx+dy*dy)) * s.color[3] * (1 - a[3])
		for k := 0; k < 3; k++ {
			a[k] += s.color[k] * weight
		}
		a[3] += weight
	}
	return a
}
func hueRotate(c vector, co, si, brightness, saturation float64) vector {
	r := (.213+.787*co-.213*si)*c[0] + (.715-.715*co-.715*si)*c[1] + (.072-.072*co+.928*si)*c[2]
	g := (.213-.213*co+.143*si)*c[0] + (.715+.285*co+.140*si)*c[1] + (.072-.072*co-.283*si)*c[2]
	b := (.213-.213*co-.787*si)*c[0] + (.715-.715*co+.715*si)*c[1] + (.072+.928*co+.072*si)*c[2]
	luma := .213*r + .715*g + .072*b
	return vector{clamp(mix(luma, r, saturation) * brightness), clamp(mix(luma, g, saturation) * brightness), clamp(mix(luma, b, saturation) * brightness), c[3]}
}
