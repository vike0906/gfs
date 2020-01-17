package download

import (
	"fmt"
	"gfs/app/component"
	"gfs/app/logger"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
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
		//TODO check token
		//TODO manager web server's token
		component.GetAccessToken(token)
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

func tokenCheck(token string) {
	//component.GetAccessToken(token)
	component.GetAuthToken(token)
}
