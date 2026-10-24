struct Params {
  time: f32,
  variant: f32,
  hueRange: f32,
  brightness: f32,
  saturation: f32,
  mode: f32,
  dark: f32,
  width: f32,
  height: f32,
  rootWidth: f32,
  rootHeight: f32,
  duration: f32,
  staticColors: f32,
  pulseScaleX: f32,
  pulseScaleY: f32,
  offsetX: f32,
  offsetY: f32,
  layer: f32,
}

@group(0) @binding(0) var<uniform> params: Params;

const MD_SHAPES = array<vec4f, 9>(vec4f(0.33000, -0.07400, 70.000, 40.000), vec4f(0.12000, -0.05000, 60.000, 35.000), vec4f(0.02100, 0.68300, 40.000, 70.000), vec4f(0.02100, 0.68300, 20.000, 35.000), vec4f(0.74400, 1.00000, 180.000, 32.000), vec4f(0.55000, 1.00000, 85.000, 26.000), vec4f(0.93900, 0.00000, 74.000, 32.000), vec4f(1.00000, 0.27100, 26.000, 42.000), vec4f(1.00000, 0.27100, 52.000, 48.000));
const SM_SHAPES = array<vec4f, 8>(vec4f(0.02000, 0.68000, 9.000, 18.000), vec4f(0.02000, 0.68000, 4.000, 8.000), vec4f(0.72000, -0.03000, 59.000, 9.000), vec4f(0.74000, 1.00000, 42.000, 7.000), vec4f(1.00000, 0.27000, 10.000, 17.000), vec4f(1.00000, 0.27000, 10.000, 18.000), vec4f(1.00000, 0.27000, 5.000, 10.000), vec4f(1.00000, 0.27000, 11.000, 12.000));
const LINE_SHAPES_DARK = array<vec4f, 9>(vec4f(0.000, 2.000, 36.000, 36.000), vec4f(39.000, 0.000, 30.000, 32.000), vec4f(-36.000, 2.000, 33.000, 28.000), vec4f(-54.000, 0.000, 29.000, 34.000), vec4f(51.000, -1.000, 27.000, 30.000), vec4f(21.000, 1.000, 36.000, 24.000), vec4f(-21.000, 0.000, 30.000, 22.000), vec4f(66.000, 1.000, 25.000, 28.000), vec4f(-66.000, -1.000, 23.000, 30.000));
const LINE_SHAPES_LIGHT = array<vec4f, 9>(vec4f(0.000, 2.000, 45.000, 36.000), vec4f(65.000, 0.000, 35.000, 32.000), vec4f(-60.000, 2.000, 40.000, 28.000), vec4f(-90.000, 0.000, 35.000, 34.000), vec4f(85.000, -1.000, 38.000, 30.000), vec4f(35.000, 1.000, 50.000, 24.000), vec4f(-35.000, 0.000, 40.000, 22.000), vec4f(110.000, 1.000, 35.000, 28.000), vec4f(-110.000, -1.000, 30.000, 30.000));
const LINE_INNER_SHAPES = array<vec4f, 9>(vec4f(0.000, 0.000, 33.000, 30.000), vec4f(39.000, -3.000, 24.000, 26.000), vec4f(-36.000, 0.000, 27.000, 24.000), vec4f(-54.000, -2.000, 23.000, 28.000), vec4f(51.000, -1.000, 24.000, 24.000), vec4f(21.000, 0.000, 30.000, 20.000), vec4f(-21.000, -2.000, 25.000, 18.000), vec4f(66.000, 0.000, 21.000, 24.000), vec4f(-66.000, -1.000, 18.000, 26.000));
const LINE_BLOOM_SHAPES = array<vec4f, 8>(vec4f(8.000, 98.000, 1.000, 92.000), vec4f(22.000, 96.000, 10.000, 35.000), vec4f(36.000, 97.000, 2.000, 72.000), vec4f(50.000, 98.000, 14.000, 28.000), vec4f(64.000, 96.000, 1.200, 85.000), vec4f(78.000, 98.000, 7.000, 45.000), vec4f(92.000, 97.000, 1.000, 60.000), vec4f(50.000, 100.000, 50.000, 32.000));
const PULSE_RING_SHAPES = array<vec4f, 9>(vec4f(0.33000, -0.07400, 70.000, 40.000), vec4f(0.12000, -0.05000, 60.000, 35.000), vec4f(0.02100, 0.68300, 40.000, 70.000), vec4f(0.02100, 0.68300, 20.000, 35.000), vec4f(0.74400, 1.00000, 180.000, 32.000), vec4f(0.55000, 1.00000, 85.000, 26.000), vec4f(0.93900, 0.00000, 74.000, 32.000), vec4f(1.00000, 0.27100, 26.000, 42.000), vec4f(1.00000, 0.27100, 52.000, 48.000));
const PULSE_INNER_SHAPES = array<vec4f, 9>(vec4f(0.33000, -0.07400, 65.000, 35.000), vec4f(0.12000, -0.05000, 55.000, 30.000), vec4f(0.02100, 0.68300, 35.000, 65.000), vec4f(0.02100, 0.68300, 15.000, 30.000), vec4f(0.74400, 1.00000, 173.000, 28.000), vec4f(0.55000, 1.00000, 80.000, 22.000), vec4f(0.93900, 0.00000, 69.000, 28.000), vec4f(1.00000, 0.27100, 22.000, 38.000), vec4f(1.00000, 0.27100, 47.000, 44.000));
const PULSE_INNER_BLOOM_SHAPES = array<vec4f, 7>(vec4f(0.33000, -0.07400, 84.000, 48.000), vec4f(0.12000, -0.05000, 72.000, 42.000), vec4f(0.02100, 0.68300, 48.000, 84.000), vec4f(0.74400, 1.00000, 216.000, 38.000), vec4f(0.55000, 1.00000, 102.000, 31.000), vec4f(0.93900, 0.00000, 89.000, 38.000), vec4f(1.00000, 0.27100, 62.000, 58.000));
const PULSE_OUTER_CORE_SHAPES = array<vec4f, 8>(vec4f(0.27000, 0.00000, 80.000, 19.000), vec4f(0.73000, -0.01000, 74.000, 11.000), vec4f(1.00000, 0.33000, 15.000, 44.000), vec4f(1.01000, 0.72000, 19.000, 38.000), vec4f(0.67000, 1.00000, 84.000, 13.000), vec4f(0.24000, 1.01000, 60.000, 21.000), vec4f(0.00000, 0.60000, 17.000, 40.000), vec4f(-0.01000, 0.28000, 13.000, 32.000));
const PULSE_OUTER_BLOOM_SHAPES = array<vec4f, 7>(vec4f(0.27000, 0.03000, 110.000, 30.000), vec4f(0.73000, 0.01000, 100.000, 20.000), vec4f(1.00000, 0.33000, 26.000, 62.000), vec4f(1.01000, 0.72000, 30.000, 56.000), vec4f(0.67000, 0.99000, 120.000, 22.000), vec4f(0.24000, 0.99000, 88.000, 32.000), vec4f(0.00000, 0.60000, 28.000, 58.000));
const MD_COLORS_COLORFUL = array<vec4f, 9>(vec4f(1.0000, 0.1961, 0.3922, 1.0000), vec4f(0.1569, 0.5490, 1.0000, 1.0000), vec4f(0.1961, 0.7843, 0.3137, 1.0000), vec4f(0.1176, 0.7255, 0.6667, 1.0000), vec4f(0.3922, 0.2745, 1.0000, 1.0000), vec4f(0.1569, 0.5490, 1.0000, 1.0000), vec4f(1.0000, 0.4706, 0.1569, 1.0000), vec4f(0.9412, 0.1961, 0.7059, 1.0000), vec4f(0.7059, 0.1569, 0.9412, 1.0000));
const SM_COLORS_COLORFUL = array<vec4f, 8>(vec4f(0.1961, 0.7843, 0.3137, 1.0000), vec4f(0.1176, 0.7255, 0.6667, 1.0000), vec4f(1.0000, 0.4706, 0.1569, 1.0000), vec4f(0.3922, 0.2745, 1.0000, 1.0000), vec4f(0.9412, 0.1961, 0.7059, 1.0000), vec4f(0.7059, 0.1569, 0.9412, 1.0000), vec4f(0.1569, 0.5490, 1.0000, 1.0000), vec4f(1.0000, 0.1961, 0.3922, 1.0000));
const MD_INNER_COLORS_COLORFUL = array<vec4f, 9>(vec4f(1.0000, 0.1961, 0.3922, 0.4500), vec4f(0.1569, 0.5490, 1.0000, 0.4500), vec4f(0.1961, 0.7843, 0.3137, 0.4500), vec4f(0.1176, 0.7255, 0.6667, 0.4500), vec4f(0.3922, 0.2745, 1.0000, 0.4500), vec4f(0.1569, 0.5490, 1.0000, 0.4500), vec4f(1.0000, 0.4706, 0.1569, 0.4500), vec4f(0.9412, 0.1961, 0.7059, 0.4500), vec4f(0.7059, 0.1569, 0.9412, 0.4500));
const SM_INNER_COLORS_COLORFUL = array<vec4f, 8>(vec4f(0.1961, 0.7843, 0.3137, 0.5000), vec4f(0.1176, 0.7255, 0.6667, 0.4500), vec4f(1.0000, 0.4706, 0.1569, 0.3500), vec4f(0.3922, 0.2745, 1.0000, 0.3500), vec4f(0.9412, 0.1961, 0.7059, 0.3000), vec4f(0.7059, 0.1569, 0.9412, 0.4000), vec4f(0.1569, 0.5490, 1.0000, 0.3000), vec4f(1.0000, 0.1961, 0.3922, 0.3000));
const LINE_COLORS_COLORFUL_DARK = array<vec4f, 9>(vec4f(1.0000, 0.1961, 0.3922, 1.0000), vec4f(0.1569, 0.7059, 0.8627, 1.0000), vec4f(0.1961, 0.7843, 0.3137, 1.0000), vec4f(0.7059, 0.1569, 0.9412, 1.0000), vec4f(1.0000, 0.6275, 0.1176, 1.0000), vec4f(0.3922, 0.2745, 1.0000, 1.0000), vec4f(0.1569, 0.5490, 1.0000, 1.0000), vec4f(0.9412, 0.1961, 0.7059, 1.0000), vec4f(0.1176, 0.7255, 0.6667, 1.0000));
const LINE_COLORS_COLORFUL_LIGHT = array<vec4f, 9>(vec4f(1.0000, 0.1961, 0.3922, 1.0000), vec4f(0.1569, 0.5490, 1.0000, 1.0000), vec4f(0.1961, 0.7843, 0.3137, 1.0000), vec4f(0.7059, 0.1569, 0.9412, 1.0000), vec4f(0.1176, 0.7255, 0.6667, 1.0000), vec4f(0.3922, 0.2745, 1.0000, 1.0000), vec4f(0.1569, 0.5490, 1.0000, 1.0000), vec4f(1.0000, 0.4706, 0.1569, 1.0000), vec4f(0.9412, 0.1961, 0.7059, 1.0000));
const LINE_INNER_COLORS_COLORFUL = array<vec4f, 9>(vec4f(1.0000, 0.1961, 0.3922, 0.4800), vec4f(0.1569, 0.7059, 0.8627, 0.4200), vec4f(0.1961, 0.7843, 0.3137, 0.4800), vec4f(0.7059, 0.1569, 0.9412, 0.4200), vec4f(1.0000, 0.6275, 0.1176, 0.5000), vec4f(0.3922, 0.2745, 1.0000, 0.4500), vec4f(0.1569, 0.5490, 1.0000, 0.4000), vec4f(0.9412, 0.1961, 0.7059, 0.4500), vec4f(0.1176, 0.7255, 0.6667, 0.5200));
const LINE_BLOOM_HEAD_COLORFUL_DARK = array<vec4f, 2>(vec4f(1.0000, 0.2353, 0.3137, 1.0000), vec4f(0.1569, 0.7451, 0.7059, 0.9800));
const LINE_BLOOM_HEAD_COLORFUL_LIGHT = array<vec4f, 2>(vec4f(0.7843, 0.1176, 0.2353, 1.0000), vec4f(0.0784, 0.5882, 0.5490, 1.0000));
const LINE_BLOOM_SPIKES_COLORFUL_DARK = array<vec4f, 5>(vec4f(0.3922, 0.2745, 1.0000, 1.0000), vec4f(1.0000, 0.6667, 0.1569, 0.5900), vec4f(0.1961, 0.7843, 0.3922, 1.0000), vec4f(0.7843, 0.1961, 0.9412, 0.9100), vec4f(0.1569, 0.5490, 1.0000, 1.0000));
const LINE_BLOOM_SPIKES_COLORFUL_LIGHT = array<vec4f, 5>(vec4f(0.3137, 0.1961, 0.7843, 1.0000), vec4f(0.8235, 0.5098, 0.0000, 0.7000), vec4f(0.1176, 0.6275, 0.2745, 1.0000), vec4f(0.6275, 0.1176, 0.7451, 1.0000), vec4f(0.1176, 0.3922, 0.7843, 1.0000));
const PULSE_RING_COLORS_COLORFUL = array<vec4f, 9>(vec4f(1.0000, 0.1961, 0.3922, 1.0000), vec4f(0.1569, 0.5490, 1.0000, 1.0000), vec4f(0.1961, 0.7843, 0.3137, 1.0000), vec4f(0.1176, 0.7255, 0.6667, 1.0000), vec4f(0.3922, 0.2745, 1.0000, 1.0000), vec4f(0.1569, 0.5490, 1.0000, 1.0000), vec4f(1.0000, 0.4706, 0.1569, 1.0000), vec4f(0.9412, 0.1961, 0.7059, 1.0000), vec4f(0.7059, 0.1569, 0.9412, 1.0000));
const PULSE_INNER_COLORS_COLORFUL = array<vec4f, 9>(vec4f(1.0000, 0.1961, 0.3922, 0.4500), vec4f(0.1569, 0.5490, 1.0000, 0.4500), vec4f(0.1961, 0.7843, 0.3137, 0.4500), vec4f(0.1176, 0.7255, 0.6667, 0.4500), vec4f(0.3922, 0.2745, 1.0000, 0.4500), vec4f(0.1569, 0.5490, 1.0000, 0.4500), vec4f(1.0000, 0.4706, 0.1569, 0.4500), vec4f(0.9412, 0.1961, 0.7059, 0.4500), vec4f(0.7059, 0.1569, 0.9412, 0.4500));
const PULSE_INNER_BLOOM_COLORS_COLORFUL = array<vec4f, 7>(vec4f(1.0000, 0.1961, 0.3922, 0.7750), vec4f(0.1569, 0.5490, 1.0000, 0.7750), vec4f(0.1961, 0.7843, 0.3137, 0.7750), vec4f(0.3922, 0.2745, 1.0000, 0.7750), vec4f(0.1569, 0.5490, 1.0000, 0.7750), vec4f(1.0000, 0.4706, 0.1569, 0.7750), vec4f(0.7059, 0.1569, 0.9412, 0.7750));
const PULSE_OUTER_CORE_COLORS_COLORFUL = array<vec4f, 8>(vec4f(1.0000, 0.1961, 0.3922, 1.0000), vec4f(1.0000, 0.4706, 0.1569, 1.0000), vec4f(0.9412, 0.1961, 0.7059, 1.0000), vec4f(0.7059, 0.1569, 0.9412, 1.0000), vec4f(0.3922, 0.2745, 1.0000, 1.0000), vec4f(0.1569, 0.5490, 1.0000, 1.0000), vec4f(0.1961, 0.7843, 0.3137, 1.0000), vec4f(0.1176, 0.7255, 0.6667, 1.0000));
const PULSE_OUTER_BLOOM_COLORS_COLORFUL = array<vec4f, 7>(vec4f(1.0000, 0.1961, 0.3922, 0.7700), vec4f(1.0000, 0.4706, 0.1569, 0.7700), vec4f(0.9412, 0.1961, 0.7059, 0.7700), vec4f(0.7059, 0.1569, 0.9412, 0.7700), vec4f(0.3922, 0.2745, 1.0000, 0.7700), vec4f(0.1569, 0.5490, 1.0000, 0.7700), vec4f(0.1961, 0.7843, 0.3137, 0.7700));
const MD_COLORS_MONO = array<vec4f, 9>(vec4f(0.7059, 0.7059, 0.7059, 1.0000), vec4f(0.5490, 0.5490, 0.5490, 1.0000), vec4f(0.6275, 0.6275, 0.6275, 1.0000), vec4f(0.5098, 0.5098, 0.5098, 1.0000), vec4f(0.6667, 0.6667, 0.6667, 1.0000), vec4f(0.5882, 0.5882, 0.5882, 1.0000), vec4f(0.7451, 0.7451, 0.7451, 1.0000), vec4f(0.5686, 0.5686, 0.5686, 1.0000), vec4f(0.6471, 0.6471, 0.6471, 1.0000));
const SM_COLORS_MONO = array<vec4f, 8>(vec4f(0.6275, 0.6275, 0.6275, 1.0000), vec4f(0.5490, 0.5490, 0.5490, 1.0000), vec4f(0.7059, 0.7059, 0.7059, 1.0000), vec4f(0.5882, 0.5882, 0.5882, 1.0000), vec4f(0.6667, 0.6667, 0.6667, 1.0000), vec4f(0.6078, 0.6078, 0.6078, 1.0000), vec4f(0.5686, 0.5686, 0.5686, 1.0000), vec4f(0.6471, 0.6471, 0.6471, 1.0000));
const MD_INNER_COLORS_MONO = array<vec4f, 9>(vec4f(0.7059, 0.7059, 0.7059, 0.2250), vec4f(0.5490, 0.5490, 0.5490, 0.2250), vec4f(0.6275, 0.6275, 0.6275, 0.2250), vec4f(0.5098, 0.5098, 0.5098, 0.2250), vec4f(0.6667, 0.6667, 0.6667, 0.2250), vec4f(0.5882, 0.5882, 0.5882, 0.2250), vec4f(0.7451, 0.7451, 0.7451, 0.2250), vec4f(0.5686, 0.5686, 0.5686, 0.2250), vec4f(0.6471, 0.6471, 0.6471, 0.2250));
const SM_INNER_COLORS_MONO = array<vec4f, 8>(vec4f(0.6275, 0.6275, 0.6275, 0.2500), vec4f(0.5490, 0.5490, 0.5490, 0.2200), vec4f(0.7059, 0.7059, 0.7059, 0.1700), vec4f(0.5882, 0.5882, 0.5882, 0.1700), vec4f(0.6667, 0.6667, 0.6667, 0.1500), vec4f(0.6078, 0.6078, 0.6078, 0.2000), vec4f(0.5686, 0.5686, 0.5686, 0.1500), vec4f(0.6471, 0.6471, 0.6471, 0.1500));
const LINE_COLORS_MONO_DARK = array<vec4f, 9>(vec4f(0.7843, 0.7843, 0.7843, 1.0000), vec4f(0.6667, 0.6667, 0.6667, 1.0000), vec4f(0.6078, 0.6078, 0.6078, 1.0000), vec4f(0.7255, 0.7255, 0.7255, 1.0000), vec4f(0.6471, 0.6471, 0.6471, 1.0000), vec4f(0.7059, 0.7059, 0.7059, 1.0000), vec4f(0.6275, 0.6275, 0.6275, 1.0000), vec4f(0.6863, 0.6863, 0.6863, 1.0000), vec4f(0.7451, 0.7451, 0.7451, 1.0000));
const LINE_COLORS_MONO_LIGHT = array<vec4f, 9>(vec4f(0.3922, 0.3922, 0.3922, 1.0000), vec4f(0.3137, 0.3137, 0.3137, 1.0000), vec4f(0.3529, 0.3529, 0.3529, 1.0000), vec4f(0.2745, 0.2745, 0.2745, 1.0000), vec4f(0.3333, 0.3333, 0.3333, 1.0000), vec4f(0.3725, 0.3725, 0.3725, 1.0000), vec4f(0.2941, 0.2941, 0.2941, 1.0000), vec4f(0.4118, 0.4118, 0.4118, 1.0000), vec4f(0.2549, 0.2549, 0.2549, 1.0000));
const LINE_INNER_COLORS_MONO = array<vec4f, 9>(vec4f(0.7843, 0.7843, 0.7843, 0.4800), vec4f(0.6667, 0.6667, 0.6667, 0.4200), vec4f(0.6078, 0.6078, 0.6078, 0.4800), vec4f(0.7255, 0.7255, 0.7255, 0.4200), vec4f(0.6471, 0.6471, 0.6471, 0.5000), vec4f(0.7059, 0.7059, 0.7059, 0.4500), vec4f(0.6275, 0.6275, 0.6275, 0.4000), vec4f(0.6863, 0.6863, 0.6863, 0.4500), vec4f(0.7451, 0.7451, 0.7451, 0.5200));
const LINE_BLOOM_HEAD_MONO_DARK = array<vec4f, 2>(vec4f(0.7843, 0.7843, 0.7843, 1.0000), vec4f(0.6667, 0.6667, 0.6667, 1.0000));
const LINE_BLOOM_HEAD_MONO_LIGHT = array<vec4f, 2>(vec4f(0.3137, 0.3137, 0.3137, 1.0000), vec4f(0.4706, 0.4706, 0.4706, 1.0000));
const LINE_BLOOM_SPIKES_MONO_DARK = array<vec4f, 5>(vec4f(0.7843, 0.7843, 0.7843, 1.0000), vec4f(0.7059, 0.7059, 0.7059, 0.5900), vec4f(0.7451, 0.7451, 0.7451, 1.0000), vec4f(0.6667, 0.6667, 0.6667, 0.9100), vec4f(0.7255, 0.7255, 0.7255, 1.0000));
const LINE_BLOOM_SPIKES_MONO_LIGHT = array<vec4f, 5>(vec4f(0.3137, 0.3137, 0.3137, 1.0000), vec4f(0.3922, 0.3922, 0.3922, 0.7000), vec4f(0.2745, 0.2745, 0.2745, 1.0000), vec4f(0.3529, 0.3529, 0.3529, 1.0000), vec4f(0.3333, 0.3333, 0.3333, 1.0000));
const PULSE_RING_COLORS_MONO = array<vec4f, 9>(vec4f(0.7059, 0.7059, 0.7059, 1.0000), vec4f(0.5490, 0.5490, 0.5490, 1.0000), vec4f(0.6275, 0.6275, 0.6275, 1.0000), vec4f(0.5098, 0.5098, 0.5098, 1.0000), vec4f(0.6667, 0.6667, 0.6667, 1.0000), vec4f(0.5882, 0.5882, 0.5882, 1.0000), vec4f(0.7451, 0.7451, 0.7451, 1.0000), vec4f(0.5686, 0.5686, 0.5686, 1.0000), vec4f(0.6471, 0.6471, 0.6471, 1.0000));
const PULSE_INNER_COLORS_MONO = array<vec4f, 9>(vec4f(0.7059, 0.7059, 0.7059, 0.2250), vec4f(0.5490, 0.5490, 0.5490, 0.2250), vec4f(0.6275, 0.6275, 0.6275, 0.2250), vec4f(0.5098, 0.5098, 0.5098, 0.2250), vec4f(0.6667, 0.6667, 0.6667, 0.2250), vec4f(0.5882, 0.5882, 0.5882, 0.2250), vec4f(0.7451, 0.7451, 0.7451, 0.2250), vec4f(0.5686, 0.5686, 0.5686, 0.2250), vec4f(0.6471, 0.6471, 0.6471, 0.2250));
const PULSE_INNER_BLOOM_COLORS_MONO = array<vec4f, 7>(vec4f(0.7059, 0.7059, 0.7059, 0.7750), vec4f(0.5490, 0.5490, 0.5490, 0.7750), vec4f(0.6275, 0.6275, 0.6275, 0.7750), vec4f(0.6667, 0.6667, 0.6667, 0.7750), vec4f(0.5882, 0.5882, 0.5882, 0.7750), vec4f(0.7451, 0.7451, 0.7451, 0.7750), vec4f(0.6471, 0.6471, 0.6471, 0.7750));
const PULSE_OUTER_CORE_COLORS_MONO = array<vec4f, 8>(vec4f(0.7059, 0.7059, 0.7059, 1.0000), vec4f(0.7451, 0.7451, 0.7451, 1.0000), vec4f(0.5686, 0.5686, 0.5686, 1.0000), vec4f(0.6471, 0.6471, 0.6471, 1.0000), vec4f(0.6667, 0.6667, 0.6667, 1.0000), vec4f(0.5490, 0.5490, 0.5490, 1.0000), vec4f(0.6275, 0.6275, 0.6275, 1.0000), vec4f(0.5098, 0.5098, 0.5098, 1.0000));
const PULSE_OUTER_BLOOM_COLORS_MONO = array<vec4f, 7>(vec4f(0.7059, 0.7059, 0.7059, 0.7700), vec4f(0.7451, 0.7451, 0.7451, 0.7700), vec4f(0.5686, 0.5686, 0.5686, 0.7700), vec4f(0.6471, 0.6471, 0.6471, 0.7700), vec4f(0.6667, 0.6667, 0.6667, 0.7700), vec4f(0.5490, 0.5490, 0.5490, 0.7700), vec4f(0.6275, 0.6275, 0.6275, 0.7700));
const MD_COLORS_OCEAN = array<vec4f, 9>(vec4f(0.3922, 0.3137, 0.8627, 1.0000), vec4f(0.2353, 0.4706, 1.0000, 1.0000), vec4f(0.3137, 0.3922, 0.7843, 1.0000), vec4f(0.1961, 0.5490, 0.8627, 1.0000), vec4f(0.4706, 0.3137, 1.0000, 1.0000), vec4f(0.2745, 0.5098, 1.0000, 1.0000), vec4f(0.5490, 0.3922, 0.9412, 1.0000), vec4f(0.3529, 0.4314, 0.9020, 1.0000), vec4f(0.5098, 0.2745, 1.0000, 1.0000));
const SM_COLORS_OCEAN = array<vec4f, 8>(vec4f(0.2353, 0.5490, 0.7843, 1.0000), vec4f(0.1961, 0.4706, 0.7059, 1.0000), vec4f(0.3922, 0.3137, 0.8627, 1.0000), vec4f(0.3137, 0.3922, 1.0000, 1.0000), vec4f(0.4706, 0.2745, 0.9412, 1.0000), vec4f(0.3529, 0.3137, 0.8627, 1.0000), vec4f(0.2745, 0.4314, 1.0000, 1.0000), vec4f(0.4314, 0.3529, 0.9020, 1.0000));
const MD_INNER_COLORS_OCEAN = array<vec4f, 9>(vec4f(0.3922, 0.3137, 0.8627, 0.4500), vec4f(0.2353, 0.4706, 1.0000, 0.4500), vec4f(0.3137, 0.3922, 0.7843, 0.4500), vec4f(0.1961, 0.5490, 0.8627, 0.4500), vec4f(0.4706, 0.3137, 1.0000, 0.4500), vec4f(0.2745, 0.5098, 1.0000, 0.4500), vec4f(0.5490, 0.3922, 0.9412, 0.4500), vec4f(0.3529, 0.4314, 0.9020, 0.4500), vec4f(0.5098, 0.2745, 1.0000, 0.4500));
const SM_INNER_COLORS_OCEAN = array<vec4f, 8>(vec4f(0.2353, 0.5490, 0.7843, 0.5000), vec4f(0.1961, 0.4706, 0.7059, 0.4500), vec4f(0.3922, 0.3137, 0.8627, 0.3500), vec4f(0.3137, 0.3922, 1.0000, 0.3500), vec4f(0.4706, 0.2745, 0.9412, 0.3000), vec4f(0.3529, 0.3137, 0.8627, 0.4000), vec4f(0.2745, 0.4314, 1.0000, 0.3000), vec4f(0.4314, 0.3529, 0.9020, 0.3000));
const LINE_COLORS_OCEAN_DARK = array<vec4f, 9>(vec4f(0.3922, 0.3137, 0.8627, 1.0000), vec4f(0.2353, 0.4706, 1.0000, 1.0000), vec4f(0.3137, 0.3922, 0.7843, 1.0000), vec4f(0.5098, 0.2745, 1.0000, 1.0000), vec4f(0.2745, 0.5098, 1.0000, 1.0000), vec4f(0.4706, 0.3137, 1.0000, 1.0000), vec4f(0.3529, 0.4314, 0.9020, 1.0000), vec4f(0.4314, 0.3529, 0.9412, 1.0000), vec4f(0.5490, 0.3922, 1.0000, 1.0000));
const LINE_COLORS_OCEAN_LIGHT = array<vec4f, 9>(vec4f(0.3137, 0.2353, 0.7843, 1.0000), vec4f(0.1961, 0.3922, 0.8627, 1.0000), vec4f(0.2745, 0.3529, 0.7451, 1.0000), vec4f(0.4314, 0.2353, 0.8627, 1.0000), vec4f(0.2353, 0.4314, 0.9020, 1.0000), vec4f(0.3922, 0.2745, 0.9412, 1.0000), vec4f(0.3137, 0.3922, 0.8235, 1.0000), vec4f(0.3529, 0.3137, 0.8824, 1.0000), vec4f(0.4706, 0.3529, 0.9608, 1.0000));
const LINE_INNER_COLORS_OCEAN = array<vec4f, 9>(vec4f(0.3922, 0.3137, 0.8627, 0.4800), vec4f(0.2353, 0.4706, 1.0000, 0.4200), vec4f(0.3137, 0.3922, 0.7843, 0.4800), vec4f(0.5098, 0.2745, 1.0000, 0.4200), vec4f(0.2745, 0.5098, 1.0000, 0.5000), vec4f(0.4706, 0.3137, 1.0000, 0.4500), vec4f(0.3529, 0.4314, 0.9020, 0.4000), vec4f(0.4314, 0.3529, 0.9412, 0.4500), vec4f(0.5490, 0.3922, 1.0000, 0.5200));
const LINE_BLOOM_HEAD_OCEAN_DARK = array<vec4f, 2>(vec4f(0.3922, 0.4706, 1.0000, 1.0000), vec4f(0.5098, 0.3922, 0.8627, 0.9800));
const LINE_BLOOM_HEAD_OCEAN_LIGHT = array<vec4f, 2>(vec4f(0.2353, 0.2353, 0.7059, 1.0000), vec4f(0.3137, 0.3922, 0.7843, 1.0000));
const LINE_BLOOM_SPIKES_OCEAN_DARK = array<vec4f, 5>(vec4f(0.3922, 0.3137, 1.0000, 1.0000), vec4f(0.3137, 0.5098, 0.8627, 0.5900), vec4f(0.2353, 0.3922, 1.0000, 1.0000), vec4f(0.3529, 0.4706, 0.7843, 0.9100), vec4f(0.4706, 0.3529, 1.0000, 1.0000));
const LINE_BLOOM_SPIKES_OCEAN_LIGHT = array<vec4f, 5>(vec4f(0.1961, 0.1569, 0.7059, 1.0000), vec4f(0.1569, 0.3137, 0.7843, 0.7000), vec4f(0.1176, 0.1961, 0.7451, 1.0000), vec4f(0.2353, 0.3529, 0.7059, 1.0000), vec4f(0.2745, 0.2353, 0.7843, 1.0000));
const PULSE_RING_COLORS_OCEAN = array<vec4f, 9>(vec4f(0.3922, 0.3137, 0.8627, 1.0000), vec4f(0.2353, 0.4706, 1.0000, 1.0000), vec4f(0.3137, 0.3922, 0.7843, 1.0000), vec4f(0.1961, 0.5490, 0.8627, 1.0000), vec4f(0.4706, 0.3137, 1.0000, 1.0000), vec4f(0.2745, 0.5098, 1.0000, 1.0000), vec4f(0.5490, 0.3922, 0.9412, 1.0000), vec4f(0.3529, 0.4314, 0.9020, 1.0000), vec4f(0.5098, 0.2745, 1.0000, 1.0000));
const PULSE_INNER_COLORS_OCEAN = array<vec4f, 9>(vec4f(0.3922, 0.3137, 0.8627, 0.4500), vec4f(0.2353, 0.4706, 1.0000, 0.4500), vec4f(0.3137, 0.3922, 0.7843, 0.4500), vec4f(0.1961, 0.5490, 0.8627, 0.4500), vec4f(0.4706, 0.3137, 1.0000, 0.4500), vec4f(0.2745, 0.5098, 1.0000, 0.4500), vec4f(0.5490, 0.3922, 0.9412, 0.4500), vec4f(0.3529, 0.4314, 0.9020, 0.4500), vec4f(0.5098, 0.2745, 1.0000, 0.4500));
const PULSE_INNER_BLOOM_COLORS_OCEAN = array<vec4f, 7>(vec4f(0.3922, 0.3137, 0.8627, 0.7750), vec4f(0.2353, 0.4706, 1.0000, 0.7750), vec4f(0.3137, 0.3922, 0.7843, 0.7750), vec4f(0.4706, 0.3137, 1.0000, 0.7750), vec4f(0.2745, 0.5098, 1.0000, 0.7750), vec4f(0.5490, 0.3922, 0.9412, 0.7750), vec4f(0.5098, 0.2745, 1.0000, 0.7750));
const PULSE_OUTER_CORE_COLORS_OCEAN = array<vec4f, 8>(vec4f(0.3922, 0.3137, 0.8627, 1.0000), vec4f(0.5490, 0.3922, 0.9412, 1.0000), vec4f(0.3529, 0.4314, 0.9020, 1.0000), vec4f(0.5098, 0.2745, 1.0000, 1.0000), vec4f(0.4706, 0.3137, 1.0000, 1.0000), vec4f(0.2353, 0.4706, 1.0000, 1.0000), vec4f(0.3137, 0.3922, 0.7843, 1.0000), vec4f(0.1961, 0.5490, 0.8627, 1.0000));
const PULSE_OUTER_BLOOM_COLORS_OCEAN = array<vec4f, 7>(vec4f(0.3922, 0.3137, 0.8627, 0.7700), vec4f(0.5490, 0.3922, 0.9412, 0.7700), vec4f(0.3529, 0.4314, 0.9020, 0.7700), vec4f(0.5098, 0.2745, 1.0000, 0.7700), vec4f(0.4706, 0.3137, 1.0000, 0.7700), vec4f(0.2353, 0.4706, 1.0000, 0.7700), vec4f(0.3137, 0.3922, 0.7843, 0.7700));
const MD_COLORS_SUNSET = array<vec4f, 9>(vec4f(1.0000, 0.3137, 0.1961, 1.0000), vec4f(1.0000, 0.6275, 0.1569, 1.0000), vec4f(1.0000, 0.4706, 0.2353, 1.0000), vec4f(1.0000, 0.7843, 0.1961, 1.0000), vec4f(1.0000, 0.3922, 0.3137, 1.0000), vec4f(1.0000, 0.7059, 0.2353, 1.0000), vec4f(1.0000, 0.2353, 0.2353, 1.0000), vec4f(1.0000, 0.5490, 0.1961, 1.0000), vec4f(1.0000, 0.3529, 0.2745, 1.0000));
const SM_COLORS_SUNSET = array<vec4f, 8>(vec4f(1.0000, 0.7059, 0.1961, 1.0000), vec4f(1.0000, 0.5882, 0.1569, 1.0000), vec4f(1.0000, 0.3137, 0.2353, 1.0000), vec4f(1.0000, 0.3922, 0.3137, 1.0000), vec4f(1.0000, 0.2353, 0.3137, 1.0000), vec4f(1.0000, 0.4706, 0.2353, 1.0000), vec4f(1.0000, 0.7843, 0.1961, 1.0000), vec4f(1.0000, 0.3529, 0.2745, 1.0000));
const MD_INNER_COLORS_SUNSET = array<vec4f, 9>(vec4f(1.0000, 0.3137, 0.1961, 0.4500), vec4f(1.0000, 0.6275, 0.1569, 0.4500), vec4f(1.0000, 0.4706, 0.2353, 0.4500), vec4f(1.0000, 0.7843, 0.1961, 0.4500), vec4f(1.0000, 0.3922, 0.3137, 0.4500), vec4f(1.0000, 0.7059, 0.2353, 0.4500), vec4f(1.0000, 0.2353, 0.2353, 0.4500), vec4f(1.0000, 0.5490, 0.1961, 0.4500), vec4f(1.0000, 0.3529, 0.2745, 0.4500));
const SM_INNER_COLORS_SUNSET = array<vec4f, 8>(vec4f(1.0000, 0.7059, 0.1961, 0.5000), vec4f(1.0000, 0.5882, 0.1569, 0.4500), vec4f(1.0000, 0.3137, 0.2353, 0.3500), vec4f(1.0000, 0.3922, 0.3137, 0.3500), vec4f(1.0000, 0.2353, 0.3137, 0.3000), vec4f(1.0000, 0.4706, 0.2353, 0.4000), vec4f(1.0000, 0.7843, 0.1961, 0.3000), vec4f(1.0000, 0.3529, 0.2745, 0.3000));
const LINE_COLORS_SUNSET_DARK = array<vec4f, 9>(vec4f(1.0000, 0.3922, 0.2353, 1.0000), vec4f(1.0000, 0.7059, 0.1961, 1.0000), vec4f(1.0000, 0.5490, 0.2745, 1.0000), vec4f(1.0000, 0.3137, 0.3137, 1.0000), vec4f(1.0000, 0.7843, 0.2353, 1.0000), vec4f(1.0000, 0.4706, 0.1961, 1.0000), vec4f(1.0000, 0.6275, 0.3137, 1.0000), vec4f(1.0000, 0.3529, 0.2353, 1.0000), vec4f(1.0000, 0.2745, 0.2745, 1.0000));
const LINE_COLORS_SUNSET_LIGHT = array<vec4f, 9>(vec4f(0.8627, 0.3137, 0.1569, 1.0000), vec4f(0.9020, 0.5882, 0.1176, 1.0000), vec4f(0.8235, 0.4314, 0.1961, 1.0000), vec4f(0.7843, 0.2353, 0.2353, 1.0000), vec4f(0.8627, 0.6667, 0.1569, 1.0000), vec4f(0.8235, 0.3922, 0.1176, 1.0000), vec4f(0.9020, 0.5098, 0.2353, 1.0000), vec4f(0.7451, 0.2745, 0.1961, 1.0000), vec4f(0.7059, 0.1961, 0.1961, 1.0000));
const LINE_INNER_COLORS_SUNSET = array<vec4f, 9>(vec4f(1.0000, 0.3922, 0.2353, 0.4800), vec4f(1.0000, 0.7059, 0.1961, 0.4200), vec4f(1.0000, 0.5490, 0.2745, 0.4800), vec4f(1.0000, 0.3137, 0.3137, 0.4200), vec4f(1.0000, 0.7843, 0.2353, 0.5000), vec4f(1.0000, 0.4706, 0.1961, 0.4500), vec4f(1.0000, 0.6275, 0.3137, 0.4000), vec4f(1.0000, 0.3529, 0.2353, 0.4500), vec4f(1.0000, 0.2745, 0.2745, 0.5200));
const LINE_BLOOM_HEAD_SUNSET_DARK = array<vec4f, 2>(vec4f(1.0000, 0.5490, 0.3137, 1.0000), vec4f(1.0000, 0.3922, 0.2353, 0.9800));
const LINE_BLOOM_HEAD_SUNSET_LIGHT = array<vec4f, 2>(vec4f(0.7843, 0.3137, 0.1569, 1.0000), vec4f(0.8627, 0.4706, 0.1176, 1.0000));
const LINE_BLOOM_SPIKES_SUNSET_DARK = array<vec4f, 5>(vec4f(1.0000, 0.3922, 0.3137, 1.0000), vec4f(1.0000, 0.5882, 0.3137, 0.5900), vec4f(1.0000, 0.3137, 0.2353, 1.0000), vec4f(1.0000, 0.4706, 0.1961, 0.9100), vec4f(1.0000, 0.5490, 0.2745, 1.0000));
const LINE_BLOOM_SPIKES_SUNSET_LIGHT = array<vec4f, 5>(vec4f(0.7843, 0.2353, 0.1176, 1.0000), vec4f(0.8627, 0.3922, 0.0784, 0.7000), vec4f(0.7059, 0.1569, 0.0784, 1.0000), vec4f(0.8235, 0.3137, 0.0392, 1.0000), vec4f(0.7451, 0.2745, 0.1176, 1.0000));
const PULSE_RING_COLORS_SUNSET = array<vec4f, 9>(vec4f(1.0000, 0.3137, 0.1961, 1.0000), vec4f(1.0000, 0.6275, 0.1569, 1.0000), vec4f(1.0000, 0.4706, 0.2353, 1.0000), vec4f(1.0000, 0.7843, 0.1961, 1.0000), vec4f(1.0000, 0.3922, 0.3137, 1.0000), vec4f(1.0000, 0.7059, 0.2353, 1.0000), vec4f(1.0000, 0.2353, 0.2353, 1.0000), vec4f(1.0000, 0.5490, 0.1961, 1.0000), vec4f(1.0000, 0.3529, 0.2745, 1.0000));
const PULSE_INNER_COLORS_SUNSET = array<vec4f, 9>(vec4f(1.0000, 0.3137, 0.1961, 0.4500), vec4f(1.0000, 0.6275, 0.1569, 0.4500), vec4f(1.0000, 0.4706, 0.2353, 0.4500), vec4f(1.0000, 0.7843, 0.1961, 0.4500), vec4f(1.0000, 0.3922, 0.3137, 0.4500), vec4f(1.0000, 0.7059, 0.2353, 0.4500), vec4f(1.0000, 0.2353, 0.2353, 0.4500), vec4f(1.0000, 0.5490, 0.1961, 0.4500), vec4f(1.0000, 0.3529, 0.2745, 0.4500));
const PULSE_INNER_BLOOM_COLORS_SUNSET = array<vec4f, 7>(vec4f(1.0000, 0.3137, 0.1961, 0.7750), vec4f(1.0000, 0.6275, 0.1569, 0.7750), vec4f(1.0000, 0.4706, 0.2353, 0.7750), vec4f(1.0000, 0.3922, 0.3137, 0.7750), vec4f(1.0000, 0.7059, 0.2353, 0.7750), vec4f(1.0000, 0.2353, 0.2353, 0.7750), vec4f(1.0000, 0.3529, 0.2745, 0.7750));
const PULSE_OUTER_CORE_COLORS_SUNSET = array<vec4f, 8>(vec4f(1.0000, 0.3137, 0.1961, 1.0000), vec4f(1.0000, 0.2353, 0.2353, 1.0000), vec4f(1.0000, 0.5490, 0.1961, 1.0000), vec4f(1.0000, 0.3529, 0.2745, 1.0000), vec4f(1.0000, 0.3922, 0.3137, 1.0000), vec4f(1.0000, 0.6275, 0.1569, 1.0000), vec4f(1.0000, 0.4706, 0.2353, 1.0000), vec4f(1.0000, 0.7843, 0.1961, 1.0000));
const PULSE_OUTER_BLOOM_COLORS_SUNSET = array<vec4f, 7>(vec4f(1.0000, 0.3137, 0.1961, 0.7700), vec4f(1.0000, 0.2353, 0.2353, 0.7700), vec4f(1.0000, 0.5490, 0.1961, 0.7700), vec4f(1.0000, 0.3529, 0.2745, 0.7700), vec4f(1.0000, 0.3922, 0.3137, 0.7700), vec4f(1.0000, 0.6275, 0.1569, 0.7700), vec4f(1.0000, 0.4706, 0.2353, 0.7700));

struct Accum {
  color: vec3f,
  alpha: f32,
}

fn over(acc: Accum, uv: vec2f, center: vec2f, radius: vec2f, source: vec4f) -> Accum {
  let distance = length((uv - center) / max(radius, vec2f(0.0001)));
  let contribution = clamp(1.0 - distance, 0.0, 1.0) * source.a * (1.0 - acc.alpha);
  return Accum(acc.color + source.rgb * contribution, acc.alpha + contribution);
}

fn beamUv(uv: vec2f) -> vec2f {
  return (uv * vec2f(params.width, params.height) + vec2f(params.offsetX, params.offsetY)) / vec2f(params.rootWidth, params.rootHeight);
}

fn osc(a: f32, b: f32, period: f32, delay: f32) -> f32 {
  let phase = (params.time - delay) / max(period, 0.001);
  return a + (b - a) * ((1.0 - cos(6.2831853 * phase)) * 0.5);
}

fn pulsePoint(acc: Accum, uv: vec2f, shape: vec4f, region: i32, quad: i32, source: vec4f, frozen: bool) -> Accum {
  let inner = params.mode == 3.0;
  let dark = params.dark > 0.5;
  let sp = 0.28;
  let drift = select(select(19.0, 14.0, dark), select(40.0, 33.0, dark), inner);
  let op = select(select(0.0, 0.46, dark), select(0.45, 0.48, dark), inner);
  let gh = select(select(0.58, 0.16, dark), select(0.22, 0.34, dark), inner);
  let scaleDuration = params.duration / 2.3;
  let bs = select(2.3, 1.9, inner) * scaleDuration;
  let ss = select(6.4, 2.6, inner) * scaleDuration;
  let ghs = select(3.8, 2.4, inner) * scaleDuration;
  let bw1 = osc(1.0 - sp, 1.0 + sp * 1.1, ss * 0.9, 0.0);
  let bw2 = osc(1.0 + sp, 1.0 - sp * 0.85, ss * 1.1, 0.0);
  let bw3 = osc(1.0 - sp * 0.6, 1.0 + sp * 1.15, ss * 0.98, 0.0);
  let bh1 = osc(1.0 + sp * 0.9, 1.0 - sp * 0.85, ss * 1.26, 0.0);
  let bh2 = osc(1.0 - sp * 0.8, 1.0 + sp * 1.05, ss * 0.81, 0.0);
  let bh3 = osc(1.0 + sp * 0.75, 1.0 - sp, ss * 1.4, 0.0);
  let bgh = osc(1.0 - gh, 1.0 + gh, ghs, 0.0);
  let bx1 = osc(-drift, drift * 0.9, bs * 1.6, 0.0);
  let by1 = osc(drift * 0.55, -drift * 0.7, bs * 1.6, 0.0);
  let bx2 = osc(drift * 0.8, -drift * 0.9, bs * 1.88, 0.0);
  let by2 = osc(-drift, drift * 0.65, bs * 1.88, 0.0);
  let bx3 = osc(-drift * 0.6, drift, bs * 1.45, 0.0);
  let by3 = osc(-drift * 0.85, drift * 0.45, bs * 1.45, 0.0);
  let widthScale = select(select(bw3, bw2, region == 2), bw1, region == 1);
  let heightScale = select(select(bh3, bh2, region == 2), bh1, region == 1) * bgh;
  let movement = select(select(vec2f(bx3, by3), vec2f(bx2, by2), region == 2), vec2f(bx1, by1), region == 1);
  let quadrantOpacity = select(select(select(1.0 - op, 1.0, quad == 0), osc(1.0 - op, 1.0, bs * 1.32, bs * 0.28), quad == 1), osc(1.0 - op, 1.0, bs * 0.84, bs * 0.55), quad == 2);
  let finalOpacity = select(osc(1.0 - op, 1.0, bs * 1.58, bs * 0.83), quadrantOpacity, quad < 3);
  let finalWidth = select(widthScale, 1.0, frozen);
  let finalHeight = select(heightScale, 1.0, frozen);
  let finalMovement = select(movement, vec2f(0.0), frozen);
  let center = shape.xy + finalMovement / vec2f(params.rootWidth, params.rootHeight);
  let radius = shape.zw * vec2f(finalWidth * params.pulseScaleX, finalHeight * params.pulseScaleY) / vec2f(params.rootWidth, params.rootHeight);
  return over(acc, uv, center, radius, vec4f(source.rgb, source.a * select(finalOpacity, 1.0, frozen)));
}

fn paletteColor(variant: i32, index: i32, kind: i32, dark: bool) -> vec4f {
  if (variant == 0) {
    if (kind == 0) { return MD_COLORS_COLORFUL[index]; }
    if (kind == 1) { return SM_COLORS_COLORFUL[index]; }
    if (kind == 2) { return select(LINE_COLORS_COLORFUL_LIGHT[index], LINE_COLORS_COLORFUL_DARK[index], dark); }
    if (kind == 3) { return LINE_INNER_COLORS_COLORFUL[index]; }
    if (kind == 4) { return MD_INNER_COLORS_COLORFUL[index]; }
    if (kind == 5) { return SM_INNER_COLORS_COLORFUL[index]; }
    if (kind == 6) { return PULSE_RING_COLORS_COLORFUL[index]; }
    return PULSE_INNER_COLORS_COLORFUL[index];
  }
  if (variant == 1) {
    if (kind == 0) { return MD_COLORS_MONO[index]; }
    if (kind == 1) { return SM_COLORS_MONO[index]; }
    if (kind == 2) { return select(LINE_COLORS_MONO_LIGHT[index], LINE_COLORS_MONO_DARK[index], dark); }
    if (kind == 3) { return LINE_INNER_COLORS_MONO[index]; }
    if (kind == 4) { return MD_INNER_COLORS_MONO[index]; }
    if (kind == 5) { return SM_INNER_COLORS_MONO[index]; }
    if (kind == 6) { return PULSE_RING_COLORS_MONO[index]; }
    return PULSE_INNER_COLORS_MONO[index];
  }
  if (variant == 2) {
    if (kind == 0) { return MD_COLORS_OCEAN[index]; }
    if (kind == 1) { return SM_COLORS_OCEAN[index]; }
    if (kind == 2) { return select(LINE_COLORS_OCEAN_LIGHT[index], LINE_COLORS_OCEAN_DARK[index], dark); }
    if (kind == 3) { return LINE_INNER_COLORS_OCEAN[index]; }
    if (kind == 4) { return MD_INNER_COLORS_OCEAN[index]; }
    if (kind == 5) { return SM_INNER_COLORS_OCEAN[index]; }
    if (kind == 6) { return PULSE_RING_COLORS_OCEAN[index]; }
    return PULSE_INNER_COLORS_OCEAN[index];
  }
  if (variant == 3) {
    if (kind == 0) { return MD_COLORS_SUNSET[index]; }
    if (kind == 1) { return SM_COLORS_SUNSET[index]; }
    if (kind == 2) { return select(LINE_COLORS_SUNSET_LIGHT[index], LINE_COLORS_SUNSET_DARK[index], dark); }
    if (kind == 3) { return LINE_INNER_COLORS_SUNSET[index]; }
    if (kind == 4) { return MD_INNER_COLORS_SUNSET[index]; }
    if (kind == 5) { return SM_INNER_COLORS_SUNSET[index]; }
    if (kind == 6) { return PULSE_RING_COLORS_SUNSET[index]; }
    return PULSE_INNER_COLORS_SUNSET[index];
  }
  return vec4f(1.0, 1.0, 1.0, 1.0);
}

fn lineBloomColor(variant: i32, index: i32, dark: bool) -> vec4f {
  if (variant == 0) {
    if (index < 2) { return select(LINE_BLOOM_HEAD_COLORFUL_LIGHT[index], LINE_BLOOM_HEAD_COLORFUL_DARK[index], dark); }
    return select(LINE_BLOOM_SPIKES_COLORFUL_LIGHT[index - 2], LINE_BLOOM_SPIKES_COLORFUL_DARK[index - 2], dark);
  }
  if (variant == 1) {
    if (index < 2) { return select(LINE_BLOOM_HEAD_MONO_LIGHT[index], LINE_BLOOM_HEAD_MONO_DARK[index], dark); }
    return select(LINE_BLOOM_SPIKES_MONO_LIGHT[index - 2], LINE_BLOOM_SPIKES_MONO_DARK[index - 2], dark);
  }
  if (variant == 2) {
    if (index < 2) { return select(LINE_BLOOM_HEAD_OCEAN_LIGHT[index], LINE_BLOOM_HEAD_OCEAN_DARK[index], dark); }
    return select(LINE_BLOOM_SPIKES_OCEAN_LIGHT[index - 2], LINE_BLOOM_SPIKES_OCEAN_DARK[index - 2], dark);
  }
  if (variant == 3) {
    if (index < 2) { return select(LINE_BLOOM_HEAD_SUNSET_LIGHT[index], LINE_BLOOM_HEAD_SUNSET_DARK[index], dark); }
    return select(LINE_BLOOM_SPIKES_SUNSET_LIGHT[index - 2], LINE_BLOOM_SPIKES_SUNSET_DARK[index - 2], dark);
  }
  return vec4f(1.0, 1.0, 1.0, 1.0);
}

fn pulseColor(variant: i32, index: i32, kind: i32) -> vec4f {
  if (variant == 0) {
    if (kind == 0) { return PULSE_INNER_BLOOM_COLORS_COLORFUL[index]; }
    if (kind == 1) { return PULSE_OUTER_CORE_COLORS_COLORFUL[index]; }
    return PULSE_OUTER_BLOOM_COLORS_COLORFUL[index];
  }
  if (variant == 1) {
    if (kind == 0) { return PULSE_INNER_BLOOM_COLORS_MONO[index]; }
    if (kind == 1) { return PULSE_OUTER_CORE_COLORS_MONO[index]; }
    return PULSE_OUTER_BLOOM_COLORS_MONO[index];
  }
  if (variant == 2) {
    if (kind == 0) { return PULSE_INNER_BLOOM_COLORS_OCEAN[index]; }
    if (kind == 1) { return PULSE_OUTER_CORE_COLORS_OCEAN[index]; }
    return PULSE_OUTER_BLOOM_COLORS_OCEAN[index];
  }
  if (variant == 3) {
    if (kind == 0) { return PULSE_INNER_BLOOM_COLORS_SUNSET[index]; }
    if (kind == 1) { return PULSE_OUTER_CORE_COLORS_SUNSET[index]; }
    return PULSE_OUTER_BLOOM_COLORS_SUNSET[index];
  }
  return vec4f(1.0, 1.0, 1.0, 1.0);
}

fn rotateHue(color: vec3f, degrees: f32) -> vec3f {
  let angle = degrees * 0.0174532925;
  let c = cos(angle);
  let s = sin(angle);
  return vec3f(
    (0.213 + 0.787 * c - 0.213 * s) * color.r + (0.715 - 0.715 * c - 0.715 * s) * color.g + (0.072 - 0.072 * c + 0.928 * s) * color.b,
    (0.213 - 0.213 * c + 0.143 * s) * color.r + (0.715 + 0.285 * c + 0.140 * s) * color.g + (0.072 + 0.072 * c - 0.283 * s) * color.b,
    (0.213 - 0.213 * c - 0.787 * s) * color.r + (0.715 - 0.715 * c + 0.715 * s) * color.g + (0.072 + 0.928 * c + 0.072 * s) * color.b
  );
}

fn lineTravel(phase: f32) -> f32 {
  if (phase < 0.1) { return 0.06 + (0.15 - 0.06) * phase / 0.1; }
  if (phase < 0.2) { return 0.15 + (0.25 - 0.15) * (phase - 0.1) / 0.1; }
  if (phase < 0.3) { return 0.25 + (0.35 - 0.25) * (phase - 0.2) / 0.1; }
  if (phase < 0.4) { return 0.35 + (0.44 - 0.35) * (phase - 0.3) / 0.1; }
  if (phase < 0.5) { return 0.44 + (0.50 - 0.44) * (phase - 0.4) / 0.1; }
  if (phase < 0.6) { return 0.50 + (0.56 - 0.50) * (phase - 0.5) / 0.1; }
  if (phase < 0.7) { return 0.56 + (0.65 - 0.56) * (phase - 0.6) / 0.1; }
  if (phase < 0.8) { return 0.65 + (0.75 - 0.65) * (phase - 0.7) / 0.1; }
  if (phase < 0.9) { return 0.75 + (0.85 - 0.75) * (phase - 0.8) / 0.1; }
  return 0.85 + (0.94 - 0.85) * (phase - 0.9) / 0.1;
}

fn lineEdge(phase: f32) -> f32 {
  if (phase < 0.125) { return 0.0; }
  if (phase < 0.325) { return (phase - 0.125) / 0.2; }
  if (phase < 0.675) { return 1.0; }
  if (phase < 0.875) { return 1.0 - (phase - 0.675) / 0.2; }
  return 0.0;
}

fn conicOverlay(uv: vec2f, dark: bool, bloom: bool, spin: f32) -> vec3f {
  let angle = fract(atan2(uv.y - 0.5, uv.x - 0.5) / 6.2831853 + 0.5 - spin);
  let center = select(0.69, 0.70, bloom);
  let distance = abs(fract(angle - center + 0.5) - 0.5);
  let intensity = 1.0 - smoothstep(0.0, select(0.13, 0.12, bloom), distance);
  return mix(vec3f(0.0), select(vec3f(1.0), vec3f(0.0), dark), intensity * select(0.55, 0.85, bloom));
}

fn renderPulse(uv: vec2f, variant: i32, layer: i32) -> Accum {
  var acc = Accum(vec3f(0.0), 0.0);
  if (params.mode == 3.0 && layer == 0) {
  acc = pulsePoint(acc, uv, PULSE_RING_SHAPES[0], 1, 0, pulseColor(variant, 0, 3), false);
  acc = pulsePoint(acc, uv, PULSE_RING_SHAPES[1], 2, 0, pulseColor(variant, 1, 3), false);
  acc = pulsePoint(acc, uv, PULSE_RING_SHAPES[2], 3, 2, pulseColor(variant, 2, 3), false);
  acc = pulsePoint(acc, uv, PULSE_RING_SHAPES[3], 1, 2, pulseColor(variant, 3, 3), false);
  acc = pulsePoint(acc, uv, PULSE_RING_SHAPES[4], 2, 3, pulseColor(variant, 4, 3), false);
  acc = pulsePoint(acc, uv, PULSE_RING_SHAPES[5], 3, 3, pulseColor(variant, 5, 3), false);
  acc = pulsePoint(acc, uv, PULSE_RING_SHAPES[6], 1, 1, pulseColor(variant, 6, 3), false);
  acc = pulsePoint(acc, uv, PULSE_RING_SHAPES[7], 2, 1, pulseColor(variant, 7, 3), false);
  acc = pulsePoint(acc, uv, PULSE_RING_SHAPES[8], 3, 1, pulseColor(variant, 8, 3), false);
  } else if (params.mode == 3.0 && layer == 1) {
  acc = pulsePoint(acc, uv, PULSE_INNER_SHAPES[0], 1, 0, pulseColor(variant, 0, 4), false);
  acc = pulsePoint(acc, uv, PULSE_INNER_SHAPES[1], 2, 0, pulseColor(variant, 1, 4), false);
  acc = pulsePoint(acc, uv, PULSE_INNER_SHAPES[2], 3, 2, pulseColor(variant, 2, 4), false);
  acc = pulsePoint(acc, uv, PULSE_INNER_SHAPES[3], 1, 2, pulseColor(variant, 3, 4), false);
  acc = pulsePoint(acc, uv, PULSE_INNER_SHAPES[4], 2, 3, pulseColor(variant, 4, 4), false);
  acc = pulsePoint(acc, uv, PULSE_INNER_SHAPES[5], 3, 3, pulseColor(variant, 5, 4), false);
  acc = pulsePoint(acc, uv, PULSE_INNER_SHAPES[6], 1, 1, pulseColor(variant, 6, 4), false);
  acc = pulsePoint(acc, uv, PULSE_INNER_SHAPES[7], 2, 1, pulseColor(variant, 7, 4), false);
  acc = pulsePoint(acc, uv, PULSE_INNER_SHAPES[8], 3, 1, pulseColor(variant, 8, 4), false);
  } else if (params.mode == 3.0 && layer == 2) {
  acc = pulsePoint(acc, uv, PULSE_INNER_BLOOM_SHAPES[0], 1, 0, pulseColor(variant, 0, 0), true);
  acc = pulsePoint(acc, uv, PULSE_INNER_BLOOM_SHAPES[1], 2, 0, pulseColor(variant, 1, 0), true);
  acc = pulsePoint(acc, uv, PULSE_INNER_BLOOM_SHAPES[2], 3, 2, pulseColor(variant, 2, 0), true);
  acc = pulsePoint(acc, uv, PULSE_INNER_BLOOM_SHAPES[3], 2, 3, pulseColor(variant, 3, 0), true);
  acc = pulsePoint(acc, uv, PULSE_INNER_BLOOM_SHAPES[4], 3, 3, pulseColor(variant, 4, 0), true);
  acc = pulsePoint(acc, uv, PULSE_INNER_BLOOM_SHAPES[5], 1, 1, pulseColor(variant, 5, 0), true);
  acc = pulsePoint(acc, uv, PULSE_INNER_BLOOM_SHAPES[6], 3, 1, pulseColor(variant, 6, 0), true);
  } else if (params.mode == 4.0 && (layer == 0 || layer == 1)) {
  acc = pulsePoint(acc, uv, PULSE_OUTER_CORE_SHAPES[0], 1, 0, pulseColor(variant, 0, 1), false);
  acc = pulsePoint(acc, uv, PULSE_OUTER_CORE_SHAPES[1], 2, 1, pulseColor(variant, 1, 1), false);
  acc = pulsePoint(acc, uv, PULSE_OUTER_CORE_SHAPES[2], 3, 1, pulseColor(variant, 2, 1), false);
  acc = pulsePoint(acc, uv, PULSE_OUTER_CORE_SHAPES[3], 1, 3, pulseColor(variant, 3, 1), false);
  acc = pulsePoint(acc, uv, PULSE_OUTER_CORE_SHAPES[4], 2, 3, pulseColor(variant, 4, 1), false);
  acc = pulsePoint(acc, uv, PULSE_OUTER_CORE_SHAPES[5], 3, 2, pulseColor(variant, 5, 1), false);
  acc = pulsePoint(acc, uv, PULSE_OUTER_CORE_SHAPES[6], 1, 2, pulseColor(variant, 6, 1), false);
  acc = pulsePoint(acc, uv, PULSE_OUTER_CORE_SHAPES[7], 2, 0, pulseColor(variant, 7, 1), false);
  } else if (params.mode == 4.0 && layer == 2) {
  acc = pulsePoint(acc, uv, PULSE_OUTER_BLOOM_SHAPES[0], 1, 0, pulseColor(variant, 0, 2), true);
  acc = pulsePoint(acc, uv, PULSE_OUTER_BLOOM_SHAPES[1], 2, 1, pulseColor(variant, 1, 2), true);
  acc = pulsePoint(acc, uv, PULSE_OUTER_BLOOM_SHAPES[2], 3, 1, pulseColor(variant, 2, 2), true);
  acc = pulsePoint(acc, uv, PULSE_OUTER_BLOOM_SHAPES[3], 1, 3, pulseColor(variant, 3, 2), true);
  acc = pulsePoint(acc, uv, PULSE_OUTER_BLOOM_SHAPES[4], 2, 3, pulseColor(variant, 4, 2), true);
  acc = pulsePoint(acc, uv, PULSE_OUTER_BLOOM_SHAPES[5], 3, 2, pulseColor(variant, 5, 2), true);
  acc = pulsePoint(acc, uv, PULSE_OUTER_BLOOM_SHAPES[6], 1, 2, pulseColor(variant, 6, 2), true);
  }
  return acc;
}

@fragment
fn fs_main(@location(0) uv: vec2f) -> @location(0) vec4f {
  let localUv = beamUv(uv);
  let mode = i32(params.mode + 0.5);
  let variant = i32(params.variant + 0.5);
  let layer = i32(params.layer + 0.5);
  let dark = params.dark > 0.5;
  var acc = Accum(vec3f(0.0), 0.0);

  if (mode >= 3) {
    acc = renderPulse(localUv, variant, layer);
  } else if (mode == 2) {
    let phase = fract(params.time / max(params.duration, 0.001));
    let x = lineTravel(phase);
    let edge = lineEdge(phase);
    if (layer == 0) {
      for (var i = 0; i < 9; i++) {
        let shape = select(LINE_SHAPES_LIGHT[i], LINE_SHAPES_DARK[i], dark);
        let center = vec2f(x + shape.x / params.rootWidth, 1.0 + shape.y / params.rootHeight);
        acc = over(acc, localUv, center, shape.zw / vec2f(params.rootWidth, params.rootHeight), vec4f(paletteColor(variant, i, 2, dark).rgb, paletteColor(variant, i, 2, dark).a * edge));
      }
      acc.color += conicOverlay(localUv, dark, false, 0.0) * acc.alpha;
    } else if (layer == 1) {
      for (var i = 0; i < 9; i++) {
        let shape = LINE_INNER_SHAPES[i];
        let source = paletteColor(variant, i, 3, dark);
        let center = vec2f(x + shape.x / params.rootWidth, 1.0 + shape.y / params.rootHeight);
        acc = over(acc, localUv, center, shape.zw / vec2f(params.rootWidth, params.rootHeight), vec4f(source.rgb, source.a * edge));
      }
    } else {
      for (var i = 0; i < 7; i++) {
        let shape = LINE_BLOOM_SHAPES[i];
        acc = over(acc, localUv, shape.xy / 100.0, shape.zw / vec2f(params.rootWidth, params.rootHeight), lineBloomColor(variant, i, dark));
      }
    }
  } else if (layer == 0) {
    let count = select(8, 9, mode == 0);
    for (var i = 0; i < 9; i++) {
      if (i >= count) { continue; }
      let shape = select(SM_SHAPES[i], MD_SHAPES[i], mode == 0);
      acc = over(acc, localUv, shape.xy, shape.zw / vec2f(params.rootWidth, params.rootHeight), paletteColor(variant, i, mode, dark));
    }
    let spin = fract(params.time / max(params.duration, 0.001));
    acc.color += conicOverlay(localUv, dark, false, spin) * acc.alpha;
  } else if (layer == 1) {
    let count = select(8, 9, mode == 0);
    for (var i = 0; i < 9; i++) {
      if (i >= count) { continue; }
      let shape = select(SM_SHAPES[i], MD_SHAPES[i], mode == 0);
      var source = paletteColor(variant, i, 4, dark);
      if (mode != 0) { source = paletteColor(variant, i, 5, dark); }
      acc = over(acc, localUv, shape.xy, shape.zw * 0.9 / vec2f(params.rootWidth, params.rootHeight), source);
    }
  } else {
    let spin = fract(params.time / max(params.duration, 0.001));
    acc = Accum(conicOverlay(localUv, dark, true, spin), 1.0);
  }

  let hue = select(-cos(6.2831853 * params.time / 12.0) * params.hueRange, fract(params.time / select(14.0, 16.0, mode == 3)) * 360.0, mode >= 3);
  let rotated = rotateHue(acc.color, select(hue, 0.0, params.staticColors > 0.5));
  let luma = dot(rotated, vec3f(0.213, 0.715, 0.072));
  let adjusted = mix(vec3f(luma), rotated, params.saturation) * params.brightness;
  let finalAlpha = clamp(acc.alpha, 0.0, 1.0);
  return vec4f(adjusted * finalAlpha, finalAlpha);
}
