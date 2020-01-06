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
	var adaptivePath string
	if osType == "windows" {
		path = strings.ReplaceAll(path, "/", "\\")
		adaptivePath = rootPath + path
	} else if osType == "linux" {
		adaptivePath = rootPath + path
	}
	return adaptivePath, nil
}
