package main

import (
	"fmt"
	"gfs/config"
	"gfs/router"
)

func main() {
	//初始化配置文件
	err := config.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//启动服务
	router.Start(&config.Info.Server.Port)
}
