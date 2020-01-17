package main

import (
	"fmt"
)

func main() {
	testMap := make(map[string]int16, 100)
	testMap["x"] = 666
	fmt.Println(testMap)
	fmt.Println(len(testMap))
	//if _,err := os.OpenFile("./osTest.txt", os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0600);err!=nil{
	//	fmt.Println("create success")
	//}
}
