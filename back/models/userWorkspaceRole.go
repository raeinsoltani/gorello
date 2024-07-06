package models

type UserWorkspaceRole struct {
	User_id      uint `gorm:"foreignKey:UserID;not null"`
	Workspace_id uint `gorm:"foreignKey:WorkspaceID;not null"`
	Role         uint `gorm:"default:0"`
}
