package manager

import (
	"gfs/app/component"
	"gfs/app/db"
	"gfs/app/repository/model"
	"gfs/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginForm struct {
	LoginName string `form:"loginName" binding:"required"`
	Password  string `form:"password" binding:"required"`
}

type ChangePsdForm struct {
	OldPassword string `form:"oldPsd" binding:"required"`
	NewPassword string `form:"newPsd" binding:"required"`
}

func Login(c *gin.Context) {

	var form LoginForm

	if c.ShouldBind(&form) != nil {
		c.JSON(http.StatusOK, response.Fail(formParamError))
		return
	}

	//query user info by loginName
	var user model.User
	if err := db.DataBase().Where("login_name = ?", form.LoginName).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, response.Fail(loginNameOrPsdError))
		return
	}

	//calculate password and comparison
	p := passwordHash(form.Password, user.Salt)
	if p != user.Password {
		c.JSON(http.StatusOK, response.Fail(loginNameOrPsdError))
		return
	}

	//create token and put to token cache
	token := util.UUID()
	userVo := user.GainVo()
	userVo.Token = token
	component.PutAuthToken(token, userVo)
	vo := *userVo
	vo.ID = 0
	c.JSON(http.StatusOK, response.SuccessWithContent(vo))
}

func ChangePassword(c *gin.Context) {

	var changePsdForm ChangePsdForm

	if c.ShouldBind(&changePsdForm) != nil {
		c.JSON(http.StatusOK, response.Fail(formParamError))
		return
	}

	if userVo := getUserVoByToken(c); userVo == nil {
		return
	} else {

		var user model.User

		if err := db.DataBase().Where("id = ?", userVo.ID).First(&user).Error; err != nil {
			c.JSON(http.StatusOK, response.Fail(userInfoError))
			return
		} else {
			//calculate password and comparison
			p := passwordHash(changePsdForm.OldPassword, user.Salt)
			if p != user.Password {
				c.JSON(http.StatusOK, response.Fail("old password is error"))
				return
			}
			np := passwordHash(changePsdForm.NewPassword, user.Salt)
			if err := db.DataBase().Model(&user).Update("password", np).Error; err != nil {
				c.JSON(http.StatusOK, response.Fail(ormOptionsFailed))
				return
			} else {
				component.DeleteAuthToken(userVo.Token)
				c.JSON(http.StatusOK, response.SuccessWithMessage("change password success, please login again"))
			}
		}
	}
	//获取用户信息
}

func Logout(c *gin.Context) {

	var userVo *model.UserVo

	if userVo = getUserVoByToken(c); userVo == nil {
		return
	}

	component.DeleteAuthToken(userVo.Token)

	c.JSON(http.StatusOK, response.SuccessWithMessage("logout success, please login again"))
}
