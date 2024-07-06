package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title          string `gorm:"type:varchar(100);not null"`
	Description    string `gorm:"type:varchar(100)"`
	Status         uint   `gorm:"default:0"`
	Estimated_time string `gorm:"type:varchar(100)"`
	Actual_time    string `gorm:"type:varchar(100)"`
	Due_date       *time.Time
	Priority       uint   `gorm:"default:0"`
	Workspace_id   uint   `gorm:"foreignKey:WorkspaceID"`
	Assignee_id    uint   `gorm:"foreignKey:AssigneeID;optional"`
	Image_url      string `gorm:"type:varchar(100)"`
}
