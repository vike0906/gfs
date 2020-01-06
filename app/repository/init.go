package repository

import (
	"gfs/app/repository/model"
	"github.com/jinzhu/gorm"
)

//检查更新
func CheckAndCreates(db *gorm.DB) error {
	var err error
	if err = checkAndCreate(db, &model.User{}); err != nil {
		return err
	}
	return nil
}

func checkAndCreate(db *gorm.DB, tableModel interface{}) error {
	if !db.HasTable(tableModel) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(tableModel).Error; err != nil {
			return err
		}
	}
	return nil
}
