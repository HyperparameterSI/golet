// Author: HyperparameterSI

package unsupported

import "github.com/HyperparameterSI/golet/backend/internal/platform"

func (*Backend) NewFileWatch(func(platform.FileWatchNotice)) (platform.FileWatch, error) {
	return nil, platform.ErrUnsupported
}

func (*Backend) ReplaceFile(target, replacement string) error { return platform.ErrUnsupported }
