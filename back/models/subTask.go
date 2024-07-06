package models

import (
	"gorm.io/gorm"
)

type SubTask struct {
	gorm.Model
	Title        string `gorm:"type:varchar(100);not null"`
	Task_id      uint   `gorm:"foreignKey:TaskID"`
	Is_completed bool   `gorm:"default:false"`
	Assignee_id  uint   `gorm:"foreignKey:AssigneeID;optional"`
}
