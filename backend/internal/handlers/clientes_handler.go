package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

type ClienteHandler struct {
	repo repository.ClienteRepo
}

func NewClienteHandler(repo repository.ClienteRepo) *ClienteHandler {
	return &ClienteHandler{repo: repo}
}

type crearClienteRequest struct {
	Nombre    string `json:"nombre" binding:"required"`
	Apellido  string `json:"apellido"`
	Tipo      string `json:"tipo" binding:"required"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	Email     string `json:"email"`
}

type actualizarClienteRequest struct {
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	Email     string `json:"email"`
}

func (h *ClienteHandler) Crear(c *gin.Context) {
	var req crearClienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Tipo != "individual" && req.Tipo != "empresa" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo debe ser 'individual' o 'empresa'"})
		return
	}

	if req.Email != "" && h.repo.EmailExists(req.Email, 0) {
		c.JSON(http.StatusConflict, gin.H{"error": "ya existe un cliente con ese email"})
		return
	}

	cliente := &models.Cliente{
		Nombre:    req.Nombre,
		Apellido:  req.Apellido,
		Tipo:      req.Tipo,
		Telefono:  req.Telefono,
		Direccion: req.Direccion,
		Email:     req.Email,
		Activo:    true,
	}

	if err := h.repo.Create(cliente); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear cliente"})
		return
	}

	c.JSON(http.StatusCreated, cliente)
}

func (h *ClienteHandler) ObtenerPorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	cliente, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener cliente"})
		return
	}
	if cliente == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cliente no encontrado"})
		return
	}

	c.JSON(http.StatusOK, cliente)
}

func (h *ClienteHandler) Listar(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	tipo := c.Query("tipo")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	clientes, total, err := h.repo.FindAll(page, perPage, tipo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar clientes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     clientes,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *ClienteHandler) Actualizar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar cliente"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cliente no encontrado"})
		return
	}

	var req actualizarClienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email != "" && h.repo.EmailExists(req.Email, id) {
		c.JSON(http.StatusConflict, gin.H{"error": "ya existe un cliente con ese email"})
		return
	}

	fields := map[string]interface{}{}
	if req.Nombre != "" {
		fields["nombre"] = req.Nombre
	}
	if req.Apellido != "" {
		fields["apellido"] = req.Apellido
	}
	if req.Telefono != "" {
		fields["telefono"] = req.Telefono
	}
	if req.Direccion != "" {
		fields["direccion"] = req.Direccion
	}
	if req.Email != "" {
		fields["email"] = req.Email
	}

	actualizado, err := h.repo.Update(id, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar cliente"})
		return
	}

	c.JSON(http.StatusOK, actualizado)
}

func (h *ClienteHandler) Desactivar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.repo.Deactivate(id); err != nil {
		if err.Error() == "record not found" || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cliente no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al desactivar cliente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "cliente desactivado correctamente"})
}
