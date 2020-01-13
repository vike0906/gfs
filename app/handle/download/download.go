package download

import (
	"fmt"
	"gfs/app/component"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Download(c *gin.Context) {
	//key
	key := c.Param("key")
	fileInfo, err := queryFileByKey(key)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(resourceNotFount))
		return
	}
	log.Println(fileInfo)
	//is public
	if fileInfo.Type == 1 {
		//download
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.FileName))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		resource := fileInfo.ResourcePath + "/" + fileInfo.ResourceName
		log.Println(resource)
		if path, err := util.PathAdaptive(resource); err != nil {
			c.JSON(http.StatusOK, response.Fail(err.Error()))
			return
		} else {
			c.File(path)
		}
	} else {
		//token
		token := c.Query("token")
		//check token
		component.GetAccessToken(token)
	}
}
