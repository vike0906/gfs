package upload

import (
	"fmt"
	"gfs/app/common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	errMessage      = "need param %s, but not found"
	fileDownloadUrl = "http://host:port/download/"
)

var (
	ulr      uploader = new(uploadHelper)
	response          = common.Response{}
)

func SmallFileUpload(c *gin.Context) {

	if err := smallFUParamCheck(c); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}

	//TODO 权限校验

	//计算hash值并校验
	fileData, err := c.FormFile("fileData")
	hash, err := ulr.hash(fileData)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//TODO 校验文件准确性

	//获取保存目录
	savePath, err := ulr.gainSavePath("userName")
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}

	//保存文件
	size, err := ulr.saveFile(fileData, savePath, hash)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}

	fileName := c.PostForm("fileName")

	//TODO 文件信息写入数据库

	c.JSON(http.StatusOK, response.SuccessWithContent(NewSmallFUResult(fileName, hash, fileDownloadUrl+hash, size)))

	//TODO 发起备份任务
	log.Println("抵达此处")
}

func BigFileUploadInit(c *gin.Context) {

}

func smallFUParamCheck(c *gin.Context) error {
	if err := paramCheck(c); err != nil {
		return err
	}
	if fileName := c.PostForm("fileName"); fileName == "" {
		return &common.GfsError{Message: fmt.Sprintf(errMessage, "fileName")}
	}
	if fileData, err := c.FormFile("fileData"); err != nil || fileData.Size == 0 {
		return &common.GfsError{Message: fmt.Sprintf(errMessage, "fileData")}
	}
	return nil
}

func paramCheck(c *gin.Context) error {

	if uploadToken := c.PostForm("uploadToken"); uploadToken == "" {
		return &common.GfsError{Message: fmt.Sprintf(errMessage, "uploadToken")}
	}
	if fileHash := c.PostForm("fileHash"); fileHash == "" {
		return &common.GfsError{Message: fmt.Sprintf(errMessage, "fileHash")}
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
