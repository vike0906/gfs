package upload

import (
	"gfs/app/db"
	"gfs/app/repository/model"
)

//check file is exist
func queryFileByHash(hash string) (*model.File, error) {
	var file model.File
	if err := db.DataBase.Where("hash_md5 =?", hash).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

//save file info
func saveFileInfo(file *model.File) bool {
	db.DataBase.Create(file)
	return db.DataBase.NewRecord(file)
}
