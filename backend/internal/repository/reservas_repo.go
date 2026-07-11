package repository

import (
	"errors"

	"gorm.io/gorm"

	"restaurantos/internal/models"
)

type ReservaRepo interface {
	Create(r *models.Reserva) error
	FindByID(id int) (*models.Reserva, error)
	FindAll(page, perPage int, estado string) ([]models.Reserva, int64, error)
	UpdateEstado(id int, estado string) error
	VincularPedido(id int, pedidoID int) error
}

type reservaRepoDB struct {
	db *gorm.DB
}

func NewReservaRepo(db *gorm.DB) ReservaRepo {
	return &reservaRepoDB{db: db}
}

func (r *reservaRepoDB) Create(reserva *models.Reserva) error {
	return r.db.Create(reserva).Error
}

func (r *reservaRepoDB) FindByID(id int) (*models.Reserva, error) {
	var reserva models.Reserva
	err := r.db.Where("id = ?", id).First(&reserva).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &reserva, err
}

func (r *reservaRepoDB) FindAll(page, perPage int, estado string) ([]models.Reserva, int64, error) {
	var reservas []models.Reserva
	var total int64

	q := r.db.Model(&models.Reserva{})
	if estado != "" {
		q = q.Where("estado = ?", estado)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := q.Order("fecha ASC, created_at DESC").Limit(perPage).Offset(offset).Find(&reservas).Error
	return reservas, total, err
}

func (r *reservaRepoDB) UpdateEstado(id int, estado string) error {
	result := r.db.Model(&models.Reserva{}).Where("id = ?", id).Update("estado", estado)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *reservaRepoDB) VincularPedido(id int, pedidoID int) error {
	result := r.db.Model(&models.Reserva{}).Where("id = ?", id).Update("pedido_id", pedidoID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
