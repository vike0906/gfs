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

//校验权限

//获取储存根目录

//获取用户目录

//获取日期

//文件信息写入数据库

//保存文件

//发起备份任务
