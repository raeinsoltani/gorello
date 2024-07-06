package models

import (
	"gorm.io/gorm"
)

type Workspace struct {
	gorm.Model
	Name        string `gorm:"unique;type:varchar(100)"`
	Description string `gorm:"type:varchar(100)"`
}
