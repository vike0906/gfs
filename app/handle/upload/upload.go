package upload

import (
	"gfs/app/common"
	"gfs/app/logger"
	"gfs/app/repository/model"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
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
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		c.Abort()
	}

	//TODO check uploadToken

	//hash and validation
	fileBinary, err := c.FormFile(paramFileBinary)
	hash, err := ulr.hash(fileBinary)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	if fileHash := c.PostForm(paramFileHash); fileHash != hash {
		c.JSON(http.StatusOK, response.Fail(fileCorrupted))
		return
	}

	//judge file is exist?
	fileInfo, err := queryFileByHash(hash)
	var size int64
	var resourceName, savePath string
	if err != nil {
		//file is not exist

		//gain save path TODO: from database
		savePath, err = ulr.gainSavePath("persist")
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusOK, response.Fail(err.Error()))
			return
		}

		//save file
		resourceName = util.UUID()
		size, err = ulr.saveFile(fileBinary, savePath, resourceName)
		if err != nil {
			logger.Error(err.Error())
			c.JSON(http.StatusOK, response.Fail(err.Error()))
			return
		}
	} else {
		//file is existed
		resourceName = fileInfo.ResourceName
		size = fileInfo.Size
		savePath = fileInfo.ResourcePath
	}
	fileName := c.PostForm(paramFileName)

	//insert file info to database
	key := util.UUID()
	newRecord := model.NewFileForDataBase(1, key, fileName, savePath, resourceName, hash, size, model.Public, model.Uploaded)
	re := saveFileInfo(newRecord)
	if re {
		c.JSON(http.StatusOK, response.Fail(saveRecordError))
		return
	}
	c.JSON(http.StatusOK, response.SuccessWithContent(newSmallFUResult(fileName, hash, fileDownloadUrl+key, size)))

	//TODO back up
}

func BigFileUploadInit(c *gin.Context) {
	if err := bigFFUInitCheck(c); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}

	//TODO uploadToken validation

	// judge file is exist?
	hash := c.PostForm(paramFileHash)
	fileName := c.PostForm(paramFileName)
	fileInfo, err := queryFileByHash(hash)
	if err == nil {
		// crate logic link
		//insert file info to database
		key := util.UUID()
		newRecord := model.NewFileForDataBase(1, key, fileName, fileInfo.ResourcePath, fileInfo.ResourceName, fileInfo.HashMd5, fileInfo.Size, model.Public, model.Uploaded)
		re := saveFileInfo(newRecord)
		if re {
			c.JSON(http.StatusOK, response.Fail(saveRecordError))
			return
		}
		c.JSON(http.StatusOK, response.SuccessWithContent(newBigFUIResultIsExist(exist, fileName, fileInfo.HashMd5, fileDownloadUrl+key)))
	} else {
		//not exist
		//return the chunk's hash which is exist
		var existedChunkArray []string
		parentFileInfo := tempInstance.getParentFileInfo(hash)
		if parentFileInfo == nil {
			bigFileInfo, err := newBigFileInfo(fileName, hash, c.PostForm(paramChunkCount))
			if err != nil {
				logger.Error(err.Error())
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
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}

	//TODO check uploadToken

	//hash and validation
	chunkBinary, _ := c.FormFile(paramChunkBinary)
	chunkHash, err := ulr.hash(chunkBinary)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	if fileHash := c.PostForm(paramChunkHash); fileHash != chunkHash {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(fileCorrupted))
		return
	}
	//gain save path
	savePath, err := util.PathAdaptive("/resource/temp")
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//save chunk file
	size, err := ulr.saveFile(chunkBinary, savePath, chunkHash)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//add or update chunk upload info
	chunkInfo, err := newChunkInfo(chunkHash, c.PostForm(paramChunkIndex), c.PostForm(paramChunkStart), c.PostForm(paramChunkEnd))
	if err != nil {
		logger.Error(err.Error())
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
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//TODO validation upload token
	//get merge info by file hash
	fileHash := c.PostForm(paramFileHash)
	fileInfo := tempInstance.getParentFileInfo(fileHash)
	chunkInfoMap := *fileInfo.chunkInfoMap
	length := len(chunkInfoMap)
	//check is all chunks unloaded
	if fileInfo.chunkCount != uint16(length) {
		c.JSON(http.StatusOK, response.Fail(chunkCountError))
	}
	//map to sorted array
	chunkSortArray := make([]*ChunkInfo, length, length)
	for _, v := range chunkInfoMap {
		chunkSortArray[v.Index] = v
	}
	//gain save path TODO: from database
	savePath, err := ulr.gainSavePath("persist")
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//get temp save path
	tempPath, err := util.PathAdaptive("/resource/temp/")
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//io（name=uuid,check start and end）
	var resourceName = util.UUID()
	hash, size, err := ulr.mergeChunk(&chunkSortArray, savePath, resourceName, tempPath)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	//check hash value
	if fileHash != hash {
		c.JSON(http.StatusOK, response.Fail(fileCorrupted))
		return
	}

	//insert into database
	key := util.UUID()
	newRecord := model.NewFileForDataBase(1, key, fileInfo.name, savePath, resourceName, hash, size, model.Public, model.Uploaded)
	re := saveFileInfo(newRecord)
	if re {
		c.JSON(http.StatusOK, response.Fail(saveRecordError))
		return
	}

	//delete big file cache info
	tempInstance.deleteParentFileInfo(fileHash)

	c.JSON(http.StatusOK, response.SuccessWithContent(newSmallFUResult(fileInfo.name, hash, fileDownloadUrl+key, size)))

	//TODO back up
}

//校验权限

//获取储存根目录

//获取用户目录

//获取日期

//文件信息写入数据库

//保存文件

//发起备份任务
