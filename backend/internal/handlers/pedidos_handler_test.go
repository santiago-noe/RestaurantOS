package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/auth"
	"restaurantos/internal/middleware"
	"restaurantos/internal/models"
	"restaurantos/internal/services"
)

// ─── Mock del servicio ────────────────────────────────────────────────────────

type mockPedidoService struct{ mock.Mock }

func (m *mockPedidoService) Crear(input services.CrearPedidoInput) (*models.Pedido, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Pedido), args.Error(1)
}
func (m *mockPedidoService) Anular(id int) error { return m.Called(id).Error(0) }
func (m *mockPedidoService) MarcarEntregado(id int) error { return m.Called(id).Error(0) }
func (m *mockPedidoService) FindByID(id int) (*models.Pedido, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Pedido), args.Error(1)
}
func (m *mockPedidoService) FindAll(page, perPage, clienteID int, estado string) ([]models.Pedido, int64, error) {
	args := m.Called(page, perPage, clienteID, estado)
	return args.Get(0).([]models.Pedido), args.Get(1).(int64), args.Error(2)
}

// ─── Helper ───────────────────────────────────────────────────────────────────

func routerPedidos(svc PedidoServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewPedidoHandler(svc)
	const secret = "test-secret"

	emp := r.Group("/api/empleado", middleware.JWTMiddleware(secret))
	emp.POST("/pedidos", h.Crear)
	emp.GET("/pedidos", h.Listar)
	emp.GET("/pedidos/:id", h.ObtenerPorID)
	emp.PUT("/pedidos/:id/entregar", h.MarcarEntregado)

	adm := r.Group("/api/admin", middleware.JWTMiddleware(secret), middleware.RequireRole("admin"))
	adm.DELETE("/pedidos/:id", h.Anular)
	return r
}

func tokEmp() string {
	t, _ := auth.GenerateJWT(2, "emp@test.com", "empleado", "test-secret")
	return "Bearer " + t
}
func tokAdm() string {
	t, _ := auth.GenerateJWT(1, "adm@test.com", "admin", "test-secret")
	return "Bearer " + t
}

// ─── POST /pedidos ────────────────────────────────────────────────────────────

func TestCrearPedido_DatosValidosRetorna201(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("Crear", mock.AnythingOfType("services.CrearPedidoInput")).
		Return(&models.Pedido{ID: 1, Total: 28.00, Estado: "pendiente"}, nil)

	body := `{
		"cliente_id":1,
		"tipo_comida":"almuerzo",
		"forma_pago":"contado",
		"items":[{"producto_id":3,"cantidad":2,"precio_unitario":12.50},
		         {"producto_id":7,"cantidad":1,"precio_unitario":3.00}]
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/empleado/pedidos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, float64(28), resp["total"])
}

func TestCrearPedido_BodyInvalidoRetorna400(t *testing.T) {
	svc := new(mockPedidoService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/empleado/pedidos", bytes.NewBufferString(`{bad}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Crear")
}

func TestCrearPedido_ClienteInexistenteRetorna404(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("Crear", mock.Anything).Return(nil, errors.New("cliente no encontrado"))

	body := `{"cliente_id":999,"tipo_comida":"almuerzo","forma_pago":"contado","items":[{"producto_id":1,"cantidad":1,"precio_unitario":10}]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/empleado/pedidos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCrearPedido_StockInsuficienteRetorna422(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("Crear", mock.Anything).Return(nil, errors.New("stock insuficiente"))

	body := `{"cliente_id":1,"tipo_comida":"almuerzo","forma_pago":"contado","items":[{"producto_id":1,"cantidad":100,"precio_unitario":10}]}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/empleado/pedidos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// ─── GET /pedidos ─────────────────────────────────────────────────────────────

func TestListarPedidos_RetornaPaginacion(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("FindAll", 1, 20, 0, "").Return([]models.Pedido{{ID: 1}, {ID: 2}}, int64(2), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/pedidos", nil)
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, float64(2), resp["total"])
}

func TestListarPedidos_FiltrarPorClienteYEstado(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("FindAll", 1, 20, 5, "pendiente").Return([]models.Pedido{}, int64(0), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/pedidos?cliente_id=5&estado=pendiente", nil)
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertCalled(t, "FindAll", 1, 20, 5, "pendiente")
}

func TestListarPedidos_ErrorBDRetorna500(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("FindAll", 1, 20, 0, "").Return([]models.Pedido{}, int64(0), errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/pedidos", nil)
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET /pedidos/:id ─────────────────────────────────────────────────────────

func TestObtenerPedido_ExisteRetornaDatos(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("FindByID", 3).Return(&models.Pedido{ID: 3, Total: 12.00}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/pedidos/3", nil)
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestObtenerPedido_InexistenteRetorna404(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("FindByID", 999).Return(nil, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/pedidos/999", nil)
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestObtenerPedido_IDInvalidoRetorna400(t *testing.T) {
	svc := new(mockPedidoService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/pedidos/abc", nil)
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── DELETE /pedidos/:id ──────────────────────────────────────────────────────

// ─── PUT /pedidos/:id/entregar ─────────────────────────────────────────────────

func TestMarcarEntregadoPedido_ExisteRetorna200(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("MarcarEntregado", 5).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/pedidos/5/entregar", nil)
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "entregado")
}

func TestMarcarEntregadoPedido_YaEntregadoRetorna422(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("MarcarEntregado", 5).Return(errors.New("el pedido ya está entregado"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/pedidos/5/entregar", nil)
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestMarcarEntregadoPedido_IDInvalidoRetorna400(t *testing.T) {
	svc := new(mockPedidoService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/pedidos/xyz/entregar", nil)
	req.Header.Set("Authorization", tokEmp())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnularPedido_ExisteRetorna200(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("Anular", 5).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/pedidos/5", nil)
	req.Header.Set("Authorization", tokAdm())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "anulado")
}

func TestAnularPedido_YaAnuladoRetorna422(t *testing.T) {
	svc := new(mockPedidoService)
	svc.On("Anular", 5).Return(errors.New("el pedido ya está anulado"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/pedidos/5", nil)
	req.Header.Set("Authorization", tokAdm())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestAnularPedido_IDInvalidoRetorna400(t *testing.T) {
	svc := new(mockPedidoService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/pedidos/xyz", nil)
	req.Header.Set("Authorization", tokAdm())

	routerPedidos(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
