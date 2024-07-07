package repository

import (
	"github.com/raeinsoltani/gorello/back/models"
)

type UserWorkspaceRole interface {
	Create(userWorkspaceRole *models.UserWorkspaceRole) error
	FindByID(id uint) (*models.UserWorkspaceRole, error)
	FindByUserID(user_id uint) ([]*models.UserWorkspaceRole, error)
	FindByWorkspaceID(workspace_id uint) ([]*models.UserWorkspaceRole, error)
}
