package upload

import (
	"gfs/app/util"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"time"
)

var bigFileHashCache = cache.New(24*time.Hour, 60*time.Minute)

var chunkHashCache = make(map[string]int64, 1000)

func putBigFileHash(bigFileHash string, bigFileChunk map[string]string) {
	bigFileHashCache.Set(bigFileHash, bigFileChunk, cache.DefaultExpiration)
}

func getBigFileHash(token string) *map[string]string {
	ex, found := bigFileHashCache.Get(token)
	if found {
		r := ex.(map[string]string)
		return &r
	}
	return nil
}

func putChunkHash(hash string, timestamp int64) {
	chunkHashCache[hash] = timestamp
}

func init() {
	chunkHashCacheClearTicker := time.NewTicker(time.Hour)
	go func(t *time.Ticker, c *map[string]int64) {
		var expireTime int64 = 24 * 60 * 60
		for {
			<-t.C
			for chunkHash, timestamp := range *c {
				if interval := time.Now().Unix() - timestamp; interval >= expireTime {
					path, err := util.PathAdaptive("/resource/temp/")
					if err != nil {
						log.Println("File routing failed")
					} else {
						var resource = path + chunkHash
						if _, err := os.Stat(resource); err != nil {
							if os.IsExist(err) {
								if err := os.Remove(resource); err != nil {
									log.Println("Temporary cache file deletion failed")
								} else {
									delete(*c, chunkHash)
								}
							}
						} else {
							log.Println("Resource addressing failed")
						}
					}

				}
			}
		}
	}(chunkHashCacheClearTicker, &chunkHashCache)
}
