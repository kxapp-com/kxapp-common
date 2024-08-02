//go:build !windows
// +build !windows

package shellz

import (
	"errors"
)

func CreateLinkWindows(src string, dest string) error {
	return errors.New("windows only support")
}
func CreateLinkJunkWindows(src string, dest string) error {
	return errors.New("windows only support")
}

func SearchRegistryKeys(keyPath string) (map[string]string, error) {
	return nil, errors.New("windows only support")

}
