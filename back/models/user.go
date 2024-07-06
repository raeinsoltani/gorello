package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;type:varchar(100)"`
	Email    string `gorm:"unique;type:varchar(100)"`
	Password string `gorm:"type:varchar(100)"`
}
