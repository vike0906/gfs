package download

import (
	"gfs/app/db"
	"gfs/app/logger"
	"gfs/app/repository/model"
)

//query file info by key
func queryFileByKey(key string) (*model.File, error) {
	var file model.File
	if err := db.DataBase().Where("file_key =?", key).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func queryUserByAppKey(appKey string) (*model.User, error) {
	var user model.User
	if err := db.DataBase().Where("app_key = ?", appKey).First(&user).Error; err != nil {
		logger.Error(err.Error)
		return nil, err
	}
	return &user, nil
}
