package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

type MenuHandler struct {
	repo repository.MenuRepo
}

func NewMenuHandler(repo repository.MenuRepo) *MenuHandler {
	return &MenuHandler{repo: repo}
}

type crearMenuRequest struct {
	Categoria   string  `json:"categoria" binding:"required"`
	Nombre      string  `json:"nombre" binding:"required"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	ImagenURL   string  `json:"imagen_url"`
	Disponible  *bool   `json:"disponible"`
	Orden       int     `json:"orden"`
	ProductoID  *int    `json:"producto_id"`
}

type actualizarMenuRequest struct {
	Categoria   string  `json:"categoria"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	ImagenURL   string  `json:"imagen_url"`
	Disponible  *bool   `json:"disponible"`
	Orden       *int    `json:"orden"`
	ProductoID  *int    `json:"producto_id"`
}

// Publico devuelve solo los items disponibles, para la landing page (sin auth).
func (h *MenuHandler) Publico(c *gin.Context) {
	items, err := h.repo.FindPublico()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener el menú"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Listar devuelve todos los items (disponibles o no), para el dashboard admin.
func (h *MenuHandler) Listar(c *gin.Context) {
	items, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar el menú"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *MenuHandler) Crear(c *gin.Context) {
	var req crearMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	disponible := true
	if req.Disponible != nil {
		disponible = *req.Disponible
	}

	item := &models.MenuPublico{
		Categoria:   req.Categoria,
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Precio:      req.Precio,
		ImagenURL:   req.ImagenURL,
		Disponible:  disponible,
		Orden:       req.Orden,
		ProductoID:  req.ProductoID,
	}

	if err := h.repo.Create(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear el item del menú"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *MenuHandler) Actualizar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar el item del menú"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item del menú no encontrado"})
		return
	}

	var req actualizarMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fields := map[string]interface{}{}
	if req.Categoria != "" {
		fields["categoria"] = req.Categoria
	}
	if req.Nombre != "" {
		fields["nombre"] = req.Nombre
	}
	if req.Descripcion != "" {
		fields["descripcion"] = req.Descripcion
	}
	if req.Precio != 0 {
		fields["precio"] = req.Precio
	}
	if req.ImagenURL != "" {
		fields["imagen_url"] = req.ImagenURL
	}
	if req.Disponible != nil {
		fields["disponible"] = *req.Disponible
	}
	if req.Orden != nil {
		fields["orden"] = *req.Orden
	}
	if req.ProductoID != nil {
		fields["producto_id"] = *req.ProductoID
	}

	actualizado, err := h.repo.Update(id, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar el item del menú"})
		return
	}

	c.JSON(http.StatusOK, actualizado)
}

func (h *MenuHandler) Eliminar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "item del menú no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar el item del menú"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "item eliminado correctamente"})
}
