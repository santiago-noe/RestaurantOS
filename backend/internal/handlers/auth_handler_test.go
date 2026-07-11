package handlers

import (
	"bytes"
	"encoding/json"
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

// ─── Mock UserRepo ────────────────────────────────────────────────────────────

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepo) FindByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

var _ repository.UserRepo = (*mockUserRepo)(nil)

// ─── Helper ───────────────────────────────────────────────────────────────────

func routerAuth(repo repository.UserRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(repo, "test-secret")
	r.POST("/api/auth/login", h.Login)
	r.GET("/api/auth/me", middleware.JWTMiddleware("test-secret"), h.Me)
	return r
}

func usuarioValido() *models.User {
	hash, _ := auth.HashPassword("password123")
	return &models.User{
		ID:       1,
		Nombre:   "Admin",
		Apellido: "Test",
		Email:    "admin@test.com",
		Password: hash,
		Rol:      "admin",
		Activo:   true,
	}
}

// ─── Tests Login ─────────────────────────────────────────────────────────────

func TestLogin_DatosValidosRetorna200ConToken(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByEmail", "admin@test.com").Return(usuarioValido(), nil)

	body := `{"email":"admin@test.com","password":"password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	routerAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotEmpty(t, resp["token"])
	assert.NotNil(t, resp["user"])
}

func TestLogin_BodyInvalidoRetorna400(t *testing.T) {
	repo := new(mockUserRepo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(`{malformed`))
	req.Header.Set("Content-Type", "application/json")

	routerAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "FindByEmail")
}

func TestLogin_SinEmailRetorna400(t *testing.T) {
	repo := new(mockUserRepo)

	body := `{"password":"password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	routerAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_UsuarioNoExisteRetorna401(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByEmail", "noexiste@test.com").Return(nil, nil)

	body := `{"email":"noexiste@test.com","password":"password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	routerAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_ContraseñaIncorrectaRetorna401(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByEmail", "admin@test.com").Return(usuarioValido(), nil)

	body := `{"email":"admin@test.com","password":"incorrecta"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	routerAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── Tests Me ────────────────────────────────────────────────────────────────

func TestMe_TokenValidoRetornaDatosUsuario(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByID", 1).Return(usuarioValido(), nil)

	token, _ := auth.GenerateJWT(1, "admin@test.com", "admin", "test-secret")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	routerAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "Admin", resp["nombre"])
}

func TestMe_UsuarioEliminadoRetorna404(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByID", 99).Return(nil, nil)

	token, _ := auth.GenerateJWT(99, "ghost@test.com", "admin", "test-secret")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	routerAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMe_SinTokenRetorna401(t *testing.T) {
	repo := new(mockUserRepo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/auth/me", nil)

	routerAuth(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
