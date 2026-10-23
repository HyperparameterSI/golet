# glass

Liquid Glass for [Golet](https://github.com/HyperparameterSI/golet) apps of native UI,
as macOS 26 and later draw it, on every platform: a material that what is
under an element shows through, frosted and bent along its edges.

```go
import "github.com/HyperparameterSI/golet/backend/plugins/glass"

ui.Row(c).Padding(8, 16).Radius(22).Material(glass.Glass{}).Children(func() {
	ui.Text(c, "On glass")
})
```

`glass.ScrollEdge` is macOS's scroll edge effect, which content scrolling
under a bar fades (soft) or frosts (hard) under, and `glass.Blur` blurs
what is under an element without the glass, evenly or, with a mask,
progressively:

```go
ui.Box(c).Absolute().Top(0).Left(0).Right(0).Height(74).PassThrough().Material(glass.ScrollEdge{})
ui.Box(c).Absolute().Top(0).Left(0).Right(0).Height(88).PassThrough().
	Material(glass.Blur{Radius: 6, Mask: &ui.LinearGradient{From: ui.RGB(0, 0, 0), To: ui.Transparent, Angle: 180}})
```

See [the documentation](../../docs/plugins/glass.md), and the Glass page of
`examples/gallery`.

The glass and the blur, which a hard scroll edge draws with, are effects
of Golet's renderers (`scene.Effect`): their shaders in Metal Shading
Language, HLSL and GLSL (`glass.metal`, `glass.hlsl`, `glass.glsl`, and
`blur.*`), and their twins for the CPU renderer (`pixels.go`, `blur.go`),
which must draw the same pixels; the tests compare each GPU's drawing
with the CPU's.
`go generate` compiles the shaders ahead of time where it can: the Metal
library on macOS, with Xcode (`shaders_darwin.go`), and the Direct3D
bytecode on Windows (`shaders_windows.go`). Renderers compile the source
when the code compiled ahead of time is older than it or than their own
part of the shader.

## Native control theme

`glass.Theme(dark)` opts styled native controls into Liquid Glass. Cache a light
and dark theme, then call `c.SetTheme(theme)` before building your view:

```go
// Author: HyperparameterSI
var lightGlass = glass.Theme(false)
var darkGlass = glass.Theme(true)

func view(c *ui.Context) {
    theme := lightGlass
    if c.Theme().Dark {
        theme = darkGlass
    }
    c.SetTheme(theme)
    ui.Column(c).Padding(24).Gap(12).Children(func() {
        ui.Text(c, "Workspace")
        ui.Button(c, "Open")
        ui.PrimaryButton(c, "Continue")
    })
}
```

Buttons, text fields, selects, autocomplete, time inputs, sidebar backgrounds,
tab/segmented tracks, popovers and dialogs adopt the material. Toggle selections,
checkboxes, radio buttons, switch tracks, steppers and color-well frames also use
the theme materials. Check marks, selection dots and color swatches retain their
semantic colors so their state stays clear. Cards and document content stay
solid unless explicitly given a material.
For a custom floating panel, use `.Material(theme.PanelMaterial)`.

The theme exposes `ControlMaterial`, `AccentMaterial` and `PanelMaterial`.
Set any slot to nil to retain solid styling for that group. High contrast
automatically skips theme materials; reduced motion uses the UI animation
system. Explicitly applied materials remain under the application's control.
To offer a “Reduce transparency” option, switch to `ui.ShadcnTheme(dark)`.

This uses Golet's existing Metal, Direct3D, OpenGL and CPU effects, not Apple's
AppKit controls or a webview. Refraction samples previously drawn content
**inside the window**, not the desktop behind the window. Place floating controls
over an image, gradient or scrolling content to see the glass; on a uniform
background the effect is intentionally subtle. Avoid covering the whole view
with nested glass panels: each material needs backdrop processing.

The Gallery → Glass page (also opened from Studio) includes a “Liquid Glass
control theme” switch. Studio's native demos share this selection and the
light/dark appearance setting: open windows update immediately, and new windows
inherit it for the current session. Gallery cards also adopt the panel material.
WebView page content and the separate Chromium browser keep their own styling.

### Emphasized glass

The control theme now uses `glass.Vivid`: more visible backdrop detail, stronger
edge refraction and a brighter hairline rim, with separate light/dark tone curves.
It reuses the existing GPU/CPU pipeline without extra blur passes or idle
animation. This is a Golet preset inspired by macOS, not an AppKit API name.
`Glass{}` remains Regular for existing applications.

For a quieter control theme, replace its material slots after creating it:

```go
theme.ControlMaterial = glass.Glass{Interactive: true}
theme.AccentMaterial = glass.Glass{Tint: theme.Accent, Interactive: true}
theme.PanelMaterial = glass.Glass{}
```

Studio's theme lab places glass above a static wallpaper so its edge refraction
can be seen. Select the plain background option to compare it against a flat
surface. High contrast continues to use solid control and panel fills.
