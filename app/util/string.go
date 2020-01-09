package util

import (
	"fmt"
	"github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

var stringArray = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "v", "u", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func UUID() string {
	uuidStr := uuid.NewV4().String()
	return strings.ReplaceAll(uuidStr, "-", "")
}

func RandomString(length int) string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	var result strings.Builder
	for i := 0; i < length; i++ {
		fmt.Fprint(&result, stringArray[r.Intn(len(stringArray))])
	}
	return result.String()
}
