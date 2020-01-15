package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gfs/app/util"
)

func main() {
	salt := util.RandomString(10)
	fmt.Printf("salt值:%s\n", salt)
	fmt.Printf("hash值:%s\n", passwordHash("123456", salt))
	uuid := util.UUID()
	fmt.Printf("uuid:%s\n", uuid)
	fmt.Printf("appSecretHash:%s\n", appSecretHash(uuid, "bxv71xnivu"))

}

func passwordHash(psd, salt string) string {
	hash := sha256.New()
	hash.Reset()
	hash.Write([]byte(salt))
	for i := 0; i < 10; i++ {
		hash.Write([]byte(psd))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func appSecretHash(appKey, salt string) string {
	hash := sha1.New()
	hash.Reset()
	hash.Write([]byte(salt))
	for i := 0; i < 10; i++ {
		hash.Write([]byte(appKey))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
