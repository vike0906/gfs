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

const format = "/2006/01/02"

type uploader interface {
	gainSavePath(userFolder string) (string, error)

	saveFile(file *multipart.FileHeader, path string, hash string) (int64, error)

	hash(file *multipart.FileHeader) (string, error)
}

type uploadHelper struct {
}

func (u *uploadHelper) gainSavePath(userFolder string) (string, error) {
	userFolder = "/resource/" + userFolder
	var dateHelper = util.NewDateHelper()
	var dateFolder = dateHelper.Format(time.Now(), format)
	savePath, err := util.PathAdaptive(userFolder + dateFolder)
	if err != nil {
		return "", err
	}
	return savePath, nil
}

func (u *uploadHelper) saveFile(file *multipart.FileHeader, path string, hash string) (int64, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				log.Println(err.Error())
				return 0, &common.GfsError{Message: "directory create failed"}
			}
		}
	}

	src, err := file.Open()
	if err != nil {
		return 0, &common.GfsError{Message: "file parse failed"}
	}
	defer src.Close()

	var resource string
	if osType := runtime.GOOS; osType == "windows" {
		resource = path + "\\" + hash
	} else {
		resource = path + "/" + hash
	}
	out, err := os.Create(resource)
	if err != nil {
		return 0, &common.GfsError{Message: "create file failed"}
	}
	defer out.Close()
	os.Chmod(resource, 0600)

	return io.Copy(out, src)
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
