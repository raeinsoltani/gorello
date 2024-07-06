package repository

import (
	"github.com/raeinsoltani/gorello/back/models"
)

// UserWorkspaceRoleRepository defines the interface for user workspace role operations
type UserWorkspaceRoleRepository interface {
	Create(userWorkspaceRole *models.UserWorkspaceRole) error
	FindByID(id uint) (*models.UserWorkspaceRole, error)
	FindByUserID(userID uint) (*models.UserWorkspaceRole, error)
	FindByWorkspaceID(workspaceID uint) (*models.UserWorkspaceRole, error)
}
