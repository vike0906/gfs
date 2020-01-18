package router

import (
	"fmt"
	"gfs/app/common"
	"gfs/app/component"
	log "gfs/app/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AuthToken")
		c.Header("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//OPTIONS pass
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// handle request
		c.Next()
	}
}

func logger() gin.HandlerFunc {
	c := gin.LoggerConfig{
		Output:    *log.GetLogWriter(),
		SkipPaths: []string{"/vue"},
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("[GIN]%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				params.ClientIP,
				params.TimeStamp.Format(time.RFC1123),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
			)
		},
	}
	return gin.LoggerWithConfig(c)
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if "/manager/login" == c.Request.RequestURI {
			c.Next()
			return
		}
		authToken := c.Request.Header.Get("AuthToken")
		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusOK, common.Response{Code: 100, Message: "Illegal request"})
			return
		}
		userVo := component.GetAuthToken(authToken)
		if userVo == nil {
			c.AbortWithStatusJSON(http.StatusOK, common.Response{Code: 100, Message: "AuthToken is Expired"})
			return
		}
		component.PutAuthToken(authToken, userVo)
		c.Next()
	}
}

func handleNotFound(c *gin.Context) {
	c.JSON(http.StatusOK, common.ResponseInstance().Fail("Not Found Error"))
	return
}
