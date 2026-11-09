# permissionflow

The permissions the system keeps for the user to give, such as
Accessibility, Screen Recording and the microphone, for
[Golet](https://github.com/HyperparameterSI/golet/backend) apps: checked, asked for with the
system's prompt, and a guide to System Settings, where a panel under its
window holds the app to drag into the list, as
[PermissionFlow](https://github.com/jaywcjlove/PermissionFlow) does for
Swift apps. macOS and Windows; it is all Go.

```go
import "github.com/HyperparameterSI/golet/backend/plugins/permissionflow"

st, err := permissionflow.Check(permissionflow.Accessibility)
if err == nil && st != permissionflow.Granted {
	st, err = permissionflow.Guide(ctx, permissionflow.Accessibility, permissionflow.GuideOptions{})
}
```

See [the documentation](../../docs/plugins/permissionflow.md), and
`examples/permissionflow`. It follows PermissionFlow (MIT, © 小弟调调;
see [NOTICE.md](NOTICE.md)).
