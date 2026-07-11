package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"restaurantos/internal/middleware"
	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

// ─── Mocks ────────────────────────────────────────────────────────────────────

type mockPedidoRepoR struct{ mock.Mock }

func (m *mockPedidoRepoR) Create(p *models.Pedido, items []models.PedidoItem) error {
	return m.Called(p, items).Error(0)
}
func (m *mockPedidoRepoR) FindByID(id int) (*models.Pedido, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Pedido), args.Error(1)
}
func (m *mockPedidoRepoR) FindAll(page, perPage, clienteID int, estado string) ([]models.Pedido, int64, error) {
	args := m.Called(page, perPage, clienteID, estado)
	return args.Get(0).([]models.Pedido), args.Get(1).(int64), args.Error(2)
}
func (m *mockPedidoRepoR) UpdateEstado(id int, estado string) error {
	return m.Called(id, estado).Error(0)
}
func (m *mockPedidoRepoR) FindEntreFechas(desde, hasta time.Time) ([]models.Pedido, error) {
	args := m.Called(desde, hasta)
	return args.Get(0).([]models.Pedido), args.Error(1)
}

type mockMovimientoRepoR struct{ mock.Mock }

func (m *mockMovimientoRepoR) Registrar(mv *models.MovimientoStock) error {
	return m.Called(mv).Error(0)
}
func (m *mockMovimientoRepoR) FindByProducto(id int) ([]models.MovimientoStock, error) {
	args := m.Called(id)
	return args.Get(0).([]models.MovimientoStock), args.Error(1)
}
func (m *mockMovimientoRepoR) FindEntreFechas(desde, hasta time.Time) ([]models.MovimientoStock, error) {
	args := m.Called(desde, hasta)
	return args.Get(0).([]models.MovimientoStock), args.Error(1)
}

type mockClienteRepoR struct{ mock.Mock }

func (m *mockClienteRepoR) Create(c *models.Cliente) error { return m.Called(c).Error(0) }
func (m *mockClienteRepoR) FindByID(id int) (*models.Cliente, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cliente), args.Error(1)
}
func (m *mockClienteRepoR) FindAll(page, perPage int, tipo string) ([]models.Cliente, int64, error) {
	args := m.Called(page, perPage, tipo)
	return args.Get(0).([]models.Cliente), args.Get(1).(int64), args.Error(2)
}
func (m *mockClienteRepoR) Update(id int, fields map[string]interface{}) (*models.Cliente, error) {
	args := m.Called(id, fields)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cliente), args.Error(1)
}
func (m *mockClienteRepoR) Deactivate(id int) error { return m.Called(id).Error(0) }
func (m *mockClienteRepoR) EmailExists(email string, excludeID int) bool {
	return m.Called(email, excludeID).Bool(0)
}

var _ repository.PedidoRepo = (*mockPedidoRepoR)(nil)
var _ repository.MovimientoRepo = (*mockMovimientoRepoR)(nil)
var _ repository.ClienteRepo = (*mockClienteRepoR)(nil)

func routerReportes(pr repository.PedidoRepo, mr repository.MovimientoRepo, cr repository.ClienteRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewReportesHandler(pr, mr, cr)

	admin := r.Group("/api/admin", middleware.JWTMiddleware("test-secret"), middleware.RequireRole("admin"))
	admin.GET("/reportes/ventas", h.Ventas)
	admin.GET("/reportes/deudores", h.Deudores)
	admin.GET("/reportes/inventario", h.Inventario)

	return r
}

// ─── Ventas ───────────────────────────────────────────────────────────────────

func TestReportesVentas_PeriodoInvalidoRetorna400(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/ventas?periodo=invalido", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReportesVentas_JSONPorDefectoRetorna200(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	pr.On("FindEntreFechas", mock.Anything, mock.Anything).Return([]models.Pedido{}, nil)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/ventas?periodo=diario", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestReportesVentas_FormatoExcelDevuelveArchivo(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	pr.On("FindEntreFechas", mock.Anything, mock.Anything).Return([]models.Pedido{}, nil)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/ventas?periodo=diario&formato=excel", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "spreadsheetml")
	assert.NotEmpty(t, w.Body.Bytes())
}

func TestReportesVentas_FormatoPDFDevuelveArchivo(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	pr.On("FindEntreFechas", mock.Anything, mock.Anything).Return([]models.Pedido{}, nil)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/ventas?periodo=diario&formato=pdf", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
}

func TestReportesVentas_ErrorBDRetorna500(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	pr.On("FindEntreFechas", mock.Anything, mock.Anything).Return([]models.Pedido{}, assert.AnError)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/ventas", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestReportesVentas_EmpleadoNoTienePermisoRetorna403(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/ventas", nil)
	req.Header.Set("Authorization", tokenEmpleado())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── Deudores ─────────────────────────────────────────────────────────────────

func TestReportesDeudores_JSONPorDefectoRetorna200(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	cr.On("FindAll", 1, 1000, "").Return([]models.Cliente{{Nombre: "Juan", DeudaTotal: 50}}, int64(1), nil)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/deudores", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReportesDeudores_FormatoPDFDevuelveArchivo(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	cr.On("FindAll", 1, 1000, "").Return([]models.Cliente{}, int64(0), nil)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/deudores?formato=pdf", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
}

// ─── Inventario ───────────────────────────────────────────────────────────────

func TestReportesInventario_JSONPorDefectoRetorna200(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	mr.On("FindEntreFechas", mock.Anything, mock.Anything).Return([]models.MovimientoStock{}, nil)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/inventario", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReportesInventario_FormatoExcelDevuelveArchivo(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	mr.On("FindEntreFechas", mock.Anything, mock.Anything).Return([]models.MovimientoStock{}, nil)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/inventario?formato=excel", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "spreadsheetml")
}

func TestReportesInventario_ErrorBDRetorna500(t *testing.T) {
	pr, mr, cr := new(mockPedidoRepoR), new(mockMovimientoRepoR), new(mockClienteRepoR)
	mr.On("FindEntreFechas", mock.Anything, mock.Anything).Return([]models.MovimientoStock{}, assert.AnError)
	r := routerReportes(pr, mr, cr)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reportes/inventario", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
