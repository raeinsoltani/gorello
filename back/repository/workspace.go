package repository

import (
	"github.com/raeinsoltani/gorello/back/models"
)

type Workspace interface {
	Create(workspace *models.Workspace) error
	FindByID(id uint) (*models.Workspace, error)
	FindByName(name string) (*models.Workspace, error)
	Update(workspace *models.Workspace) error
	Delete(id uint) error
}
