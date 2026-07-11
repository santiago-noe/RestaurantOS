package repository

import (
	"gorm.io/gorm"

	"restaurantos/internal/models"
)

type MenuRepo interface {
	FindPublico() ([]models.MenuPublico, error)
	FindAll() ([]models.MenuPublico, error)
	FindByID(id int) (*models.MenuPublico, error)
	Create(m *models.MenuPublico) error
	Update(id int, fields map[string]interface{}) (*models.MenuPublico, error)
	Delete(id int) error
}

type menuRepoDB struct {
	db *gorm.DB
}

func NewMenuRepo(db *gorm.DB) MenuRepo {
	return &menuRepoDB{db: db}
}

func (r *menuRepoDB) FindPublico() ([]models.MenuPublico, error) {
	var items []models.MenuPublico
	err := r.db.Where("disponible = true").Order("orden ASC, nombre ASC").Find(&items).Error
	return items, err
}

func (r *menuRepoDB) FindAll() ([]models.MenuPublico, error) {
	var items []models.MenuPublico
	err := r.db.Order("orden ASC, nombre ASC").Find(&items).Error
	return items, err
}

func (r *menuRepoDB) FindByID(id int) (*models.MenuPublico, error) {
	var m models.MenuPublico
	err := r.db.Where("id = ?", id).First(&m).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &m, err
}

func (r *menuRepoDB) Create(m *models.MenuPublico) error {
	return r.db.Create(m).Error
}

func (r *menuRepoDB) Update(id int, fields map[string]interface{}) (*models.MenuPublico, error) {
	if err := r.db.Model(&models.MenuPublico{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *menuRepoDB) Delete(id int) error {
	result := r.db.Delete(&models.MenuPublico{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
