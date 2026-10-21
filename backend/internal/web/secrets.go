// Author: HyperparameterSI

//go:build js

package web

import "github.com/HyperparameterSI/golet/backend/internal/platform"

func (*Backend) Secrets() platform.Secrets { return secrets{} }

type secrets struct{}

func (secrets) SetSecret(string, string, []byte) error { return platform.ErrUnsupported }
func (secrets) Secret(string, string) ([]byte, error)  { return nil, platform.ErrUnsupported }
func (secrets) DeleteSecret(string, string) error      { return platform.ErrUnsupported }
