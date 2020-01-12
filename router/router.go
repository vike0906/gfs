package router

import (
	"gfs/app/handle/download"
	"gfs/app/handle/upload"
	"github.com/gin-gonic/gin"
)

func Start(p *string) {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors())

	//文件服务
	router.GET("/download", download.HandleDownload)
	router.POST("/upload", upload.SmallFileUpload)
	router.POST("/init", upload.BigFileUploadInit)
	router.POST("/chunk", upload.BigFileUploadChunk)
	router.POST("/merge", upload.BigFileUploadMerge)

	//账户授权
	accreditServer := router.Group("/accredit")
	{
		accreditServer.POST("/upload")
		accreditServer.POST("/download")

	}

	//系统管理
	managerServer := router.Group("/manager")
	{
		managerServer.POST("/login")
		managerServer.POST("/logout")
		managerServer.POST("/change-psd")

		managerServer.GET("/users")
		user := managerServer.Group("/user")
		user.POST("/gain")
		user.POST("/add")
		user.POST("/edit")
		user.POST("/delete")

		file := managerServer.Group("upload")
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
