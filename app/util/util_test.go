package util

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPathAdaptive(t *testing.T) {
	re, _ := PathAdaptive("/dist/")
	fmt.Println(re)
}

func TestUUID(t *testing.T) {
	var sysMap sync.Map
	var wg sync.WaitGroup
	wg.Add(1000000)
	var begin int64 = time.Now().UnixNano()
	for i := 0; i < 1000000; i++ {
		go func() {
			sysMap.Store(UUID(), 0)
			wg.Done()
		}()
	}
	wg.Wait()
	var process int64 = time.Now().UnixNano()
	fmt.Printf("create end, time consuming：%d ms\n", (process-begin)/1000000)
	var count uint64 = 0
	sysMap.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	fmt.Println(count)
	var end int64 = time.Now().UnixNano()
	fmt.Printf("count end, time consuming：%d ms\n", (end-process)/1000000)
	fmt.Printf("total time consuming：%d ms\n", (end-begin)/1000000)

}

func TestRandomString(t *testing.T) {
	fmt.Println(RandomString(10))
}
