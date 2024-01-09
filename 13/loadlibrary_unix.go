//go:build linux || darwin
// +build linux darwin

package gccjit

import "github.com/ebitengine/purego"

func loadLibrary(path string) (uintptr, error) {
	return purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}
