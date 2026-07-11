package repository

import (
	"errors"

	"gorm.io/gorm"

	"restaurantos/internal/models"
)

type ClienteRepo interface {
	Create(c *models.Cliente) error
	FindByID(id int) (*models.Cliente, error)
	FindAll(page, perPage int, tipo string) ([]models.Cliente, int64, error)
	Update(id int, fields map[string]interface{}) (*models.Cliente, error)
	Deactivate(id int) error
	EmailExists(email string, excludeID int) bool
}

type clienteRepoDB struct {
	db *gorm.DB
}

func NewClienteRepo(db *gorm.DB) ClienteRepo {
	return &clienteRepoDB{db: db}
}

func (r *clienteRepoDB) Create(c *models.Cliente) error {
	return r.db.Create(c).Error
}

func (r *clienteRepoDB) FindByID(id int) (*models.Cliente, error) {
	var c models.Cliente
	err := r.db.Where("id = ? AND activo = true", id).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *clienteRepoDB) FindAll(page, perPage int, tipo string) ([]models.Cliente, int64, error) {
	var clientes []models.Cliente
	var total int64

	q := r.db.Model(&models.Cliente{}).Where("activo = true")
	if tipo != "" {
		q = q.Where("tipo = ?", tipo)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := q.Order("nombre ASC").Limit(perPage).Offset(offset).Find(&clientes).Error
	return clientes, total, err
}

func (r *clienteRepoDB) Update(id int, fields map[string]interface{}) (*models.Cliente, error) {
	if err := r.db.Model(&models.Cliente{}).Where("id = ? AND activo = true", id).Updates(fields).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *clienteRepoDB) Deactivate(id int) error {
	result := r.db.Model(&models.Cliente{}).Where("id = ? AND activo = true", id).Update("activo", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *clienteRepoDB) EmailExists(email string, excludeID int) bool {
	if email == "" {
		return false
	}
	var count int64
	r.db.Model(&models.Cliente{}).
		Where("email = ? AND id != ? AND activo = true", email, excludeID).
		Count(&count)
	return count > 0
}
