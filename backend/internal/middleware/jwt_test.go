package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"restaurantos/internal/auth"
)

func setupRouter(secret string, requireAdmin bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", JWTMiddleware(secret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.GET("/admin", JWTMiddleware(secret), RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func TestJWTMiddleware_SinHeaderDevuelve401(t *testing.T) {
	r := setupRouter("secret", false)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddleware_TokenInvalidoDevuelve401(t *testing.T) {
	r := setupRouter("secret", false)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token.invalido.aqui")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddleware_TokenValidoPasaAlSiguiente(t *testing.T) {
	secret := "mi-secreto"
	token, _ := auth.GenerateJWT(1, "user@test.com", "empleado", secret)

	r := setupRouter(secret, false)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTMiddleware_RolAdminAccedeARutaAdmin(t *testing.T) {
	secret := "mi-secreto"
	token, _ := auth.GenerateJWT(1, "admin@test.com", "admin", secret)

	r := setupRouter(secret, true)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTMiddleware_RolEmpleadoNoAccedeARutaAdmin(t *testing.T) {
	secret := "mi-secreto"
	token, _ := auth.GenerateJWT(2, "empleado@test.com", "empleado", secret)

	r := setupRouter(secret, true)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetClaims_RetornaClaimsDelContexto(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "mi-secreto"
	var capturedClaims *auth.Claims

	r := gin.New()
	r.GET("/me", JWTMiddleware(secret), func(c *gin.Context) {
		capturedClaims = GetClaims(c)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token, _ := auth.GenerateJWT(7, "user@test.com", "admin", secret)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedClaims)
	assert.Equal(t, 7, capturedClaims.UserID)
	assert.Equal(t, "admin", capturedClaims.Rol)
}

func TestGetClaims_SinClaimsRetornaNil(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var capturedClaims *auth.Claims

	r := gin.New()
	r.GET("/noop", func(c *gin.Context) {
		capturedClaims = GetClaims(c)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/noop", nil)
	r.ServeHTTP(w, req)

	assert.Nil(t, capturedClaims)
}

func TestRequireRole_SinAutenticacionDevuelve401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// RequireRole sin JWTMiddleware previo → no hay claims en el contexto
	r.GET("/solo-admin", RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/solo-admin", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
