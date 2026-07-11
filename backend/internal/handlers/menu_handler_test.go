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
	"gorm.io/gorm"

	"restaurantos/internal/middleware"
	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

// ─── Mock del repositorio ────────────────────────────────────────────────────

type mockMenuRepo struct {
	mock.Mock
}

func (m *mockMenuRepo) FindPublico() ([]models.MenuPublico, error) {
	args := m.Called()
	return args.Get(0).([]models.MenuPublico), args.Error(1)
}

func (m *mockMenuRepo) FindAll() ([]models.MenuPublico, error) {
	args := m.Called()
	return args.Get(0).([]models.MenuPublico), args.Error(1)
}

func (m *mockMenuRepo) FindByID(id int) (*models.MenuPublico, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MenuPublico), args.Error(1)
}

func (m *mockMenuRepo) Create(item *models.MenuPublico) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *mockMenuRepo) Update(id int, fields map[string]interface{}) (*models.MenuPublico, error) {
	args := m.Called(id, fields)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MenuPublico), args.Error(1)
}

func (m *mockMenuRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

var _ repository.MenuRepo = (*mockMenuRepo)(nil)

// ─── Helper para router de tests ─────────────────────────────────────────────

func routerMenu(repo repository.MenuRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewMenuHandler(repo)

	const secret = "test-secret"

	r.GET("/api/public/menu", h.Publico)

	admin := r.Group("/api/admin", middleware.JWTMiddleware(secret), middleware.RequireRole("admin"))
	admin.GET("/menu", h.Listar)
	admin.POST("/menu", h.Crear)
	admin.PUT("/menu/:id", h.Actualizar)
	admin.DELETE("/menu/:id", h.Eliminar)

	return r
}

// ─── GET /api/public/menu ─────────────────────────────────────────────────────

func TestMenuPublico_RetornaSoloDisponibles(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("FindPublico").Return([]models.MenuPublico{{ID: 1, Nombre: "Puca Picante", Disponible: true}}, nil)

	r := routerMenu(repo)
	req := httptest.NewRequest(http.MethodGet, "/api/public/menu", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestMenuPublico_NoRequiereAutenticacion(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("FindPublico").Return([]models.MenuPublico{}, nil)

	r := routerMenu(repo)
	req := httptest.NewRequest(http.MethodGet, "/api/public/menu", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMenuPublico_ErrorBDRetorna500(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("FindPublico").Return([]models.MenuPublico{}, errors.New("db error"))

	r := routerMenu(repo)
	req := httptest.NewRequest(http.MethodGet, "/api/public/menu", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── POST /api/admin/menu ─────────────────────────────────────────────────────

func TestCrearMenu_DatosValidosRetorna201(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("Create", mock.AnythingOfType("*models.MenuPublico")).Return(nil)

	r := routerMenu(repo)
	body, _ := json.Marshal(map[string]interface{}{"categoria": "Fondos", "nombre": "Cuy Chactado", "precio": 78.0})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/menu", bytes.NewReader(body))
	req.Header.Set("Authorization", tokenAdmin())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
}

func TestCrearMenu_SinNombreRetorna400(t *testing.T) {
	repo := new(mockMenuRepo)

	r := routerMenu(repo)
	body, _ := json.Marshal(map[string]interface{}{"categoria": "Fondos"})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/menu", bytes.NewReader(body))
	req.Header.Set("Authorization", tokenAdmin())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCrearMenu_SinDisponibleExplicitoQuedaDisponibleTrue(t *testing.T) {
	repo := new(mockMenuRepo)
	var itemCreado *models.MenuPublico
	repo.On("Create", mock.AnythingOfType("*models.MenuPublico")).
		Run(func(args mock.Arguments) { itemCreado = args.Get(0).(*models.MenuPublico) }).
		Return(nil)

	r := routerMenu(repo)
	body, _ := json.Marshal(map[string]interface{}{"categoria": "Fondos", "nombre": "Cuy Chactado"})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/menu", bytes.NewReader(body))
	req.Header.Set("Authorization", tokenAdmin())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.NotNil(t, itemCreado)
	assert.True(t, itemCreado.Disponible)
}

func TestCrearMenu_EmpleadoNoTienePermisoRetorna403(t *testing.T) {
	repo := new(mockMenuRepo)

	r := routerMenu(repo)
	body, _ := json.Marshal(map[string]interface{}{"categoria": "Fondos", "nombre": "Cuy Chactado"})
	req := httptest.NewRequest(http.MethodPost, "/api/admin/menu", bytes.NewReader(body))
	req.Header.Set("Authorization", tokenEmpleado())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── GET /api/admin/menu ──────────────────────────────────────────────────────

func TestListarMenuAdmin_RetornaTodosLosItems(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("FindAll").Return([]models.MenuPublico{
		{ID: 1, Nombre: "Disponible", Disponible: true},
		{ID: 2, Nombre: "Agotado", Disponible: false},
	}, nil)

	r := routerMenu(repo)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/menu", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── PUT /api/admin/menu/:id ───────────────────────────────────────────────────

func TestActualizarMenu_DatosValidosRetorna200(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("FindByID", 1).Return(&models.MenuPublico{ID: 1, Nombre: "Original"}, nil)
	repo.On("Update", 1, mock.Anything).Return(&models.MenuPublico{ID: 1, Nombre: "Editado"}, nil)

	r := routerMenu(repo)
	body, _ := json.Marshal(map[string]interface{}{"nombre": "Editado"})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/menu/1", bytes.NewReader(body))
	req.Header.Set("Authorization", tokenAdmin())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestActualizarMenu_InexistenteRetorna404(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("FindByID", 999).Return(nil, nil)

	r := routerMenu(repo)
	body, _ := json.Marshal(map[string]interface{}{"nombre": "Editado"})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/menu/999", bytes.NewReader(body))
	req.Header.Set("Authorization", tokenAdmin())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestActualizarMenu_PuedeMarcarComoNoDisponible(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("FindByID", 1).Return(&models.MenuPublico{ID: 1, Nombre: "Original", Disponible: true}, nil)
	var fieldsRecibidos map[string]interface{}
	repo.On("Update", 1, mock.Anything).
		Run(func(args mock.Arguments) { fieldsRecibidos = args.Get(1).(map[string]interface{}) }).
		Return(&models.MenuPublico{ID: 1, Nombre: "Original", Disponible: false}, nil)

	r := routerMenu(repo)
	body, _ := json.Marshal(map[string]interface{}{"disponible": false})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/menu/1", bytes.NewReader(body))
	req.Header.Set("Authorization", tokenAdmin())
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, false, fieldsRecibidos["disponible"])
}

// ─── DELETE /api/admin/menu/:id ────────────────────────────────────────────────

func TestEliminarMenu_ExisteRetorna200(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("Delete", 1).Return(nil)

	r := routerMenu(repo)
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/menu/1", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEliminarMenu_InexistenteRetorna404(t *testing.T) {
	repo := new(mockMenuRepo)
	repo.On("Delete", 999).Return(gorm.ErrRecordNotFound)

	r := routerMenu(repo)
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/menu/999", nil)
	req.Header.Set("Authorization", tokenAdmin())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
