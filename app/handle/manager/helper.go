package manager

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"gfs/app/component"
	"gfs/app/repository/model"
	"github.com/gin-gonic/gin"
)

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

func getUserVoByToken(c *gin.Context) *model.UserVo {
	authToken := c.Request.Header.Get("AuthToken")
	userVo := component.GetAuthToken(authToken)
	return userVo
}
