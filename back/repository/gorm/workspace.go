package gorm

import (
	"errors"

	"github.com/raeinsoltani/gorello/back/models"
	"gorm.io/gorm"
)

type Workspace struct {
	db *gorm.DB
}

func NewWorkspaceRepo(db *gorm.DB) *Workspace {
	return &Workspace{db: db}
}

func (repo *Workspace) Create(workspace *models.Workspace) error {
	result := repo.db.Create(workspace)
	return result.Error
}

func (repo *Workspace) FindByID(id uint) (*models.Workspace, error) {
	var workspace models.Workspace
	result := repo.db.First(&workspace, "id = ?", id)
	return &workspace, result.Error
}

func (repo *Workspace) FindByName(name string) (*models.Workspace, error) {
	var workspace models.Workspace
	result := repo.db.First(&workspace, "name = ?", name)
	if errors.Is(gorm.ErrRecordNotFound, result.Error) {
		return nil, nil
	}
	return &workspace, result.Error
}

func (repo *Workspace) Update(workspace *models.Workspace) error {
	result := repo.db.Save(workspace)
	return result.Error
}

func (repo *Workspace) Delete(name string) error {
	result := repo.db.Where("username = ?", name).Delete(&models.Workspace{})
	return result.Error
}
