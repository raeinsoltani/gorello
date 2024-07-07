package models

import (
	"gorm.io/gorm"
)

type Workspace struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(100)"`
}
