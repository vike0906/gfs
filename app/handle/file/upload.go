package file

import (
	"fmt"
	"gfs/app/common"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func UploadHandle(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	path, _ := util.PathAdaptive("/resource/")
	path = path + file.Filename
	// 上传文件至指定目录
	err := c.SaveUploadedFile(file, path)
	if err != nil {
		log.Println(err.Error())
	}
	var response = common.Response{}
	c.JSON(http.StatusOK, response.SuccessWithContent(fmt.Sprintf("'%s' uploaded!", file.Filename)))
}
