package upload

import (
	"crypto/md5"
	"encoding/hex"
	"gfs/app/common"
	"gfs/app/util"
	"io"
	"log"
	"mime/multipart"
	"os"
	"runtime"
	"time"
)

const (
	format  = "/2006/01/02"
	bufSize = 64 * 1024
)

type uploader interface {
	gainSavePath(persistFolder string) (string, error)

	saveFile(file *multipart.FileHeader, path string, resourceName string) (int64, error)

	mergeChunk(chunkSortArray []*ChunkInfo, path string, resourceName string) (int64, error)

	hash(file *multipart.FileHeader) (string, error)
}

type uploadHelper struct {
}

func (u *uploadHelper) gainSavePath(persistFolder string) (string, error) {
	persistFolder = "/resource/" + persistFolder
	var dateHelper = util.NewDateHelper()
	var dateFolder = dateHelper.Format(time.Now(), format)
	savePath, err := util.PathAdaptive(persistFolder + dateFolder)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(savePath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
				log.Println(err.Error())
				return "", &common.GfsError{Message: "directory create failed"}
			}
		}
	}
	return savePath, nil
}

func (u *uploadHelper) saveFile(file *multipart.FileHeader, path string, resourceName string) (int64, error) {

	src, err := file.Open()
	if err != nil {
		return 0, &common.GfsError{Message: "file parse failed"}
	}
	defer src.Close()

	var resource string
	if osType := runtime.GOOS; osType == "windows" {
		resource = path + "\\" + resourceName
	} else {
		resource = path + "/" + resourceName
	}
	out, err := os.Create(resource)
	if err != nil {
		return 0, &common.GfsError{Message: "create file failed"}
	}
	defer out.Close()
	os.Chmod(resource, 0600)

	return io.Copy(out, src)
}

func (u *uploadHelper) mergeChunk(chunkSortArray []*ChunkInfo, path string, resourceName, tempPath string) (int64, error) {
	//创建待写入文件
	var resource string
	if osType := runtime.GOOS; osType == "windows" {
		resource = path + "\\" + resourceName
	} else {
		resource = path + "/" + resourceName
	}
	file, err := os.OpenFile(resource, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0600)
	if err != nil {
		return 0, &common.GfsError{Message: "create file failed"}
	}
	for _, v := range chunkSortArray {
		var buf = make([]byte, bufSize)
		tempPath = tempPath + v.Hash
		chunkFile, err := os.OpenFile(tempPath, os.O_RDONLY, 0600)
		if err != nil {
			return 0, &common.GfsError{Message: "open temp file failed"}
		} else {
			io.CopyBuffer(file, chunkFile, buf)
		}
		chunkFile.Close()
	}
	defer file.Close()
	//检查当前文件size
	//继续写入
	//
	return io.Copy(file, file)
}
func (u *uploadHelper) hash(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", &common.GfsError{Message: "file parse failed"}
	}
	defer src.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", &common.GfsError{Message: "create hash(md5) fail"}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
