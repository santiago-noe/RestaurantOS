package services

import (
	"errors"
	"math"
	"time"

	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

type ItemInput struct {
	ProductoID     int     `json:"producto_id"`
	Cantidad       float64 `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
}

type CrearPedidoInput struct {
	ClienteID  int
	UserID     int
	Fecha      time.Time
	TipoComida string
	FormaPago  string
	Notas      string
	Items      []ItemInput
}

// CalcularTotal es una función pura — sin efectos secundarios, fácil de testear.
func CalcularTotal(items []ItemInput) float64 {
	var total float64
	for _, item := range items {
		total += item.Cantidad * item.PrecioUnitario
	}
	return math.Round(total*100) / 100
}

// ─── PedidoService ────────────────────────────────────────────────────────────

type PedidoService struct {
	pedidoRepo  repository.PedidoRepo
	productoRepo repository.ProductoRepo
	clienteRepo  repository.ClienteRepo
}

func NewPedidoService(pr repository.PedidoRepo, prod repository.ProductoRepo, cl repository.ClienteRepo) *PedidoService {
	return &PedidoService{pedidoRepo: pr, productoRepo: prod, clienteRepo: cl}
}

func (s *PedidoService) Crear(input CrearPedidoInput) (*models.Pedido, error) {
	// 1. Validar que hay items
	if len(input.Items) == 0 {
		return nil, errors.New("el pedido debe tener al menos un items")
	}

	// 2. Validar cliente
	cliente, err := s.clienteRepo.FindByID(input.ClienteID)
	if err != nil {
		return nil, err
	}
	if cliente == nil {
		return nil, errors.New("cliente no encontrado")
	}

	// 3. Validar productos y stock
	for _, item := range input.Items {
		producto, err := s.productoRepo.FindByID(item.ProductoID)
		if err != nil {
			return nil, err
		}
		if producto == nil {
			return nil, errors.New("producto no encontrado")
		}
		if producto.StockActual < item.Cantidad {
			return nil, errors.New("stock insuficiente para el producto")
		}
	}

	// 4. Calcular total y construir entidades
	total := CalcularTotal(input.Items)
	pedido := &models.Pedido{
		ClienteID:  input.ClienteID,
		UserID:     input.UserID,
		Fecha:      input.Fecha,
		TipoComida: input.TipoComida,
		Estado:     "pendiente",
		FormaPago:  input.FormaPago,
		Total:      total,
		Notas:      input.Notas,
	}

	items := make([]models.PedidoItem, len(input.Items))
	for i, it := range input.Items {
		items[i] = models.PedidoItem{
			ProductoID:     it.ProductoID,
			Cantidad:       it.Cantidad,
			PrecioUnitario: it.PrecioUnitario,
			Subtotal:       math.Round(it.Cantidad*it.PrecioUnitario*100) / 100,
		}
	}

	// 5. Guardar pedido + items en transacción
	if err := s.pedidoRepo.Create(pedido, items); err != nil {
		return nil, err
	}

	// 6. Descontar stock
	for _, it := range input.Items {
		if err := s.productoRepo.AjustarStock(it.ProductoID, -it.Cantidad); err != nil {
			return nil, err
		}
	}

	pedido.Items = items
	return pedido, nil
}

func (s *PedidoService) MarcarEntregado(pedidoID int) error {
	pedido, err := s.pedidoRepo.FindByID(pedidoID)
	if err != nil {
		return err
	}
	if pedido == nil {
		return errors.New("pedido no encontrado")
	}
	if pedido.Estado == "anulado" {
		return errors.New("no se puede marcar como entregado un pedido anulado")
	}
	if pedido.Estado == "entregado" {
		return errors.New("el pedido ya está entregado")
	}

	return s.pedidoRepo.UpdateEstado(pedidoID, "entregado")
}

func (s *PedidoService) Anular(pedidoID int) error {
	pedido, err := s.pedidoRepo.FindByID(pedidoID)
	if err != nil {
		return err
	}
	if pedido == nil {
		return errors.New("pedido no encontrado")
	}
	if pedido.Estado == "anulado" {
		return errors.New("el pedido ya está anulado")
	}

	// Cambiar estado
	if err := s.pedidoRepo.UpdateEstado(pedidoID, "anulado"); err != nil {
		return err
	}

	// Devolver stock
	for _, item := range pedido.Items {
		if err := s.productoRepo.AjustarStock(item.ProductoID, item.Cantidad); err != nil {
			return err
		}
	}

	return nil
}

func (s *PedidoService) FindByID(id int) (*models.Pedido, error) {
	return s.pedidoRepo.FindByID(id)
}

func (s *PedidoService) FindAll(page, perPage, clienteID int, estado string) ([]models.Pedido, int64, error) {
	return s.pedidoRepo.FindAll(page, perPage, clienteID, estado)
}
