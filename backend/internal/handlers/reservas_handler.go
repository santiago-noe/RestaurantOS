package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

type ReservaHandler struct {
	repo       repository.ReservaRepo
	pedidoRepo repository.PedidoRepo
}

func NewReservaHandler(repo repository.ReservaRepo, pedidoRepo repository.PedidoRepo) *ReservaHandler {
	return &ReservaHandler{repo: repo, pedidoRepo: pedidoRepo}
}

type crearReservaRequest struct {
	Nombre   string `json:"nombre" binding:"required"`
	Whatsapp string `json:"whatsapp" binding:"required"`
	Fecha    string `json:"fecha" binding:"required"`
	Personas string `json:"personas" binding:"required"`
	Ocasion  string `json:"ocasion"`
}

type actualizarEstadoReservaRequest struct {
	Estado string `json:"estado" binding:"required"`
}

type vincularPedidoRequest struct {
	PedidoID int `json:"pedido_id" binding:"required"`
}

// Crear — POST /api/public/reservas (sin autenticación)
func (h *ReservaHandler) Crear(c *gin.Context) {
	var req crearReservaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha inválida, formato esperado YYYY-MM-DD"})
		return
	}

	reserva := &models.Reserva{
		Nombre:   req.Nombre,
		Whatsapp: req.Whatsapp,
		Fecha:    fecha,
		Personas: req.Personas,
		Ocasion:  req.Ocasion,
		Estado:   "pendiente",
	}

	if err := h.repo.Create(reserva); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear la reserva"})
		return
	}

	c.JSON(http.StatusCreated, reserva)
}

// Listar — GET /api/empleado/reservas
func (h *ReservaHandler) Listar(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	estado := c.Query("estado")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	reservas, total, err := h.repo.FindAll(page, perPage, estado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar reservas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     reservas,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// ActualizarEstado — PUT /api/empleado/reservas/:id/estado
func (h *ReservaHandler) ActualizarEstado(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req actualizarEstadoReservaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Estado != "pendiente" && req.Estado != "confirmada" && req.Estado != "cancelada" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "estado debe ser 'pendiente', 'confirmada' o 'cancelada'"})
		return
	}

	if err := h.repo.UpdateEstado(id, req.Estado); err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "reserva no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar la reserva"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "reserva actualizada correctamente"})
}

// VincularPedido — PUT /api/empleado/reservas/:id/pedido
// Se usa cuando el cliente llega y el mesero registra el Pedido real: queda
// enlazado a la reserva solo para trazabilidad, sin que la reserva controle stock ni totales.
func (h *ReservaHandler) VincularPedido(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req vincularPedidoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pedido, err := h.pedidoRepo.FindByID(req.PedidoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar el pedido"})
		return
	}
	if pedido == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pedido no encontrado"})
		return
	}

	if err := h.repo.VincularPedido(id, req.PedidoID); err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "reserva no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al vincular el pedido"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "reserva vinculada al pedido correctamente"})
}
