package upload

import (
	"fmt"
	"gfs/app/common"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var ulr uploader = new(uploadHelper)
var response = common.Response{}

func FileUploadHandle(c *gin.Context) {
	if err := paramCheck(c); err != nil {
		c.JSON(http.StatusOK, response.Fail(err.Error()))
	}
	file, _ := c.FormFile("data")
	log.Println(c.PostForm("fileName"))
	log.Println(c.PostForm("index"))
	log.Println(file.Filename)
	//TODO 权限校验

	path, _ := util.PathAdaptive("/resource/")

	var uploader uploader = new(uploadHelper)
	p, er := uploader.gainSavePath("userName")
	if er != nil {
		log.Println(er.Error())
		c.JSON(http.StatusOK, response.Fail(er.Error()))
		return
	}
	fmt.Println(*p)
	path = path + file.Filename
	// 上传文件至指定目录
	err := c.SaveUploadedFile(file, path)
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(http.StatusOK, response.SuccessWithContent(fmt.Sprintf("'%s' uploaded!", file.Filename)))
}

func SmallFileUpload(c *gin.Context) {
	if err := paramCheck(c); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
	}
	//TODO 权限校验

	var (
		savePath *string
		err      error
	)
	if savePath, err = ulr.gainSavePath("userName"); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
	}
	fmt.Println(savePath)

}

func BigFileUploadInit(c *gin.Context) {

}

func paramCheck(c *gin.Context) error {

	var errMessage = "need param %s, but not found"

	if uploadToken := c.PostForm("uploadToken"); uploadToken == "" {
		return &common.GfsError{Message: fmt.Sprintf(errMessage, "uploadToken")}
	}
	if fileSize := c.PostForm("fileName"); fileSize == "" {
		return &common.GfsError{Message: fmt.Sprintf(errMessage, "fileName")}
	}
	return nil
}

func uploadTokenCheck(token *string) error {
	return nil
}

//校验权限

//获取储存根目录

//获取用户目录

//获取日期

//文件信息写入数据库

//保存文件

//发起备份任务
