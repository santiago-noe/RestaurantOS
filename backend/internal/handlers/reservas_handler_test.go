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

	"restaurantos/internal/middleware"
	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

// ─── Mock del repositorio ────────────────────────────────────────────────────

type mockReservaRepo struct {
	mock.Mock
}

func (m *mockReservaRepo) Create(r *models.Reserva) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *mockReservaRepo) FindByID(id int) (*models.Reserva, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reserva), args.Error(1)
}

func (m *mockReservaRepo) FindAll(page, perPage int, estado string) ([]models.Reserva, int64, error) {
	args := m.Called(page, perPage, estado)
	return args.Get(0).([]models.Reserva), args.Get(1).(int64), args.Error(2)
}

func (m *mockReservaRepo) UpdateEstado(id int, estado string) error {
	args := m.Called(id, estado)
	return args.Error(0)
}

func (m *mockReservaRepo) VincularPedido(id int, pedidoID int) error {
	args := m.Called(id, pedidoID)
	return args.Error(0)
}

var _ repository.ReservaRepo = (*mockReservaRepo)(nil)

// ─── Helper para router de tests ─────────────────────────────────────────────

func routerReservas(repo repository.ReservaRepo, pedidoRepo repository.PedidoRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewReservaHandler(repo, pedidoRepo)

	const secret = "test-secret"

	r.POST("/api/public/reservas", h.Crear)

	empleado := r.Group("/api/empleado", middleware.JWTMiddleware(secret))
	empleado.GET("/reservas", h.Listar)
	empleado.PUT("/reservas/:id/estado", h.ActualizarEstado)
	empleado.PUT("/reservas/:id/pedido", h.VincularPedido)

	return r
}

// ─── POST /api/public/reservas ────────────────────────────────────────────────

func TestCrearReserva_DatosValidosRetorna201(t *testing.T) {
	repo := new(mockReservaRepo)
	repo.On("Create", mock.AnythingOfType("*models.Reserva")).Return(nil)

	body := `{"nombre":"Juan","whatsapp":"987654321","fecha":"2026-08-01","personas":"4"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/public/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	repo.AssertExpectations(t)
}

func TestCrearReserva_SinNombreRetorna400(t *testing.T) {
	repo := new(mockReservaRepo)

	body := `{"whatsapp":"987654321","fecha":"2026-08-01","personas":"4"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/public/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "Create")
}

func TestCrearReserva_FechaInvalidaRetorna400(t *testing.T) {
	repo := new(mockReservaRepo)

	body := `{"nombre":"Juan","whatsapp":"987654321","fecha":"01-08-2026","personas":"4"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/public/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "Create")
}

func TestCrearReserva_ErrorBDRetorna500(t *testing.T) {
	repo := new(mockReservaRepo)
	repo.On("Create", mock.AnythingOfType("*models.Reserva")).Return(errors.New("db error"))

	body := `{"nombre":"Juan","whatsapp":"987654321","fecha":"2026-08-01","personas":"4"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/public/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GET /api/empleado/reservas ───────────────────────────────────────────────

func TestListarReservas_RetornaPaginacion(t *testing.T) {
	repo := new(mockReservaRepo)
	reservas := []models.Reserva{
		{ID: 1, Nombre: "Ana", Estado: "pendiente"},
		{ID: 2, Nombre: "Luis", Estado: "confirmada"},
	}
	repo.On("FindAll", 1, 20, "").Return(reservas, int64(2), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/reservas", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, float64(2), resp["total"])
}

func TestListarReservas_FiltrarPorEstado(t *testing.T) {
	repo := new(mockReservaRepo)
	repo.On("FindAll", 1, 20, "pendiente").Return([]models.Reserva{}, int64(0), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/reservas?estado=pendiente", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertCalled(t, "FindAll", 1, 20, "pendiente")
}

func TestListarReservas_ErrorBDRetorna500(t *testing.T) {
	repo := new(mockReservaRepo)
	repo.On("FindAll", 1, 20, "").Return([]models.Reserva{}, int64(0), errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/reservas", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── PUT /api/empleado/reservas/:id/estado ───────────────────────────────────

func TestActualizarEstadoReserva_EstadoValidoRetorna200(t *testing.T) {
	repo := new(mockReservaRepo)
	repo.On("UpdateEstado", 3, "confirmada").Return(nil)

	body := `{"estado":"confirmada"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/3/estado", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestActualizarEstadoReserva_EstadoInvalidoRetorna400(t *testing.T) {
	repo := new(mockReservaRepo)

	body := `{"estado":"vip"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/3/estado", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "UpdateEstado")
}

func TestActualizarEstadoReserva_IDInexistenteRetorna404(t *testing.T) {
	repo := new(mockReservaRepo)
	repo.On("UpdateEstado", 999, "cancelada").Return(errors.New("record not found"))

	body := `{"estado":"cancelada"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/999/estado", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestActualizarEstadoReserva_IDInvalidoRetorna400(t *testing.T) {
	repo := new(mockReservaRepo)

	body := `{"estado":"confirmada"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/abc/estado", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, new(mockPedidoRepoR)).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── PUT /api/empleado/reservas/:id/pedido ───────────────────────────────────

func TestVincularPedidoReserva_PedidoExisteRetorna200(t *testing.T) {
	repo := new(mockReservaRepo)
	pedidoRepo := new(mockPedidoRepoR)
	pedidoRepo.On("FindByID", 7).Return(&models.Pedido{ID: 7}, nil)
	repo.On("VincularPedido", 3, 7).Return(nil)

	body := `{"pedido_id":7}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/3/pedido", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, pedidoRepo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVincularPedidoReserva_PedidoInexistenteRetorna404(t *testing.T) {
	repo := new(mockReservaRepo)
	pedidoRepo := new(mockPedidoRepoR)
	pedidoRepo.On("FindByID", 999).Return(nil, nil)

	body := `{"pedido_id":999}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/3/pedido", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, pedidoRepo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertNotCalled(t, "VincularPedido")
}

func TestVincularPedidoReserva_ReservaInexistenteRetorna404(t *testing.T) {
	repo := new(mockReservaRepo)
	pedidoRepo := new(mockPedidoRepoR)
	pedidoRepo.On("FindByID", 7).Return(&models.Pedido{ID: 7}, nil)
	repo.On("VincularPedido", 999, 7).Return(errors.New("record not found"))

	body := `{"pedido_id":7}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/999/pedido", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, pedidoRepo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVincularPedidoReserva_SinPedidoIDRetorna400(t *testing.T) {
	repo := new(mockReservaRepo)
	pedidoRepo := new(mockPedidoRepoR)

	body := `{}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/3/pedido", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, pedidoRepo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	pedidoRepo.AssertNotCalled(t, "FindByID")
}

func TestVincularPedidoReserva_IDInvalidoRetorna400(t *testing.T) {
	repo := new(mockReservaRepo)
	pedidoRepo := new(mockPedidoRepoR)

	body := `{"pedido_id":7}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/abc/pedido", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, pedidoRepo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVincularPedidoReserva_ErrorBuscarPedidoRetorna500(t *testing.T) {
	repo := new(mockReservaRepo)
	pedidoRepo := new(mockPedidoRepoR)
	pedidoRepo.On("FindByID", 7).Return(nil, errors.New("db error"))

	body := `{"pedido_id":7}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/empleado/reservas/3/pedido", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerReservas(repo, pedidoRepo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
