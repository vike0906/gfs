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
	bufSize = 128 * 1024
)

type uploader interface {
	gainSavePath(persistFolder string) (string, error)

	saveFile(file *multipart.FileHeader, path, resourceName string) (int64, error)

	mergeChunk(chunkSortArray *[]*ChunkInfo, path, resourceName, tempPath string) (string, int64, error)

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

func (u *uploadHelper) saveFile(file *multipart.FileHeader, path, resourceName string) (int64, error) {

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

func (u *uploadHelper) mergeChunk(chunkSortArray *[]*ChunkInfo, path, resourceName, tempPath string) (string, int64, error) {
	//创建待写入文件
	var resource string
	if osType := runtime.GOOS; osType == "windows" {
		resource = path + "\\" + resourceName
	} else {
		resource = path + "/" + resourceName
	}
	file, err := os.OpenFile(resource, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", 0, &common.GfsError{Message: "create file failed"}
	}
	defer file.Close()
	var size int64
	hash := md5.New()
	for _, v := range *chunkSortArray {
		tempResource := tempPath + v.Hash
		chunkFile, err := os.OpenFile(tempResource, os.O_RDONLY, 0600)
		//检查当前文件size
		if err != nil {
			chunkFile.Close()
			return "", 0, &common.GfsError{Message: "open temp file failed"}
		} else {
			var buf = make([]byte, bufSize)
			if written, err := io.CopyBuffer(file, chunkFile, buf); err != nil {
				chunkFile.Close()
				if err := os.Remove(resource); err != nil {
					log.Printf("delete temp main file:%s error", resource)
				}
				return "", 0, &common.GfsError{Message: "copy temp file failed"}
			} else {
				chunkFile, _ := os.OpenFile(tempResource, os.O_RDONLY, 0600)
				if _, err := io.CopyBuffer(hash, chunkFile, buf); err != nil {
					chunkFile.Close()
					return "", 0, &common.GfsError{Message: "copy temp file for hash failed"}
				}
				chunkFile.Close()
				size = size + written
			}
		}
		chunkFile.Close()

		if err := os.Remove(tempResource); err != nil {
			log.Println(err.Error())
			log.Printf("delete temp chunk file:%s error", resource)
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), size, nil
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
	fileHash := hex.EncodeToString(hash.Sum(nil))
	hash.Reset()
	return fileHash, nil
}
