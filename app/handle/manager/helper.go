package manager

import (
	"crypto/sha256"
	"encoding/hex"
	"gfs/app/common"
	"gfs/app/component"
	"gfs/app/repository/model"
	"github.com/gin-gonic/gin"
	"net/http"
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

func getUserVoByToken(c *gin.Context) *model.UserVo {
	if authToken := c.Request.Header.Get("AuthToken"); authToken == "" {
		c.JSON(http.StatusOK, common.Response{Code: 100, Message: IllegalRequest})
		return nil
	} else {
		if userVo := component.GetAuthToken(authToken); userVo == nil {
			c.JSON(http.StatusOK, common.Response{Code: 100, Message: tokenError})
			return nil
		} else {
			return userVo
		}
	}
}
