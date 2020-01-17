package model

import "github.com/jinzhu/gorm"

const (
	Public    = 1
	Private   = 2
	Uploaded  = 1
	BackUpIng = 2
	BackUpEnd = 3
	Deleted   = 4
)

type File struct {
	gorm.Model
	UserId       uint   `gorm:"not null"`
	FileKey      string `gorm:"size:64;unique_index;not null"`
	FileName     string `gorm:"size:255;not null"`
	ResourcePath string `gorm:"size:255;not null"`
	ResourceName string `gorm:"size:64;not null"`
	HashMd5      string `gorm:"size:32;index;not null"`
	Size         int64  `gorm:"default:0;not null"`
	Type         uint8  `gorm:"default:0;not null"` //1:public 2:private
	Status       uint8  `gorm:"default:0;not null"` //1:uploaded 2:backupIng 3:backupEd 4:deleted
}

type FileVo struct {
	gorm.Model
	UserId       uint   `gorm:"not null"`
	FileKey      string `gorm:"size:64;unique_index;not null"`
	FileName     string `gorm:"size:255;not null"`
	ResourcePath string `gorm:"size:255;not null"`
	ResourceName string `gorm:"size:64;not null"`
	HashMd5      string `gorm:"size:32;index;not null"`
	Size         int64  `gorm:"default:0;not null"`
	Type         uint8  `gorm:"default:0;not null"` //1:public 2:private
	Status       uint8  `gorm:"default:0;not null"` //1:uploaded 2:backupIng 3:backupEd 4:deleted
	UserName     string
}

func NewFileForDataBase(userId uint, key, name, resourcePath, resourceName, hashMd5 string, size int64, fileType, status uint8) *File {
	return &File{UserId: userId, FileKey: key, FileName: name, ResourcePath: resourcePath, ResourceName: resourceName, HashMd5: hashMd5, Size: size, Type: fileType, Status: status}
}
