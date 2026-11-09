# sfsymbols

Apple's [SF Symbols](https://developer.apple.com/sf-symbols/) for
[Golet](https://github.com/HyperparameterSI/golet/backend) apps of native UI on macOS: the
symbols the running system has, drawn by AppKit and handed over as an
image.

```go
import "github.com/HyperparameterSI/golet/backend/plugins/sfsymbols"

sfsymbols.Symbol(c, "gearshape")
sfsymbols.Symbol(c, "trash").Size(20).Color(t.Accent).Weight(sfsymbols.Semibold).Fill()
```

See [the documentation](../../docs/plugins/sfsymbols.md), and
`examples/sf-symbols`. Nothing of
Apple's artwork is in the plugin: the symbols come from the system.
