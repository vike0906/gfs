package main

import (
	"gfs/app/db"
	"gfs/config"
	"gfs/router"
	"log"
)

func main() {

	//init config document
	err := config.Init()
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("init config info success")

	//init database connection pool
	err = db.Init(&config.Info.Mysql)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("init database success")

	//http server start
	router.Start(&config.Info.Server)
}
