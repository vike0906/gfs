package manager

import (
	"gfs/app/common"
	"gfs/app/db"
	"gfs/app/logger"
	"gfs/app/repository/model"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ResourceGain(c *gin.Context) {

	var userVo *model.UserVo

	if userVo = getUserVoByToken(c); userVo == nil {
		return
	}

	sqlConnect := db.DataBase()

	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "15"))
	direction := c.DefaultQuery("direction", "ascend")

	resources := make([]*model.FileVo, 0)
	var total int32
	sqlConnect = sqlConnect.Table("gfs_file").Select("gfs_file.*, gfs_user.name as user_name")
	if userVo.Role == 2 {
		sqlConnect = sqlConnect.Where("gfs_file.user_id = ?", userVo.ID)
	}

	if queryStr, isExist := c.GetQuery("queryStr"); isExist {
		sqlConnect = sqlConnect.Where("gfs_file.file_key like ?", "%"+queryStr+"%")
		sqlConnect = sqlConnect.Or("gfs_file.file_name like ?", "%"+queryStr+"%")
	}
	if order, isExist := c.GetQuery("order"); isExist {
		order = util.SnakeString(order)
		if direction == "descend" {
			order = order + " desc"
		}
		sqlConnect = sqlConnect.Order("gfs_file." + order)
	} else {
		sqlConnect = sqlConnect.Order("gfs_file.created_at desc")
	}
	sqlConnect = sqlConnect.Joins("LEFT JOIN gfs_user ON gfs_user.id = gfs_file.user_id")
	if err := sqlConnect.Count(&total).Error; err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(databaseError))
		return
	}

	if err := sqlConnect.Limit(limit).Offset((pageNo - 1) * limit).Scan(&resources).Error; err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(databaseError))
		return
	}

	if userVo.Role == 2 {
		for _, v := range resources {
			v.ResourcePath = ""
			v.ResourceName = ""
		}
	}
	c.JSON(http.StatusOK, response.SuccessWithContent(common.NewPage(total, &resources)))
}

func ResourceDelete(c *gin.Context) {
	if id, _ := strconv.Atoi(c.DefaultQuery("id", "0")); id <= 0 {
		c.JSON(http.StatusOK, response.Fail(paramError))
		return
	} else {
		sqlConnect := db.DataBase()
		if err := sqlConnect.Delete(&model.File{}, "id = ?", id).Error; err != nil {
			logger.Info(err.Error())
			c.JSON(http.StatusOK, response.Fail(ormOptionsFailed))
			return
		}
		c.JSON(http.StatusOK, response.SuccessWithMessage("Delete File Success"))
	}
}
