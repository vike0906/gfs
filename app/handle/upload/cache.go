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

func (f *fileCache) getParentFileInfo(hash string) *bigFileInfo {
	ex, found := bigFileHashCache.Get(hash)
	if found {
		return ex.(*bigFileInfo)
	}
	return nil
}
func (f *fileCache) deleteParentFileInfo(hash string) {
	bigFileHashCache.Delete(hash)
}

func (f *fileCache) putChunkHash(hash string, timestamp int64) {
	chunkHashCache[hash] = timestamp
}

func init() {
	chunkHashCacheClearTicker := time.NewTicker(1 * time.Hour)
	go func(t *time.Ticker, c *map[string]int64) {
		var expireTime int64 = 24 * 60 * 60
		for {
			<-t.C
			log.Println("clear temp cache file task start")
			for chunkHash, timestamp := range *c {
				if interval := time.Now().Unix() - timestamp; interval >= expireTime {
					path, err := util.PathAdaptive("/resource/temp/")
					if err != nil {
						log.Println("File routing failed")
					} else {
						var resource = path + chunkHash

						if _, err := os.Stat(resource); err != nil {
							if os.IsNotExist(err) {
								log.Println("Temporary cache file is not exist")
							} else {
								log.Println("Resource addressing failed")
							}
						} else {
							if err := os.Remove(resource); err != nil {
								log.Println("Temporary cache file deletion failed")
							} else {
								delete(*c, chunkHash)
							}
						}
					}

				}
			}
		}
	}(chunkHashCacheClearTicker, &chunkHashCache)
}
