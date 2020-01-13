package util

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

func PathAdaptive(path string) (string, error) {
	if index := strings.Index(path, "/"); index != 0 {
		return "", errors.New("file routing param illegal param")
	}
	var osType = runtime.GOOS
	rootPath, err := os.Getwd()
	if err != nil {
		return "", errors.New("system os.pwd() error")
	}
	var adaptivePath = rootPath + path
	if osType == "windows" {
		adaptivePath = strings.ReplaceAll(adaptivePath, "/", "\\")
	}
	return adaptivePath, nil
}
func ResourcePathAdaptive(path string) string {
	var osType = runtime.GOOS
	if osType == "windows" {
		return strings.ReplaceAll(path, "/", "\\")
	}
	return path
}
