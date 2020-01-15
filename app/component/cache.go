package component

import (
	"gfs/app/repository/model"
	"github.com/patrickmn/go-cache"
	"time"
)

//缓存登录认证信息
var authCache = cache.New(30*time.Minute, 10*time.Minute)

//缓存资源认证信息
var accessCache = cache.New(30*time.Minute, 10*time.Minute)

func PutAuthToken(token string, userVo *model.UserVo) {
	authCache.Set(token, userVo, cache.DefaultExpiration)
}

func GetAuthToken(token string) *model.UserVo {
	ex, found := authCache.Get(token)
	if found {
		authCache.Set(token, ex.(*model.UserVo), cache.DefaultExpiration)
		return ex.(*model.UserVo)
	}
	return nil
}

func DeleteAuthToken(token string) {
	authCache.Delete(token)
}

func PutAccessToken(token string, expiration time.Duration) {
	accessCache.Set(token, expiration, cache.DefaultExpiration)
}

func GetAccessToken(token string) *time.Duration {
	ex, found := authCache.Get(token)
	if found {
		r := ex.(time.Duration)
		return &r
	}
	return nil
}
