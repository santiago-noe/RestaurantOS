package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/middleware"
	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

// ─── Mocks ────────────────────────────────────────────────────────────────────

type mockProductoRepoH struct{ mock.Mock }

func (m *mockProductoRepoH) FindByID(id int) (*models.Producto, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Producto), args.Error(1)
}
func (m *mockProductoRepoH) AjustarStock(id int, cantidad float64) error {
	return m.Called(id, cantidad).Error(0)
}
func (m *mockProductoRepoH) FindAll(soloActivos bool) ([]models.Producto, error) {
	args := m.Called(soloActivos)
	return args.Get(0).([]models.Producto), args.Error(1)
}
func (m *mockProductoRepoH) Create(p *models.Producto) error { return m.Called(p).Error(0) }
func (m *mockProductoRepoH) Update(id int, f map[string]interface{}) (*models.Producto, error) {
	args := m.Called(id, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Producto), args.Error(1)
}

type mockMovimientoRepoH struct{ mock.Mock }

func (m *mockMovimientoRepoH) Registrar(mv *models.MovimientoStock) error {
	return m.Called(mv).Error(0)
}
func (m *mockMovimientoRepoH) FindByProducto(id int) ([]models.MovimientoStock, error) {
	args := m.Called(id)
	return args.Get(0).([]models.MovimientoStock), args.Error(1)
}
func (m *mockMovimientoRepoH) FindEntreFechas(desde, hasta time.Time) ([]models.MovimientoStock, error) {
	args := m.Called(desde, hasta)
	return args.Get(0).([]models.MovimientoStock), args.Error(1)
}

var _ repository.ProductoRepo  = (*mockProductoRepoH)(nil)
var _ repository.MovimientoRepo = (*mockMovimientoRepoH)(nil)

// ─── Router helper ────────────────────────────────────────────────────────────

func routerInventario(pr repository.ProductoRepo, mr repository.MovimientoRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewInventarioHandler(pr, mr)
	const secret = "test-secret"

	emp := r.Group("/api/empleado", middleware.JWTMiddleware(secret))
	emp.GET("/productos", h.Listar)

	adm := r.Group("/api/admin", middleware.JWTMiddleware(secret), middleware.RequireRole("admin"))
	adm.GET("/productos/alertas", h.Alertas)
	adm.GET("/productos/:id", h.ObtenerPorID)
	adm.POST("/productos", h.Crear)
	adm.PUT("/productos/:id", h.Actualizar)
	adm.POST("/productos/:id/restock", h.Restock)
	return r
}

// ─── Tests Listar ─────────────────────────────────────────────────────────────

func TestListarProductos_RetornaLista(t *testing.T) {
	pr := new(mockProductoRepoH)
	mr := new(mockMovimientoRepoH)
	pr.On("FindAll", true).Return([]models.Producto{
		{ID: 1, Nombre: "Arroz", StockActual: 10},
		{ID: 2, Nombre: "Pollo", StockActual: 5},
	}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/productos", nil)
	req.Header.Set("Authorization", tokEmp())

	routerInventario(pr, mr).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []models.Producto
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Len(t, resp, 2)
}

func TestListarProductos_ErrorBDRetorna500(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindAll", true).Return([]models.Producto{}, errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/productos", nil)
	req.Header.Set("Authorization", tokEmp())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Tests Alertas ────────────────────────────────────────────────────────────

func TestAlertasStock_MuestraProductosBajoMinimo(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindAll", false).Return([]models.Producto{
		{ID: 1, Nombre: "Arroz", StockActual: 1.0, StockMinimo: 5.0},  // bajo
		{ID: 2, Nombre: "Pollo", StockActual: 8.0, StockMinimo: 3.0},  // ok
	}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/productos/alertas", nil)
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	alertas := resp["alertas"].([]interface{})
	assert.Len(t, alertas, 1)
}

func TestAlertasStock_SinProductosBajos(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindAll", false).Return([]models.Producto{
		{ID: 1, Nombre: "Arroz", StockActual: 10.0, StockMinimo: 5.0},
	}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/productos/alertas", nil)
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	alertas := resp["alertas"].([]interface{})
	assert.Len(t, alertas, 0)
}

// ─── Tests ObtenerPorID ───────────────────────────────────────────────────────

func TestObtenerProducto_ExisteRetornaDatosYMovimientos(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 3).Return(&models.Producto{ID: 3, Nombre: "Aceite"}, nil)
	mr.On("FindByProducto", 3).Return([]models.MovimientoStock{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/productos/3", nil)
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestObtenerProducto_InexistenteRetorna404(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 99).Return(nil, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/productos/99", nil)
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestObtenerProducto_IDInvalidoRetorna400(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/productos/abc", nil)
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Tests Crear ──────────────────────────────────────────────────────────────

func TestCrearProducto_DatosValidosRetorna201(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("Create", mock.AnythingOfType("*models.Producto")).Return(nil)

	body := `{"nombre":"Tomate","unidad":"kg","stock_minimo":2.0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCrearProducto_SinNombreRetorna400(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)

	body := `{"unidad":"kg"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCrearProducto_UnidadInvalidaRetorna400(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)

	body := `{"nombre":"X","unidad":"tonelada"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Tests Restock ────────────────────────────────────────────────────────────

func TestRestock_AumentaStockYRegistraMovimiento(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 2).Return(&models.Producto{ID: 2, StockActual: 3.0}, nil)
	pr.On("AjustarStock", 2, 10.5).Return(nil)
	mr.On("Registrar", mock.AnythingOfType("*models.MovimientoStock")).Return(nil)

	body := `{"cantidad":10.5,"notas":"Compra del mercado"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos/2/restock", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	pr.AssertCalled(t, "AjustarStock", 2, 10.5)
	mr.AssertCalled(t, "Registrar", mock.Anything)
}

func TestRestock_CantidadCeroRetorna400(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)

	body := `{"cantidad":0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos/2/restock", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRestock_ProductoInexistenteRetorna404(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 99).Return(nil, nil)

	body := `{"cantidad":5}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos/99/restock", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Tests Actualizar ─────────────────────────────────────────────────────────

func TestActualizarProducto_DatosValidosRetorna200(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 1).Return(&models.Producto{ID: 1}, nil)
	pr.On("Update", 1, mock.Anything).Return(&models.Producto{ID: 1, Nombre: "Actualizado"}, nil)

	body := `{"nombre":"Actualizado","stock_minimo":3.0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/productos/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestActualizarProducto_InexistenteRetorna404(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 99).Return(nil, nil)

	body := `{"nombre":"X"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/productos/99", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestActualizarProducto_IDInvalidoRetorna400(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/productos/abc", bytes.NewBufferString(`{"nombre":"X"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())
	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRestockIDInvalidoRetorna400(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos/abc/restock", bytes.NewBufferString(`{"cantidad":5}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())
	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCrearProducto_ErrorBDRetorna500(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("Create", mock.AnythingOfType("*models.Producto")).Return(errors.New("db error"))
	body := `{"nombre":"X","unidad":"kg"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())
	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRestockErrorAjusteRetorna500(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 2).Return(&models.Producto{ID: 2, StockActual: 3.0}, nil)
	pr.On("AjustarStock", 2, 5.0).Return(errors.New("db error"))
	body := `{"cantidad":5}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos/2/restock", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())
	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestActualizarProducto_ErrorBuscarRetorna500(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 1).Return(nil, errors.New("db error"))
	body := `{"nombre":"X"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/productos/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())
	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestActualizarProducto_ErrorActualizarRetorna500(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 1).Return(&models.Producto{ID: 1}, nil)
	pr.On("Update", 1, mock.Anything).Return(nil, errors.New("db error"))
	body := `{"nombre":"Nuevo"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/productos/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())
	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestObtenerProducto_ErrorBDRetorna500(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 1).Return(nil, errors.New("db error"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/productos/1", nil)
	req.Header.Set("Authorization", tokAdm())
	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRestockErrorBuscarProductoRetorna500(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindByID", 2).Return(nil, errors.New("db error"))
	body := `{"cantidad":5}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/productos/2/restock", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())
	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAlertasErrorBDRetorna500(t *testing.T) {
	pr, mr := new(mockProductoRepoH), new(mockMovimientoRepoH)
	pr.On("FindAll", false).Return([]models.Producto{}, errors.New("db error"))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/productos/alertas", nil)
	req.Header.Set("Authorization", tokAdm())
	routerInventario(pr, mr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
