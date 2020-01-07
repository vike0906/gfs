package upload

import (
	"gfs/app/util"
	"time"
)

const format = "/2006/01/02/"

type uploader interface {
	gainSavePath(userFolder string) (*string, error)
}

type uploadHelper struct {
}

func (u *uploadHelper) gainSavePath(userFolder string) (*string, error) {
	userFolder = "/resource/" + userFolder
	var dateHelper = util.NewDateHelper()
	var dateFolder = dateHelper.Format(time.Now(), format)
	savePath, err := util.PathAdaptive(userFolder + dateFolder)
	if err != nil {
		return nil, err
	}
	return &savePath, nil
}
