package router

import (
	"gfs/app/handle/download"
	"gfs/app/handle/manager"
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

	//file server
	router.GET("/download/:key", download.Download)
	router.POST("/upload", upload.SmallFileUpload)
	router.POST("/init", upload.BigFileUploadInit)
	router.POST("/chunk", upload.BigFileUploadChunk)
	router.POST("/merge", upload.BigFileUploadMerge)

	//manager api
	managerServer := router.Group("/manager")
	managerServer.Use(auth())
	{
		managerServer.POST("/login", manager.Login)
		managerServer.POST("/change-psd", manager.ChangePassword)
		managerServer.POST("/logout", manager.Logout)

		user := managerServer.Group("/user")
		user.GET("/gain", manager.UserGain)
		user.POST("/save", manager.UserSave)
		user.DELETE("/delete", manager.UserDelete)

		file := managerServer.Group("resource")
		file.GET("/gain", manager.ResourceGain)
		file.DELETE("/delete", manager.ResourceDelete)
		file.GET("/accredit/upload", manager.AccreditUpload)
		file.GET("/accredit/download", manager.AccreditDownload)

	}
	return router
}

func web() http.Handler {
	router := gin.New()
	router.MaxMultipartMemory = 1 << 20
	gin.SetMode(gin.ReleaseMode)
	router.Use(gin.Recovery())
	router.Static("/static", "./vue/static")
	router.StaticFile("/index", "vue/index.html")
	router.NoRoute(func(c *gin.Context) {
		c.Request.URL.Path = "/index"
		router.HandleContext(c)
	})
	router.NoMethod(handleNotFound)
	return router
}
