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
		return
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
		//not exist return the chunk's hash which is exist
		var chunkInfoMap map[string]ChunkInfo
		bigFileInfo, err := newBigFileInfo(fileName, hash, c.PostForm(paramChunkCount), &chunkInfoMap)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, response.Fail(err.Error()))
			return
		}

		parentFileInfo := tempInstance.getParentFileInfo(hash)
		if parentFileInfo == nil {
			tempInstance.putBigFileInfo(hash, bigFileInfo)
		} else {
			chunkInfoMap = *parentFileInfo.chunkInfoMap
		}

		chunkInfoMapLength := len(chunkInfoMap)
		chunkArray := make([]string, chunkInfoMapLength, chunkInfoMapLength)
		for k := range chunkInfoMap {
			chunkArray = append(chunkArray, k)
		}
		c.JSON(http.StatusOK, response.SuccessWithContent(newBigFUIResultUnExist(unExit, &chunkArray)))

		//var chunkHashArray = c.PostForm(paramChunkInfoArray)
		//var cha []ChunkInfo
		//if err := json.Unmarshal([]byte(chunkHashArray),&cha);err!=nil{
		//	log.Println("解析失败"+err.Error())
		//}else {
		//	log.Println(cha)
		//}
		//log.Println(chunkHashArray)
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
	fileBinary, _ := c.FormFile(paramFileBinary)
	hash, err := ulr.hash(fileBinary)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	if fileHash := c.PostForm(paramFileHash); fileHash != hash {
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
	//save chunk info
	size, err := ulr.saveFile(fileBinary, savePath, hash)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	log.Println(size)

	//add or update chunk upload info
	chunkInfo, err := newChunkInfo(hash, c.PostForm(paramChunkIndex), c.PostForm(paramChunkStart), c.PostForm(paramChunkEnd))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	parentFileInfo := tempInstance.getParentFileInfo(c.PostForm(paramParentFileHash))
	var m = *parentFileInfo.chunkInfoMap
	m[hash] = *chunkInfo

	tempInstance.putChunkHash(hash, time.Now().Unix())

	c.JSON(http.StatusOK, response.SuccessWithMessage(""))
}

func BigFileUploadMerge(c *gin.Context) {

}

//校验权限

//获取储存根目录

//获取用户目录

//获取日期

//文件信息写入数据库

//保存文件

//发起备份任务
