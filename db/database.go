package db

import (
	"fmt"
	"gfs/app/repository"
	"gfs/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DataBase *gorm.DB

func Init(m *config.Mysql) error {

	var err error

	DataBase, err = gorm.Open("mysql", connectUrl(m))
	if err != nil {
		return err
	}

	DataBase.DB().SetMaxOpenConns(100)
	DataBase.DB().SetMaxIdleConns(10)
	DataBase.SingularTable(true)
	//DataBase.LogMode(true)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "gfs_" + defaultTableName
	}

	err = repository.CheckAndCreates(DataBase)
	if err != nil {
		return err
	}

	return nil
}

func connectUrl(m *config.Mysql) string {
	//var url = "userName:password@(host:port)/dbname?charset=utf8&parseTime=True&loc=Local"
	//url = strings.Replace(url,"userName",m.UserName,1)
	//url = strings.Replace(url,"password",m.Password,1)
	//url = strings.Replace(url,"host",m.Host,1)
	//url = strings.Replace(url,"port",m.Port,1)
	//url = strings.Replace(url,"dbname",m.DataBase,1)
	var url = "%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	url = fmt.Sprintf(url, m.UserName, m.Password, m.Host, m.Port, m.DataBase)
	return url
}
