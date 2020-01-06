package file

import (
	"gfs/app/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DownloadHandle(c *gin.Context) {
	var response = common.Response{}
	c.JSON(http.StatusOK, response.Success())
}
