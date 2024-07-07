package gorm

import (
	"github.com/raeinsoltani/gorello/back/models"
	"gorm.io/gorm"
)

type Task struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *Task {
	return &Task{db: db}
}

func (repo *Task) Create(task *models.Task) error {
	result := repo.db.Create(task)
	return result.Error
}

func (repo *Task) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	result := repo.db.First(&task, "id = ?", id)
	return &task, result.Error
}

func (repo *Task) FindByWorkspaceID(id uint) ([]*models.Task, error) {
	var tasks []*models.Task
	result := repo.db.Find(&tasks, "workspace_id = ?", id)
	return tasks, result.Error
}

func (repo *Task) Update(task *models.Task) error {
	result := repo.db.Save(task)
	return result.Error
}

func (repo *Task) Delete(id uint) error {
	result := repo.db.Where("id = ?", id).Delete(&models.Task{})
	return result.Error
}
