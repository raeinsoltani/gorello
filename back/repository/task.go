package repository

import "github.com/raeinsoltani/gorello/back/models"

type Task interface {
	Create(task *models.Task) error
	FindByID(id uint) (*models.Task, error)
	FindByWorkspaceID(id uint) ([]*models.Task, error)
	Update(task *models.Task) error
	Delete(id uint) error
}
