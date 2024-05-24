//go:build !windows
// +build !windows

package shellz

import (
	"errors"
	"golang.org/x/sys/windows/registry"
)

func CreateLinkWindows(src string, dest string) error {
	return errors.New("windows only support")
}
func CreateLinkJunkWindows(src string, dest string) error {
	return errors.New("windows only support")
}

func SearchRegistryKeys(baseKey registry.Key, keyPath string) (map[string]string, error) {
	return nil, errors.New("windows only support")

}
