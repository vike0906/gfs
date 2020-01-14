package router

import (
	"gfs/app/handle/download"
	"gfs/app/handle/upload"
	"gfs/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

var (
	g errgroup.Group
)

func Start(server *config.Server) {

	apiServer := &http.Server{
		Addr:         server.ApiPort,
		Handler:      api(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	webServer := &http.Server{
		Addr:         server.WebPort,
		Handler:      web(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return apiServer.ListenAndServe()
	})

	g.Go(func() error {
		return webServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}
func api() http.Handler {
	router := gin.New()
	router.MaxMultipartMemory = 2 << 20
	router.Use(gin.Recovery())
	//gin.SetMode(gin.ReleaseMode)
	router.Use(logger())
	router.Use(cors())
	router.NoRoute(handleNotFound)
	router.NoMethod(handleNotFound)

	//文件服务
	router.GET("/download/:key", download.Download)
	router.POST("/upload", upload.SmallFileUpload)
	router.POST("/init", upload.BigFileUploadInit)
	router.POST("/chunk", upload.BigFileUploadChunk)
	router.POST("/merge", upload.BigFileUploadMerge)

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
	return router
}

func web() http.Handler {
	router := gin.New()
	router.MaxMultipartMemory = 1 << 20
	gin.SetMode(gin.ReleaseMode)
	router.Use(gin.Recovery())
	router.StaticFile("/index", "vue/index.html")
	router.Static("/static", "./vue/static")
	router.NoRoute(func(c *gin.Context) {
		c.Request.URL.Path = "/index"
		router.HandleContext(c)
	})
	router.NoMethod(handleNotFound)
	return router
}
