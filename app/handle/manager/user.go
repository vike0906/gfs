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
		order = util.SnakeString(order)
		if direction == "descend" {
			order = order + " desc"
		}
		sqlConnect = sqlConnect.Order(order)
	}

	if err := sqlConnect.Limit(limit).Offset((pageNo - 1) * limit).Find(&users).Error; err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, response.Fail(databaseError))
		return
	}
	c.JSON(http.StatusOK, response.SuccessWithContent(users))
}

func UserSave(c *gin.Context) {

	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	name := c.Query("name")
	loginName := c.Query("loginName")
	roleId := c.Query("roleId")
	status := c.Query("status")

	if id <= 0 {
		//add user
		if name == "" || roleId == "" || status == "" {
			c.JSON(http.StatusOK, response.Fail(paramError))
			return
		}
		r, _ := strconv.Atoi(roleId)
		s, _ := strconv.Atoi(status)
		if err := addUser(name, loginName, uint8(r), uint8(s)); err != nil {
			c.JSON(http.StatusOK, response.Fail(ormOptionsFailed))
			return
		}
		c.JSON(http.StatusOK, response.SuccessWithMessage("Add User Success"))
	} else {
		//update user
		if name == "" || roleId == "" {
			c.JSON(http.StatusOK, response.Fail(paramError))
			return
		}
		r, _ := strconv.Atoi(roleId)
		if err := updateUser(id, name, uint8(r)); err != nil {
			c.JSON(http.StatusOK, response.Fail(ormOptionsFailed))
			return
		}
		c.JSON(http.StatusOK, response.SuccessWithMessage("Update User Success"))
	}
}

func UserDelete(c *gin.Context) {
	if id, _ := strconv.Atoi(c.DefaultQuery("id", "0")); id <= 0 {
		c.JSON(http.StatusOK, response.Fail(paramError))
		return
	} else {
		sqlConnect := db.DataBase()
		if err := sqlConnect.Delete(&model.User{}, "id = ?", id).Error; err != nil {
			logger.Info(err.Error())
			c.JSON(http.StatusOK, response.Fail(ormOptionsFailed))
			return
		}
		c.JSON(http.StatusOK, response.SuccessWithMessage("Delete User Success"))
	}
}

func addUser(name, loginName string, role, status uint8) error {
	sqlConnect := db.DataBase()
	salt := util.RandomString(10)
	p := passwordHash(initPassword, salt)
	appKey := util.UUID()
	appSecret := appSecretHash(appKey, salt)
	user := model.NewUser(name, loginName, p, salt, appKey, appSecret, role, status)
	if err := sqlConnect.Create(user).Error; err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func updateUser(id int, name string, status uint8) error {
	sqlConnect := db.DataBase()
	var user model.User
	if err := sqlConnect.First(&user, id).Error; err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := sqlConnect.Model(&user).Update(map[string]interface{}{"name": name, "status": status}).Error; err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
