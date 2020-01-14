package download

import (
	"fmt"
	"gfs/app/component"
	"gfs/app/logger"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
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
		resource := fileInfo.ResourcePath + "/" + fileInfo.ResourceName
		resourcePath := util.ResourcePathAdaptive(resource)
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.FileName))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(resourcePath)
	} else {
		//token
		token := c.Query("token")
		//TODO check token
		component.GetAccessToken(token)
	}
}
