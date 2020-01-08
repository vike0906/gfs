package upload

import (
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func BigFileUploadChunk(c *gin.Context) {
	//TODO 检查参数

	//TODO 检查权限
	//已上传
	bigFileHash := c.PostForm("parentFileHash")
	fileBinary, _ := c.FormFile("fileBinary")

	//计算hash值并校验
	hash, err := ulr.hash(fileBinary)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//获取保存目录
	savePath, err := util.PathAdaptive("/resource/temp")
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//保存文件
	size, err := ulr.saveFile(fileBinary, savePath, hash)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	log.Println(size)
	//TODO 添加/修改记录
	bigFileHashMap := getBigFileHash(bigFileHash)
	if bigFileHashMap == nil {
		c.JSON(http.StatusOK, response.Fail("parent file hash is error"))
		return
	}
	var b = *bigFileHashMap
	b[hash] = ""
	putChunkHash(hash, time.Now().Unix())

	c.JSON(http.StatusOK, response.SuccessWithMessage(""))
}

//TODO file routing param illegal param
//TODO fileBinary, _ := c.FormFile("fileBinary")

func BigFileUploadMerge(c *gin.Context) {

}
