package db

import (
	"fmt"
	"gfs/app/repository"
	"gfs/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var dataBase *gorm.DB

func DataBase() *gorm.DB {
	db := dataBase
	return db
}

func Init(m *config.Mysql) error {

	var err error

	dataBase, err = gorm.Open("mysql", connectUrl(m))
	if err != nil {
		return err
	}

	dataBase.DB().SetMaxOpenConns(100)
	dataBase.DB().SetMaxIdleConns(10)
	dataBase.SingularTable(true)
	dataBase.LogMode(true)
	//dataBase.SetLogger(logger.GetLogger("ORM"))

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "gfs_" + defaultTableName
	}

	err = repository.CheckAndCreates(dataBase)
	if err != nil {
		return err
	}

	return nil
}

func connectUrl(m *config.Mysql) string {
	var url = "%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	url = fmt.Sprintf(url, m.UserName, m.Password, m.Host, m.Port, m.DataBase)
	return url
}
