package main

import (
	"fmt"
	"gfs/app/util"
)

func main() {
	randomStr := util.RandomString(10)
	fmt.Println(randomStr)
	bytes := []byte(randomStr)
	for k, v := range bytes {
		fmt.Printf("k: %d,v: %d \n", k, v)
	}
}
