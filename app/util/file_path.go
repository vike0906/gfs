package util

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

func PathAdaptive(path string) (string, error) {
	if index := strings.Index(path, "/"); index != 0 {
		return "", errors.New("illegal param")
	}
	var osType = runtime.GOOS
	rootPath, _ := os.Getwd()
	var adaptivePath = rootPath + path
	if osType == "windows" {
		adaptivePath = strings.ReplaceAll(adaptivePath, "/", "\\")
	}
	return adaptivePath, nil
}
