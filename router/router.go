package router

import (
	fileHandle "gfs/app/handle/file"
	"github.com/gin-gonic/gin"
)

func Start(p *string) {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors())
	//账户验证
	accredit := router.Group("/accredit")
	{
		accredit.POST("/download")
		accredit.POST("/upload")
	}

	//文件服务
	file := router.Group("/file")
	{
		file.GET("/download", fileHandle.DownloadHandle)
		file.POST("/upload", fileHandle.UploadHandle)
	}

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

	//vue
	router.StaticFile("/index", "vue/index.html")
	router.Static("/static", "./vue/static")
	router.Run(":" + *p)
}
