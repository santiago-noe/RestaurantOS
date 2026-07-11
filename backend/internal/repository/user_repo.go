package repository

import (
	"errors"

	"gorm.io/gorm"

	"restaurantos/internal/models"
)

type UserRepo interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(id int) (*models.User, error)
}

type userRepoDB struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepoDB{db: db}
}

func (r *userRepoDB) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where("email = ? AND activo = true", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *userRepoDB) FindByID(id int) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}
