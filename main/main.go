package main

import (
	"fmt"
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
	//初始化数据库连接
	err = db.Init(&config.Info.Mysql)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//启动路由服务
	router.Start(&config.Info.Server.Port)

	fmt.Println("GFS Start Success")
}
