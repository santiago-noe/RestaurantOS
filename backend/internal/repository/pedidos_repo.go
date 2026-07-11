package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"restaurantos/internal/models"
)

// ─── PedidoRepo ──────────────────────────────────────────────────────────────

type PedidoRepo interface {
	Create(pedido *models.Pedido, items []models.PedidoItem) error
	FindByID(id int) (*models.Pedido, error)
	FindAll(page, perPage, clienteID int, estado string) ([]models.Pedido, int64, error)
	UpdateEstado(id int, estado string) error
	FindEntreFechas(desde, hasta time.Time) ([]models.Pedido, error)
}

type pedidoRepoDB struct{ db *gorm.DB }

func NewPedidoRepo(db *gorm.DB) PedidoRepo { return &pedidoRepoDB{db: db} }

func (r *pedidoRepoDB) Create(pedido *models.Pedido, items []models.PedidoItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(pedido).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].PedidoID = pedido.ID
		}
		return tx.Create(&items).Error
	})
}

func (r *pedidoRepoDB) FindByID(id int) (*models.Pedido, error) {
	var p models.Pedido
	err := r.db.Preload("Cliente").Preload("Items").Preload("Items.Producto").
		First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *pedidoRepoDB) FindAll(page, perPage, clienteID int, estado string) ([]models.Pedido, int64, error) {
	var pedidos []models.Pedido
	var total int64

	q := r.db.Model(&models.Pedido{})
	if clienteID > 0 {
		q = q.Where("cliente_id = ?", clienteID)
	}
	if estado != "" {
		q = q.Where("estado = ?", estado)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := q.Preload("Cliente").Order("fecha DESC").Limit(perPage).Offset(offset).Find(&pedidos).Error
	return pedidos, total, err
}

func (r *pedidoRepoDB) FindEntreFechas(desde, hasta time.Time) ([]models.Pedido, error) {
	var pedidos []models.Pedido
	err := r.db.Where("fecha BETWEEN ? AND ?", desde, hasta).
		Order("fecha ASC").Find(&pedidos).Error
	return pedidos, err
}

func (r *pedidoRepoDB) UpdateEstado(id int, estado string) error {
	result := r.db.Model(&models.Pedido{}).Where("id = ?", id).
		Updates(map[string]interface{}{"estado": estado})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ─── ProductoRepo ─────────────────────────────────────────────────────────────

type ProductoRepo interface {
	FindByID(id int) (*models.Producto, error)
	AjustarStock(id int, cantidad float64) error // positivo=entrada, negativo=salida
	FindAll(soloActivos bool) ([]models.Producto, error)
	Create(p *models.Producto) error
	Update(id int, fields map[string]interface{}) (*models.Producto, error)
}

type productoRepoDB struct{ db *gorm.DB }

func NewProductoRepo(db *gorm.DB) ProductoRepo { return &productoRepoDB{db: db} }

func (r *productoRepoDB) FindByID(id int) (*models.Producto, error) {
	var p models.Producto
	err := r.db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *productoRepoDB) AjustarStock(id int, cantidad float64) error {
	return r.db.Model(&models.Producto{}).Where("id = ?", id).
		Update("stock_actual", gorm.Expr("stock_actual + ?", cantidad)).Error
}

func (r *productoRepoDB) FindAll(soloActivos bool) ([]models.Producto, error) {
	var productos []models.Producto
	q := r.db.Model(&models.Producto{})
	if soloActivos {
		q = q.Where("activo = true")
	}
	err := q.Order("nombre ASC").Find(&productos).Error
	return productos, err
}

func (r *productoRepoDB) Create(p *models.Producto) error {
	return r.db.Create(p).Error
}

func (r *productoRepoDB) Update(id int, fields map[string]interface{}) (*models.Producto, error) {
	if err := r.db.Model(&models.Producto{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

// ─── MovimientoStockRepo ──────────────────────────────────────────────────────

type MovimientoRepo interface {
	Registrar(m *models.MovimientoStock) error
	FindByProducto(productoID int) ([]models.MovimientoStock, error)
	FindEntreFechas(desde, hasta time.Time) ([]models.MovimientoStock, error)
}

type movimientoRepoDB struct{ db *gorm.DB }

func NewMovimientoRepo(db *gorm.DB) MovimientoRepo { return &movimientoRepoDB{db: db} }

func (r *movimientoRepoDB) Registrar(m *models.MovimientoStock) error {
	return r.db.Create(m).Error
}

func (r *movimientoRepoDB) FindByProducto(productoID int) ([]models.MovimientoStock, error) {
	var movs []models.MovimientoStock
	err := r.db.Where("producto_id = ?", productoID).Order("created_at DESC").Find(&movs).Error
	return movs, err
}

func (r *movimientoRepoDB) FindEntreFechas(desde, hasta time.Time) ([]models.MovimientoStock, error) {
	var movs []models.MovimientoStock
	err := r.db.Preload("Producto").
		Where("fecha BETWEEN ? AND ?", desde, hasta).
		Order("fecha ASC").Find(&movs).Error
	return movs, err
}

// helper para tests: simula gorm.ErrRecordNotFound con un error genérico
func ErrNotFound() error {
	return gorm.ErrRecordNotFound
}

// PagoRepo ─────────────────────────────────────────────────────────────────────

type PagoRepo interface {
	Create(p *models.Pago) error
	FindByCliente(clienteID int) ([]models.Pago, error)
	SumaPagadoPorCliente(clienteID int) (float64, error)
}

type pagoRepoDB struct{ db *gorm.DB }

func NewPagoRepo(db *gorm.DB) PagoRepo { return &pagoRepoDB{db: db} }

func (r *pagoRepoDB) Create(p *models.Pago) error {
	p.Fecha = time.Now()
	return r.db.Create(p).Error
}

func (r *pagoRepoDB) FindByCliente(clienteID int) ([]models.Pago, error) {
	var pagos []models.Pago
	err := r.db.Where("cliente_id = ?", clienteID).Order("fecha DESC").Find(&pagos).Error
	return pagos, err
}

func (r *pagoRepoDB) SumaPagadoPorCliente(clienteID int) (float64, error) {
	var suma float64
	err := r.db.Model(&models.Pago{}).
		Where("cliente_id = ?", clienteID).
		Select("COALESCE(SUM(monto), 0)").Scan(&suma).Error
	return suma, err
}
