package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	resource := "D:\\Download\\Google Download\\ps.zip"
	hash := md5.New()
	if file, err := os.OpenFile(resource, os.O_RDWR, 0600); err != nil {
		log.Println(err.Error())
	} else {
		start := time.Now().UnixNano()
		//log.Println(start)g
		var bufForHash = make([]byte, 64*1024)
		if _, err := io.CopyBuffer(hash, file, bufForHash); err != nil {
			log.Println(err.Error())
		}

		log.Println(hex.EncodeToString(hash.Sum(nil)))
		end := time.Now().UnixNano()
		//log.Println(end)
		//log.Println(end-start)
		log.Printf("time:%dms\n", (end-start)/1000000)
	}

}
