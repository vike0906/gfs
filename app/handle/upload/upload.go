package upload

import (
	"gfs/app/common"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	ulr          uploader = new(uploadHelper)
	response              = common.ResponseInstance()
	tempInstance          = fileCacheInstance()
)

func SmallFileUpload(c *gin.Context) {
	//check params
	if err := smallFUParamCheck(c); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		c.Abort()
	}

	//TODO check uploadToken

	//hash and validation
	fileBinary, err := c.FormFile(paramFileBinary)
	hash, err := ulr.hash(fileBinary)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	if fileHash := c.PostForm(paramFileHash); fileHash != hash {
		c.JSON(http.StatusOK, response.Fail(fileCorrupted))
		return
	}

	//gain save path TODO: from database
	savePath, err := ulr.gainSavePath("persist")
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}

	//save file
	var resourceName = util.UUID()
	size, err := ulr.saveFile(fileBinary, savePath, resourceName)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}

	fileName := c.PostForm(paramFileName)

	//TODO insert file info to database

	c.JSON(http.StatusOK, response.SuccessWithContent(newSmallFUResult(fileName, hash, fileDownloadUrl+resourceName, size)))

	//TODO back up
}

func BigFileUploadInit(c *gin.Context) {
	if err := bigFFUInitCheck(c); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}

	//TODO uploadToken validation

	//TODO judge file is exist?
	hash := c.PostForm(paramFileHash)
	fileName := c.PostForm(paramFileName)

	var isExist bool

	if isExist {
		//TODO crate logic link
		var resourceName = util.UUID()
		//TODO insert file info to database
		c.JSON(http.StatusOK, response.SuccessWithContent(newBigFUIResultIsExist(exist, fileName, hash, fileDownloadUrl+resourceName)))
	} else {
		//not exist
		//return the chunk's hash which is exist
		var existedChunkArray []string
		parentFileInfo := tempInstance.getParentFileInfo(hash)
		if parentFileInfo == nil {
			bigFileInfo, err := newBigFileInfo(fileName, hash, c.PostForm(paramChunkCount))
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusOK, response.Fail(err.Error()))
				return
			}
			tempInstance.putBigFileInfo(hash, bigFileInfo)
			existedChunkArray = make([]string, 0)
		} else {
			chunkInfoMap := *parentFileInfo.chunkInfoMap
			existedChunkArray = make([]string, 0, len(chunkInfoMap))
			for k := range chunkInfoMap {
				existedChunkArray = append(existedChunkArray, k)
			}
		}
		c.JSON(http.StatusOK, response.SuccessWithContent(newBigFUIResultUnExist(unExit, &existedChunkArray)))
	}

}

func BigFileUploadChunk(c *gin.Context) {
	//check params
	if err := bigFFUChunkCheck(c); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}

	//TODO check uploadToken

	//hash and validation
	chunkBinary, _ := c.FormFile(paramChunkBinary)
	chunkHash, err := ulr.hash(chunkBinary)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	if fileHash := c.PostForm(paramChunkHash); fileHash != chunkHash {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(fileCorrupted))
		return
	}
	//gain save path
	savePath, err := util.PathAdaptive("/resource/temp")
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//save chunk file
	size, err := ulr.saveFile(chunkBinary, savePath, chunkHash)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//add or update chunk upload info
	chunkInfo, err := newChunkInfo(chunkHash, c.PostForm(paramChunkIndex), c.PostForm(paramChunkStart), c.PostForm(paramChunkEnd))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	fileHash := c.PostForm(paramFileHash)
	fileInfo := tempInstance.getParentFileInfo(fileHash)
	var m = *fileInfo.chunkInfoMap
	m[chunkHash] = chunkInfo
	fileInfo.chunkInfoMap = &m

	tempInstance.putChunkHash(chunkHash, time.Now().Unix())

	c.JSON(http.StatusOK, response.SuccessWithContent(newBigFUCResult(fileHash, chunkHash, size)))
}

func BigFileUploadMerge(c *gin.Context) {
	//param check
	if err := bigFFUMergeCheck(c); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//TODO validation upload token
	//获取file hash对应的待合并信息
	fileHash := c.PostForm(paramFileHash)
	//检查块是否全部上传
	fileInfo := tempInstance.getParentFileInfo(fileHash)
	chunkInfoMap := *fileInfo.chunkInfoMap
	length := len(chunkInfoMap)
	if fileInfo.chunkCount != uint16(length) {
		c.JSON(http.StatusOK, response.Fail(chunkCountError))
	}
	//将map按照index转为有序数组
	chunkSortArray := make([]*ChunkInfo, length, length)
	for _, v := range chunkInfoMap {
		chunkSortArray[v.Index] = v
	}
	//按顺序写入磁盘（name=uuid,检查start and end）
	//计算文件hash并核对
	//响应 newSmallFUResult(fileName, hash, fileDownloadUrl+resourceName, size)
}

//校验权限

//获取储存根目录

//获取用户目录

//获取日期

//文件信息写入数据库

//保存文件

//发起备份任务
