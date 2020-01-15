package manager

import (
	"gfs/app/db"
	"gfs/app/logger"
	"gfs/app/repository/model"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserGain(c *gin.Context) {
	sqlConnect := db.DataBase()

	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	direction := c.DefaultQuery("direction", "ascend")

	users := make([]model.User, 0)

	if queryStr, isExist := c.GetQuery("queryStr"); isExist {
		sqlConnect = sqlConnect.Where("name like ?", "%"+queryStr+"%").Or("login_name like ?", "%"+queryStr+"%")
	}
	if order, isExist := c.GetQuery("order"); isExist {
		if direction == "descend" {
			order = util.SnakeString(order) + "desc"
		}
		sqlConnect = sqlConnect.Order("age desc")
	}

	if err := sqlConnect.Limit(limit).Offset((pageNo - 1) * limit).Find(&users).Error; err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(databaseError))
		return
	}
	c.JSON(http.StatusOK, response.SuccessWithContent(users))
}

func UserSave(c *gin.Context) {

}

func UserDelete(c *gin.Context) {

}
