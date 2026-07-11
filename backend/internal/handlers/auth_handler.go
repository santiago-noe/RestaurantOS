package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"restaurantos/internal/auth"
	"restaurantos/internal/middleware"
	"restaurantos/internal/repository"
)

type AuthHandler struct {
	repo      repository.UserRepo
	jwtSecret string
}

func NewAuthHandler(repo repository.UserRepo, jwtSecret string) *AuthHandler {
	return &AuthHandler{repo: repo, jwtSecret: jwtSecret}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string      `json:"token"`
	User  userPublico `json:"user"`
}

type userPublico struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
	Rol      string `json:"rol"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos: " + err.Error()})
		return
	}

	user, err := h.repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, user.Rol, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generando token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Token: token,
		User: userPublico{
			ID:       user.ID,
			Nombre:   user.Nombre,
			Apellido: user.Apellido,
			Email:    user.Email,
			Rol:      user.Rol,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sin autenticación"})
		return
	}

	user, err := h.repo.FindByID(claims.UserID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, userPublico{
		ID:       user.ID,
		Nombre:   user.Nombre,
		Apellido: user.Apellido,
		Email:    user.Email,
		Rol:      user.Rol,
	})
}
