package gorm

import (
	"github.com/raeinsoltani/gorello/back/models"
	"gorm.io/gorm"
)

type UserWorkspaceRole struct {
	db *gorm.DB
}

func NewUserWorkspaceRoleRepo(db *gorm.DB) *UserWorkspaceRole {
	return &UserWorkspaceRole{db: db}
}

func (repo *UserWorkspaceRole) Create(UserWorkspaceRole *models.UserWorkspaceRole) error {
	result := repo.db.Create(UserWorkspaceRole)
	return result.Error
}

func (repo *UserWorkspaceRole) FindByID(id uint) (*models.UserWorkspaceRole, error) {
	var UserWorkspaceRole models.UserWorkspaceRole
	result := repo.db.First(&UserWorkspaceRole, "id = ?", id)
	return &UserWorkspaceRole, result.Error
}

func (repo *UserWorkspaceRole) FindByUserID(user_id uint) ([]*models.UserWorkspaceRole, error) {
	var userWorkspaceRoles []*models.UserWorkspaceRole
	result := repo.db.Find(&userWorkspaceRoles, "user_id = ?", user_id)
	return userWorkspaceRoles, result.Error
}

func (repo *UserWorkspaceRole) FindByWorkspaceID(workspace_id uint) ([]*models.UserWorkspaceRole, error) {
	var userWorkspaceRoles []*models.UserWorkspaceRole
	result := repo.db.Find(&userWorkspaceRoles, "workspace_id = ?", workspace_id)
	return userWorkspaceRoles, result.Error
}
