package upload

import (
	"fmt"
	"gfs/app/common"
	"github.com/gin-gonic/gin"
)

func smallFUParamCheck(c *gin.Context) error {
	if err := paramCheck(c); err != nil {
		return err
	}
	if fileName := c.PostForm(paramFileName); fileName == "" {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramFileName))
		return &err
	}
	if fileBinary, err := c.FormFile(paramFileBinary); err != nil || fileBinary.Size == 0 {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramFileBinary))
		return &err
	}
	return nil
}
func bigFFUInitCheck(c *gin.Context) error {
	if err := paramCheck(c); err != nil {
		return err
	}
	if fileName := c.PostForm(paramFileName); fileName == "" {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramFileName))
		return &err
	}
	if chunkCount := c.PostForm(paramChunkCount); chunkCount == "" {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramChunkCount))
		return &err
	}
	return nil
}

func bigFFUChunkCheck(c *gin.Context) error {
	if err := paramCheck(c); err != nil {
		return err
	}
	if fileHashMap := tempInstance.getParentFileInfo(c.PostForm(paramFileHash)); fileHashMap == nil {
		err := common.NewGfsError("file hash is error, you need restart the upload")
		return &err
	}
	if chunkHash := c.PostForm(paramChunkHash); chunkHash == "" {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramChunkHash))
		return &err
	}
	if chunkIndex := c.PostForm(paramChunkIndex); chunkIndex == "" {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramChunkIndex))
		return &err
	}
	if chunkStart := c.PostForm(paramChunkStart); chunkStart == "" {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramChunkStart))
		return &err
	}
	if chunkEnd := c.PostForm(paramChunkEnd); chunkEnd == "" {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramChunkEnd))
		return &err
	}
	if chunkBinary, err := c.FormFile(paramChunkBinary); err != nil || chunkBinary.Size == 0 {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramChunkBinary))
		return &err
	}
	return nil
}

func bigFFUMergeCheck(c *gin.Context) error {
	if err := paramCheck(c); err != nil {
		return err
	}
	return nil
}

func paramCheck(c *gin.Context) error {

	if uploadToken := c.PostForm(paramUploadToken); uploadToken == "" {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramUploadToken))
		return &err
	}
	if fileHash := c.PostForm(paramFileHash); fileHash == "" {
		err := common.NewGfsError(fmt.Sprintf(errMessage, paramFileHash))
		return &err
	}
	return nil
}

func uploadTokenCheck(token *string) error {
	return nil
}
