package main

import (
	"gfs/config"
	"gfs/db"
	"gfs/router"
	"log"
)

func main() {

	//初始化配置文件
	err := config.Init()
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("init config info success")

	//初始化数据库连接
	err = db.Init(&config.Info.Mysql)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("init database success")

	//启动路由服务
	router.Start(&config.Info.Server.Port)
}
