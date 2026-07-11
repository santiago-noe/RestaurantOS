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
	"restaurantos/internal/repository"
)

// ─── Mock del repositorio ────────────────────────────────────────────────────

type mockClienteRepo struct {
	mock.Mock
}

func (m *mockClienteRepo) Create(c *models.Cliente) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *mockClienteRepo) FindByID(id int) (*models.Cliente, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cliente), args.Error(1)
}

func (m *mockClienteRepo) FindAll(page, perPage int, tipo string) ([]models.Cliente, int64, error) {
	args := m.Called(page, perPage, tipo)
	return args.Get(0).([]models.Cliente), args.Get(1).(int64), args.Error(2)
}

func (m *mockClienteRepo) Update(id int, fields map[string]interface{}) (*models.Cliente, error) {
	args := m.Called(id, fields)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cliente), args.Error(1)
}

func (m *mockClienteRepo) Deactivate(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockClienteRepo) EmailExists(email string, excludeID int) bool {
	args := m.Called(email, excludeID)
	return args.Bool(0)
}

// Verificar en tiempo de compilación que el mock implementa la interfaz
var _ repository.ClienteRepo = (*mockClienteRepo)(nil)

// ─── Helper para router de tests ─────────────────────────────────────────────

func routerConAuth(repo repository.ClienteRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewClienteHandler(repo)

	const secret = "test-secret"
	empleado := r.Group("/api/empleado", middleware.JWTMiddleware(secret))
	empleado.GET("/clientes", h.Listar)
	empleado.GET("/clientes/:id", h.ObtenerPorID)

	admin := r.Group("/api/admin", middleware.JWTMiddleware(secret), middleware.RequireRole("admin"))
	admin.POST("/clientes", h.Crear)
	admin.PUT("/clientes/:id", h.Actualizar)
	admin.DELETE("/clientes/:id", h.Desactivar)

	return r
}

func tokenAdmin() string {
	t, _ := auth.GenerateJWT(1, "admin@test.com", "admin", "test-secret")
	return "Bearer " + t
}

func tokenEmpleado() string {
	t, _ := auth.GenerateJWT(2, "emp@test.com", "empleado", "test-secret")
	return "Bearer " + t
}

// ─── POST /api/admin/clientes ─────────────────────────────────────────────────

func TestCrearCliente_DatosValidosRetorna201(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("EmailExists", "juan@test.com", 0).Return(false)
	repo.On("Create", mock.AnythingOfType("*models.Cliente")).Return(nil)

	body := `{"nombre":"Juan","tipo":"individual","email":"juan@test.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/clientes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	repo.AssertExpectations(t)
}

func TestCrearCliente_SinNombreRetorna400(t *testing.T) {
	repo := new(mockClienteRepo)

	body := `{"tipo":"individual"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/clientes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "Create")
}

func TestCrearCliente_TipoInvalidoRetorna400(t *testing.T) {
	repo := new(mockClienteRepo)

	body := `{"nombre":"Juan","tipo":"vip"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/clientes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "Create")
}

func TestCrearCliente_EmailDuplicadoRetorna409(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("EmailExists", "duplicado@test.com", 0).Return(true)

	body := `{"nombre":"Juan","tipo":"individual","email":"duplicado@test.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/clientes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	repo.AssertNotCalled(t, "Create")
}

func TestCrearCliente_EmpleadoNoTienePermisoRetorna403(t *testing.T) {
	repo := new(mockClienteRepo)

	body := `{"nombre":"Juan","tipo":"individual"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/clientes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenEmpleado())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── GET /api/empleado/clientes/:id ──────────────────────────────────────────

func TestObtenerCliente_IDExisteRetornaDatos(t *testing.T) {
	repo := new(mockClienteRepo)
	cliente := &models.Cliente{ID: 5, Nombre: "Pedro", Tipo: "individual"}
	repo.On("FindByID", 5).Return(cliente, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/clientes/5", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Cliente
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "Pedro", resp.Nombre)
}

func TestObtenerCliente_IDInexistenteRetorna404(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("FindByID", 999).Return(nil, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/clientes/999", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestObtenerCliente_IDNoNumericoRetorna400(t *testing.T) {
	repo := new(mockClienteRepo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/clientes/abc", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── GET /api/empleado/clientes ───────────────────────────────────────────────

func TestListarClientes_RetornaPaginacion(t *testing.T) {
	repo := new(mockClienteRepo)
	clientes := []models.Cliente{
		{ID: 1, Nombre: "Ana", Tipo: "individual"},
		{ID: 2, Nombre: "Constructora Norte", Tipo: "empresa"},
	}
	repo.On("FindAll", 1, 20, "").Return(clientes, int64(2), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/clientes", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, float64(2), resp["total"])
	assert.Equal(t, float64(1), resp["page"])
}

func TestListarClientes_FiltrarPorTipo(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("FindAll", 1, 20, "empresa").Return([]models.Cliente{}, int64(0), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/clientes?tipo=empresa", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertCalled(t, "FindAll", 1, 20, "empresa")
}

// ─── PUT /api/admin/clientes/:id ─────────────────────────────────────────────

func TestEditarCliente_DatosValidosActualiza(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("FindByID", 3).Return(&models.Cliente{ID: 3}, nil)
	repo.On("EmailExists", "nuevo@test.com", 3).Return(false)
	repo.On("Update", 3, mock.Anything).Return(&models.Cliente{ID: 3, Nombre: "Nuevo"}, nil)

	body := `{"nombre":"Nuevo","email":"nuevo@test.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/clientes/3", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEditarCliente_IDInexistenteRetorna404(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("FindByID", 999).Return(nil, nil)

	body := `{"nombre":"X"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/clientes/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── DELETE /api/admin/clientes/:id ──────────────────────────────────────────

func TestDesactivarCliente_SoftDeleteNoElimina(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("Deactivate", 4).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/clientes/4", nil)
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Verificar que dice "desactivado", no "eliminado"
	assert.Contains(t, w.Body.String(), "desactivado")
}

func TestDesactivarCliente_IDInexistenteRetorna404(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("Deactivate", 999).Return(gormErrNotFound())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/clientes/999", nil)
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── Paths de error (cobertura de ramas internas) ────────────────────────────

func TestCrearCliente_ErrorBDRetorna500(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("EmailExists", "", 0).Return(false)
	repo.On("Create", mock.AnythingOfType("*models.Cliente")).Return(errors.New("db error"))

	body := `{"nombre":"Juan","tipo":"individual"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/clientes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestObtenerCliente_ErrorBDRetorna500(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("FindByID", 1).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/clientes/1", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestListarClientes_ErrorBDRetorna500(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("FindAll", 1, 20, "").Return([]models.Cliente{}, int64(0), errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/empleado/clientes", nil)
	req.Header.Set("Authorization", tokenEmpleado())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEditarCliente_ErrorBuscarRetorna500(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("FindByID", 5).Return(nil, errors.New("db error"))

	body := `{"nombre":"X"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/clientes/5", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEditarCliente_EmailDuplicadoRetorna409(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("FindByID", 3).Return(&models.Cliente{ID: 3}, nil)
	repo.On("EmailExists", "usado@test.com", 3).Return(true)

	body := `{"email":"usado@test.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/clientes/3", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestEditarCliente_ErrorActualizarRetorna500(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("FindByID", 3).Return(&models.Cliente{ID: 3}, nil)
	repo.On("EmailExists", "", 3).Return(false)
	repo.On("Update", 3, mock.Anything).Return(nil, errors.New("db error"))

	body := `{"nombre":"Nuevo"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/clientes/3", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDesactivarCliente_ErrorGenericoRetorna500(t *testing.T) {
	repo := new(mockClienteRepo)
	repo.On("Deactivate", 4).Return(errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/clientes/4", nil)
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDesactivarCliente_IDInvalidoRetorna400(t *testing.T) {
	repo := new(mockClienteRepo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/clientes/xyz", nil)
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEditarCliente_IDInvalidoRetorna400(t *testing.T) {
	repo := new(mockClienteRepo)

	body := `{"nombre":"X"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/admin/clientes/abc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenAdmin())

	routerConAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// helper
func gormErrNotFound() error {
	// Simula gorm.ErrRecordNotFound
	return &notFoundError{}
}

type notFoundError struct{}

func (e *notFoundError) Error() string { return "record not found" }
