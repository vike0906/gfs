package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"default:'guest';not null"`
	LoginName string `gorm:"unique_index;not null"`
	Password  string `gorm:"size:64;not null"`
	Salt      string `gorm:"size:10;not null"`
	Role      uint8  `gorm:"default:0;not null"`
	Status    uint8  `gorm:"default:0;not null"`
	Token     string `gorm:"size:32"`
}

type UserVo struct {
	gorm.Model
	Name      string
	LoginName string
	Role      uint8
	Status    uint8
	Token     string
}

func (u *User) GainVo() *UserVo {
	var userVo = new(UserVo)
	userVo.Name = u.Name
	userVo.LoginName = u.LoginName
	userVo.Role = u.Role
	userVo.Status = u.Status
	userVo.Token = u.Token
	userVo.ID = u.ID
	userVo.CreatedAt = u.CreatedAt
	userVo.UpdatedAt = u.UpdatedAt
	return userVo
}

////query by loginName
//func (u *User)QueryByLoginName(loginName *string) (*User, error) {
//	var user User
//	if err := db.DataBase.Where("login_name = ?", loginName).First(&user).Error; err!=nil{
//		return nil, err
//	}
//	return &user,nil
//}
////add User
//func (u *User)Add(user *User) (*User, error) {
//	if err := db.DataBase.Create(user).Error; err!=nil{
//		return nil, err
//	}
//	return user, nil
//}
