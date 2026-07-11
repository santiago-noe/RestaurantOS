package services

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

// ─── Mocks ───────────────────────────────────────────────────────────────────

type mockPedidoRepo struct{ mock.Mock }

func (m *mockPedidoRepo) Create(p *models.Pedido, items []models.PedidoItem) error {
	args := m.Called(p, items)
	if args.Error(0) == nil {
		p.ID = 1 // simula autoincrement
	}
	return args.Error(0)
}
func (m *mockPedidoRepo) FindByID(id int) (*models.Pedido, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Pedido), args.Error(1)
}
func (m *mockPedidoRepo) FindAll(page, perPage, clienteID int, estado string) ([]models.Pedido, int64, error) {
	args := m.Called(page, perPage, clienteID, estado)
	return args.Get(0).([]models.Pedido), args.Get(1).(int64), args.Error(2)
}
func (m *mockPedidoRepo) UpdateEstado(id int, estado string) error {
	return m.Called(id, estado).Error(0)
}
func (m *mockPedidoRepo) FindEntreFechas(desde, hasta time.Time) ([]models.Pedido, error) {
	args := m.Called(desde, hasta)
	return args.Get(0).([]models.Pedido), args.Error(1)
}

type mockProductoRepo struct{ mock.Mock }

func (m *mockProductoRepo) FindByID(id int) (*models.Producto, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Producto), args.Error(1)
}
func (m *mockProductoRepo) AjustarStock(id int, cantidad float64) error {
	return m.Called(id, cantidad).Error(0)
}
func (m *mockProductoRepo) FindAll(soloActivos bool) ([]models.Producto, error) {
	args := m.Called(soloActivos)
	return args.Get(0).([]models.Producto), args.Error(1)
}
func (m *mockProductoRepo) Create(p *models.Producto) error { return m.Called(p).Error(0) }
func (m *mockProductoRepo) Update(id int, fields map[string]interface{}) (*models.Producto, error) {
	args := m.Called(id, fields)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Producto), args.Error(1)
}

type mockClienteRepo struct{ mock.Mock }

func (m *mockClienteRepo) Create(c *models.Cliente) error      { return m.Called(c).Error(0) }
func (m *mockClienteRepo) FindByID(id int) (*models.Cliente, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cliente), args.Error(1)
}
func (m *mockClienteRepo) FindAll(p, pp int, t string) ([]models.Cliente, int64, error) {
	args := m.Called(p, pp, t)
	return args.Get(0).([]models.Cliente), args.Get(1).(int64), args.Error(2)
}
func (m *mockClienteRepo) Update(id int, f map[string]interface{}) (*models.Cliente, error) {
	args := m.Called(id, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cliente), args.Error(1)
}
func (m *mockClienteRepo) Deactivate(id int) error          { return m.Called(id).Error(0) }
func (m *mockClienteRepo) EmailExists(e string, id int) bool { return m.Called(e, id).Bool(0) }

// Verificar interfaces en tiempo de compilación
var _ repository.PedidoRepo   = (*mockPedidoRepo)(nil)
var _ repository.ProductoRepo  = (*mockProductoRepo)(nil)
var _ repository.ClienteRepo   = (*mockClienteRepo)(nil)

// ─── CalcularTotal (función pura) ─────────────────────────────────────────────

func TestCalcularTotal_ItemsMultiplesRetornaSumaCorrecta(t *testing.T) {
	items := []ItemInput{
		{Cantidad: 2, PrecioUnitario: 12.50},
		{Cantidad: 1, PrecioUnitario: 3.00},
	}
	assert.Equal(t, 28.00, CalcularTotal(items))
}

func TestCalcularTotal_ListaVaciaRetornaCero(t *testing.T) {
	assert.Equal(t, 0.00, CalcularTotal([]ItemInput{}))
}

func TestCalcularTotal_UnItemRetornaSubtotalCorrecto(t *testing.T) {
	items := []ItemInput{{Cantidad: 3, PrecioUnitario: 12.00}}
	assert.Equal(t, 36.00, CalcularTotal(items))
}

func TestCalcularTotal_PrecisionDecimalCorrecta(t *testing.T) {
	items := []ItemInput{{Cantidad: 3, PrecioUnitario: 0.10}}
	assert.Equal(t, 0.30, CalcularTotal(items))
}

// ─── PedidoService.Crear ─────────────────────────────────────────────────────

func newService(pr repository.PedidoRepo, prod repository.ProductoRepo, cl repository.ClienteRepo) *PedidoService {
	return NewPedidoService(pr, prod, cl)
}

func inputValido() CrearPedidoInput {
	return CrearPedidoInput{
		ClienteID:  1,
		UserID:     1,
		Fecha:      time.Now(),
		TipoComida: "almuerzo",
		FormaPago:  "contado",
		Items: []ItemInput{
			{ProductoID: 10, Cantidad: 2, PrecioUnitario: 12.50},
		},
	}
}

func TestCrearPedido_DatosValidosCreaPedido(t *testing.T) {
	pr := new(mockPedidoRepo)
	prod := new(mockProductoRepo)
	cl := new(mockClienteRepo)

	cl.On("FindByID", 1).Return(&models.Cliente{ID: 1, Activo: true}, nil)
	prod.On("FindByID", 10).Return(&models.Producto{ID: 10, StockActual: 5, Activo: true}, nil)
	pr.On("Create", mock.AnythingOfType("*models.Pedido"), mock.Anything).Return(nil)
	prod.On("AjustarStock", 10, -2.0).Return(nil)

	pedido, err := newService(pr, prod, cl).Crear(inputValido())

	require.NoError(t, err)
	assert.Equal(t, 25.00, pedido.Total)
	assert.Equal(t, "pendiente", pedido.Estado)
	pr.AssertExpectations(t)
	prod.AssertExpectations(t)
}

func TestCrearPedido_ClienteInexistenteRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	cl.On("FindByID", 1).Return(nil, nil)

	_, err := newService(pr, prod, cl).Crear(inputValido())

	assert.ErrorContains(t, err, "cliente")
	pr.AssertNotCalled(t, "Create")
}

func TestCrearPedido_ProductoInexistenteRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	cl.On("FindByID", 1).Return(&models.Cliente{ID: 1, Activo: true}, nil)
	prod.On("FindByID", 10).Return(nil, nil)

	_, err := newService(pr, prod, cl).Crear(inputValido())

	assert.ErrorContains(t, err, "producto")
}

func TestCrearPedido_StockInsuficienteRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	cl.On("FindByID", 1).Return(&models.Cliente{ID: 1, Activo: true}, nil)
	// Stock actual = 1, se piden 2
	prod.On("FindByID", 10).Return(&models.Producto{ID: 10, StockActual: 1, Activo: true}, nil)

	_, err := newService(pr, prod, cl).Crear(inputValido())

	assert.ErrorContains(t, err, "stock")
}

func TestCrearPedido_ItemsVaciosRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	cl.On("FindByID", 1).Return(&models.Cliente{ID: 1, Activo: true}, nil)

	input := inputValido()
	input.Items = []ItemInput{}

	_, err := newService(pr, prod, cl).Crear(input)

	assert.ErrorContains(t, err, "items")
}

func TestCrearPedido_DescuentaStockCorrectamente(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	cl.On("FindByID", 1).Return(&models.Cliente{ID: 1, Activo: true}, nil)
	prod.On("FindByID", 10).Return(&models.Producto{ID: 10, StockActual: 10, Activo: true}, nil)
	pr.On("Create", mock.Anything, mock.Anything).Return(nil)
	// Verifica que descuenta exactamente -2
	prod.On("AjustarStock", 10, -2.0).Return(nil)

	_, err := newService(pr, prod, cl).Crear(inputValido())

	require.NoError(t, err)
	prod.AssertCalled(t, "AjustarStock", 10, -2.0)
}

func TestCrearPedido_ErrorAlGuardarRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	cl.On("FindByID", 1).Return(&models.Cliente{ID: 1, Activo: true}, nil)
	prod.On("FindByID", 10).Return(&models.Producto{ID: 10, StockActual: 10, Activo: true}, nil)
	pr.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	_, err := newService(pr, prod, cl).Crear(inputValido())

	assert.Error(t, err)
}

// ─── PedidoService.MarcarEntregado ────────────────────────────────────────────

func TestMarcarEntregado_PendienteCambiaAEntregado(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 5).Return(&models.Pedido{ID: 5, Estado: "pendiente"}, nil)
	pr.On("UpdateEstado", 5, "entregado").Return(nil)

	err := newService(pr, prod, cl).MarcarEntregado(5)

	require.NoError(t, err)
	pr.AssertCalled(t, "UpdateEstado", 5, "entregado")
}

func TestMarcarEntregado_YaEntregadoRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 5).Return(&models.Pedido{ID: 5, Estado: "entregado"}, nil)

	err := newService(pr, prod, cl).MarcarEntregado(5)

	assert.ErrorContains(t, err, "entregado")
	pr.AssertNotCalled(t, "UpdateEstado")
}

func TestMarcarEntregado_AnuladoRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 5).Return(&models.Pedido{ID: 5, Estado: "anulado"}, nil)

	err := newService(pr, prod, cl).MarcarEntregado(5)

	assert.ErrorContains(t, err, "anulado")
	pr.AssertNotCalled(t, "UpdateEstado")
}

func TestMarcarEntregado_InexistenteRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 99).Return(nil, nil)

	err := newService(pr, prod, cl).MarcarEntregado(99)

	assert.ErrorContains(t, err, "pedido")
}

func TestMarcarEntregado_ErrorAlActualizarRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 5).Return(&models.Pedido{ID: 5, Estado: "pendiente"}, nil)
	pr.On("UpdateEstado", 5, "entregado").Return(errors.New("db error"))

	err := newService(pr, prod, cl).MarcarEntregado(5)

	assert.Error(t, err)
}

// ─── PedidoService.Anular ────────────────────────────────────────────────────

func TestAnularPedido_EntregadoDevuelveStockYCambiaEstado(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pedido := &models.Pedido{
		ID:     5,
		Estado: "entregado",
		Items:  []models.PedidoItem{{ProductoID: 10, Cantidad: 2}},
	}
	pr.On("FindByID", 5).Return(pedido, nil)
	pr.On("UpdateEstado", 5, "anulado").Return(nil)
	prod.On("AjustarStock", 10, 2.0).Return(nil) // devuelve stock

	err := newService(pr, prod, cl).Anular(5)

	require.NoError(t, err)
	pr.AssertCalled(t, "UpdateEstado", 5, "anulado")
	prod.AssertCalled(t, "AjustarStock", 10, 2.0)
}

func TestAnularPedido_YaAnuladoRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 5).Return(&models.Pedido{ID: 5, Estado: "anulado"}, nil)

	err := newService(pr, prod, cl).Anular(5)

	assert.ErrorContains(t, err, "anulado")
	pr.AssertNotCalled(t, "UpdateEstado")
}

func TestAnularPedido_InexistenteRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 99).Return(nil, nil)

	err := newService(pr, prod, cl).Anular(99)

	assert.ErrorContains(t, err, "pedido")
}

func TestAnularPedido_ErrorAlActualizarEstadoRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 5).Return(&models.Pedido{ID: 5, Estado: "pendiente", Items: []models.PedidoItem{}}, nil)
	pr.On("UpdateEstado", 5, "anulado").Return(errors.New("db error"))

	err := newService(pr, prod, cl).Anular(5)

	assert.Error(t, err)
}

func TestAnularPedido_ErrorAlDevolverStockRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 5).Return(&models.Pedido{
		ID:     5,
		Estado: "entregado",
		Items:  []models.PedidoItem{{ProductoID: 10, Cantidad: 2}},
	}, nil)
	pr.On("UpdateEstado", 5, "anulado").Return(nil)
	prod.On("AjustarStock", 10, 2.0).Return(errors.New("db error"))

	err := newService(pr, prod, cl).Anular(5)

	assert.Error(t, err)
}

func TestCrearPedido_ErrorAlDescontarStockRetornaError(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	cl.On("FindByID", 1).Return(&models.Cliente{ID: 1, Activo: true}, nil)
	prod.On("FindByID", 10).Return(&models.Producto{ID: 10, StockActual: 10, Activo: true}, nil)
	pr.On("Create", mock.Anything, mock.Anything).Return(nil)
	prod.On("AjustarStock", 10, -2.0).Return(errors.New("db error"))

	_, err := newService(pr, prod, cl).Crear(inputValido())

	assert.Error(t, err)
}

func TestFindByID_RetornaPedidoCorrecto(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindByID", 7).Return(&models.Pedido{ID: 7, Total: 50.00}, nil)

	pedido, err := newService(pr, prod, cl).FindByID(7)

	require.NoError(t, err)
	assert.Equal(t, 7, pedido.ID)
}

func TestFindAll_RetornaListaConTotal(t *testing.T) {
	pr, prod, cl := new(mockPedidoRepo), new(mockProductoRepo), new(mockClienteRepo)
	pr.On("FindAll", 1, 20, 0, "").Return([]models.Pedido{{ID: 1}, {ID: 2}}, int64(2), nil)

	pedidos, total, err := newService(pr, prod, cl).FindAll(1, 20, 0, "")

	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, pedidos, 2)
}
