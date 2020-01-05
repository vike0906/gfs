package router

import (
	"github.com/gin-gonic/gin"
)

func Start(p *string) {

	router := gin.Default()

	//系统管理
	manager := router.Group("/manager")
	{
		manager.POST("/login")
		manager.POST("/logout")
		manager.POST("/change-psd")

		manager.GET("/users")
		user := manager.Group("/user")
		user.POST("/gain")
		user.POST("/add")
		user.POST("/edit")
		user.POST("/delete")

		file := manager.Group("file")
		file.POST("/gain")
		file.POST("/add")
		file.POST("/edit")
		file.POST("/delete")
	}

	//账户验证
	account := router.Group("/accredit")
	{
		account.POST("/download")
		account.POST("/upload")
	}

	//文件服务
	file := router.Group("/file")
	{
		file.GET("/download")
		file.POST("/upload")
	}

	//vue
	router.StaticFile("/index", "vue/index.html")
	router.Static("/static", "./vue/static")
	router.Run(":" + *p)
}
