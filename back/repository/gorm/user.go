package gorm

import (
	"errors"

	"github.com/raeinsoltani/gorello/back/models"
	"github.com/raeinsoltani/gorello/back/repository"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *User {
	return &User{db: db}
}

func (repo *User) Create(user *models.User) error {
	result := repo.db.Create(user)
	return result.Error
}

func (repo *User) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := repo.db.First(&user, "id = ?", id)
	return &user, result.Error
}

func (repo *User) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := repo.db.First(&user, "username = ?", username)
	if errors.Is(gorm.ErrRecordNotFound, result.Error) {
		return nil, nil
	}
	return &user, result.Error
}

func (repo *User) FindByKeyWord(keyword string) ([]*repository.UserSearchResultDTO, error) {
	var users []*repository.UserSearchResultDTO
	result := repo.db.Model(&models.User{}).Select("id, username, email").
		Where("email ILIKE ? OR username ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Scan(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (repo *User) Update(user *models.User) error {
	result := repo.db.Save(user)
	return result.Error
}

func (repo *User) Delete(username string) error {
	result := repo.db.Where("username = ?", username).Delete(&models.User{})
	return result.Error
}

func (repo *User) FindAll() ([]*models.User, error) {
	var users []*models.User
	result := repo.db.Find(&users)
	return users, result.Error
}
