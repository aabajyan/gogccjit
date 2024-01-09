//go:build windows
// +build windows

package gccjit

import "golang.org/x/sys/windows"

func loadLibrary(path string) (uintptr, error) {
	ptr, err := windows.LoadLibrary(path)

	return uintptr(ptr), err
}
