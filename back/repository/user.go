package repository

import (
	"github.com/raeinsoltani/gorello/back/models"
)

type UserSearchResultDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type User interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByKeyWord(keyword string) ([]*UserSearchResultDTO, error)
	Update(user *models.User) error
	Delete(username string) error
}
