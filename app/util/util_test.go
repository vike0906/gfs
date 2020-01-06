package util

import (
	"fmt"
	"testing"
)

func TestPathAdaptive(t *testing.T) {
	re, _ := PathAdaptive("/dist/")
	fmt.Println(re)
}

func TestUUID(t *testing.T) {
	fmt.Println(UUID())
}

func TestRandomString(t *testing.T) {
	fmt.Println(RandomString(10))
}
