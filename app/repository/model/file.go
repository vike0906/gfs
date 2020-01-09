package model

import "github.com/jinzhu/gorm"

type File struct {
	gorm.Model
	UserId       string `gorm:"default:'guest';not null"`
	Name         string `gorm:"size:10;not null"`
	ResourcePath string `gorm:"unique_index;not null"`
	ResourceName string `gorm:"size:64;not null"`
	HashMd5      string `gorm:"size:10;not null"`
	Size         string `gorm:"size:10;not null"`
	Type         uint8  `gorm:"default:0;not null"` //1:public 2:private
	Status       uint8  `gorm:"default:0;not null"` //1:uploaded 2:backupIng 3:backupEd 4:deleted
}
