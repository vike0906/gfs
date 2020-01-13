package download

import (
	"gfs/app/repository/model"
	"gfs/db"
)

//query file info by key
func queryFileByKey(key string) (*model.File, error) {
	var file model.File
	if err := db.DataBase.Where("file_key =?", key).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}
