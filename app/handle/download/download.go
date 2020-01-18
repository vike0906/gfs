package download

import (
	"fmt"
	"gfs/app/component"
	"gfs/app/logger"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Download(c *gin.Context) {
	//key
	key := c.Param("key")
	fileInfo, err := queryFileByKey(key)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(resourceNotFount))
		return
	}
	//is public
	if fileInfo.Type == 1 {
		//download
		resourceTran(fileInfo.FileName, fileInfo.ResourcePath, fileInfo.ResourceName, c)
	} else {
		//token
		token := c.Query("token")
		ex := c.Query("e")
		if re, err := tokenCheck(token, key, ex); re != true {
			if err != nil {
				c.JSON(http.StatusOK, response.Fail(err.Error()))
			} else {
				c.JSON(http.StatusOK, response.Fail(tokenError))
			}
			return
		}
		if exInt64, err := strconv.ParseInt(ex, 10, 64); err != nil {
			c.JSON(http.StatusOK, response.Fail("expired timestamp is error"))
			return
		} else if exInt64 < time.Now().Unix() {
			c.JSON(http.StatusOK, response.Fail(exError))
			return
		}
		resourceTran(fileInfo.FileName, fileInfo.ResourcePath, fileInfo.ResourceName, c)
	}
}

func resourceTran(fileName, resourcePath, ResourceName string, c *gin.Context) {
	resource := resourcePath + "/" + ResourceName
	resource = util.ResourcePathAdaptive(resource)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	if _, err := os.Stat(resource); err != nil {
		c.JSON(http.StatusOK, response.Fail(resourceIsBreak))
		return
	}
	c.File(resource)
}

func tokenCheck(token, key, ex string) (bool, error) {

	component.GetAuthToken(token)

	tokens := strings.Split(token, ":")
	if len(tokens) != 2 {
		return false, nil
	}
	user, err := queryUserByAppKey(tokens[0])
	if err != nil {
		return false, err
	}
	//TODO serverUrl from database
	var serverUrl = "http://localhost/9090/" + key + "?e=" + ex
	sign := util.CreateSign([]byte(serverUrl), []byte(user.AppSecret))
	if sign == token {
		return true, nil
	}
	return false, nil
}
