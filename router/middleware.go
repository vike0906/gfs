package router

import (
	"gfs/app/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AuthToken")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
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

func handleNotFound(c *gin.Context) {
	c.JSON(http.StatusOK, common.ResponseInstance().Fail("Not Found Error"))
	return
}
