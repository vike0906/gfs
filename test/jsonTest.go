package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func main() {
	paramMap := make(map[string]string, 2)
	paramMap["deadline"] = strconv.FormatInt(time.Now().Add(30*time.Minute).Unix(), 10)
	paramMap["permissionType"] = "public"
	s, _ := json.Marshal(paramMap)
	fmt.Println(string(s))
	fmt.Println(base64.RawURLEncoding.EncodeToString(s))
}
