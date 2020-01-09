package upload

import (
	"gfs/app/util"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"time"
)

var (
	bigFileHashCache = cache.New(24*time.Hour, 60*time.Minute)
	chunkHashCache   = make(map[string]int64, 1000)
)

type fileCache struct {
}

func fileCacheInstance() *fileCache {
	return new(fileCache)
}

func (f *fileCache) putBigFileInfo(bigFileHash string, bigFileInfo *bigFileInfo) {
	bigFileHashCache.Set(bigFileHash, bigFileInfo, cache.DefaultExpiration)
}

func (f *fileCache) getParentFileInfo(token string) *bigFileInfo {
	ex, found := bigFileHashCache.Get(token)
	if found {
		return ex.(*bigFileInfo)
	}
	return nil
}

func (f *fileCache) putChunkHash(hash string, timestamp int64) {
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
