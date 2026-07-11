package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"restaurantos/internal/middleware"
	"restaurantos/internal/models"
	"restaurantos/internal/services"
)

// PedidoServiceInterface permite mockear el servicio en los tests.
type PedidoServiceInterface interface {
	Crear(input services.CrearPedidoInput) (*models.Pedido, error)
	Anular(pedidoID int) error
	MarcarEntregado(pedidoID int) error
	FindByID(id int) (*models.Pedido, error)
	FindAll(page, perPage, clienteID int, estado string) ([]models.Pedido, int64, error)
}

type PedidoHandler struct {
	svc PedidoServiceInterface
}

func NewPedidoHandler(svc PedidoServiceInterface) *PedidoHandler {
	return &PedidoHandler{svc: svc}
}

type crearPedidoRequest struct {
	ClienteID  int                  `json:"cliente_id" binding:"required"`
	Fecha      string               `json:"fecha"`
	TipoComida string               `json:"tipo_comida" binding:"required"`
	FormaPago  string               `json:"forma_pago" binding:"required"`
	Notas      string               `json:"notas"`
	Items      []itemPedidoRequest  `json:"items" binding:"required,min=1"`
}

type itemPedidoRequest struct {
	ProductoID     int     `json:"producto_id" binding:"required"`
	Cantidad       float64 `json:"cantidad" binding:"required,gt=0"`
	PrecioUnitario float64 `json:"precio_unitario" binding:"required,gt=0"`
}

func (h *PedidoHandler) Crear(c *gin.Context) {
	var req crearPedidoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fecha := time.Now()
	if req.Fecha != "" {
		if parsed, err := time.Parse("2006-01-02", req.Fecha); err == nil {
			fecha = parsed
		}
	}

	claims := middleware.GetClaims(c)
	items := make([]services.ItemInput, len(req.Items))
	for i, it := range req.Items {
		items[i] = services.ItemInput{
			ProductoID:     it.ProductoID,
			Cantidad:       it.Cantidad,
			PrecioUnitario: it.PrecioUnitario,
		}
	}

	pedido, err := h.svc.Crear(services.CrearPedidoInput{
		ClienteID:  req.ClienteID,
		UserID:     claims.UserID,
		Fecha:      fecha,
		TipoComida: req.TipoComida,
		FormaPago:  req.FormaPago,
		Notas:      req.Notas,
		Items:      items,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pedido)
}

func (h *PedidoHandler) Listar(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	clienteID, _ := strconv.Atoi(c.Query("cliente_id"))
	estado := c.Query("estado")

	if page < 1 {
		page = 1
	}

	pedidos, total, err := h.svc.FindAll(page, perPage, clienteID, estado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar pedidos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     pedidos,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *PedidoHandler) ObtenerPorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	pedido, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener pedido"})
		return
	}
	if pedido == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pedido no encontrado"})
		return
	}

	c.JSON(http.StatusOK, pedido)
}

func (h *PedidoHandler) MarcarEntregado(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.svc.MarcarEntregado(id); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "pedido marcado como entregado"})
}

func (h *PedidoHandler) Anular(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.svc.Anular(id); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "pedido anulado correctamente"})
}
