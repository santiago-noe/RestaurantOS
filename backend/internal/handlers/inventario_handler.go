package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

var unidadesValidas = map[string]bool{
	"kg": true, "litro": true, "unidad": true, "porcion": true,
}

type InventarioHandler struct {
	productoRepo   repository.ProductoRepo
	movimientoRepo repository.MovimientoRepo
}

func NewInventarioHandler(pr repository.ProductoRepo, mr repository.MovimientoRepo) *InventarioHandler {
	return &InventarioHandler{productoRepo: pr, movimientoRepo: mr}
}

func (h *InventarioHandler) Listar(c *gin.Context) {
	productos, err := h.productoRepo.FindAll(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar productos"})
		return
	}
	c.JSON(http.StatusOK, productos)
}

func (h *InventarioHandler) Alertas(c *gin.Context) {
	todos, err := h.productoRepo.FindAll(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener alertas"})
		return
	}

	type alerta struct {
		ID          int     `json:"id"`
		Nombre      string  `json:"nombre"`
		StockActual float64 `json:"stock_actual"`
		StockMinimo float64 `json:"stock_minimo"`
		Unidad      string  `json:"unidad"`
		Deficit     float64 `json:"deficit"`
	}

	alertas := []alerta{}
	for _, p := range todos {
		if p.StockActual < p.StockMinimo {
			alertas = append(alertas, alerta{
				ID:          p.ID,
				Nombre:      p.Nombre,
				StockActual: p.StockActual,
				StockMinimo: p.StockMinimo,
				Unidad:      p.Unidad,
				Deficit:     p.StockMinimo - p.StockActual,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"alertas": alertas})
}

func (h *InventarioHandler) ObtenerPorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	producto, err := h.productoRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener producto"})
		return
	}
	if producto == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
		return
	}

	movimientos, err := h.movimientoRepo.FindByProducto(id)
	if err != nil {
		movimientos = []models.MovimientoStock{}
	}

	c.JSON(http.StatusOK, gin.H{"producto": producto, "movimientos": movimientos})
}

type crearProductoReq struct {
	Nombre      string  `json:"nombre" binding:"required"`
	Unidad      string  `json:"unidad" binding:"required"`
	StockMinimo float64 `json:"stock_minimo"`
	PrecioVenta float64 `json:"precio_venta"`
}

func (h *InventarioHandler) Crear(c *gin.Context) {
	var req crearProductoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !unidadesValidas[req.Unidad] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unidad inválida (kg, litro, unidad, porcion)"})
		return
	}

	p := &models.Producto{
		Nombre:      req.Nombre,
		Unidad:      req.Unidad,
		StockMinimo: req.StockMinimo,
		PrecioVenta: req.PrecioVenta,
		Activo:      true,
	}
	if err := h.productoRepo.Create(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear producto"})
		return
	}
	c.JSON(http.StatusCreated, p)
}

type actualizarProductoReq struct {
	Nombre      string  `json:"nombre"`
	StockMinimo float64 `json:"stock_minimo"`
	PrecioVenta float64 `json:"precio_venta"`
}

func (h *InventarioHandler) Actualizar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	existing, err := h.productoRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar producto"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
		return
	}

	var req actualizarProductoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fields := map[string]interface{}{}
	if req.Nombre != "" {
		fields["nombre"] = req.Nombre
	}
	if req.StockMinimo > 0 {
		fields["stock_minimo"] = req.StockMinimo
	}
	if req.PrecioVenta > 0 {
		fields["precio_venta"] = req.PrecioVenta
	}

	actualizado, err := h.productoRepo.Update(id, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar"})
		return
	}
	c.JSON(http.StatusOK, actualizado)
}

type restockReq struct {
	Cantidad float64 `json:"cantidad" binding:"required,gt=0"`
	Notas    string  `json:"notas"`
}

func (h *InventarioHandler) Restock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req restockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cantidad debe ser mayor a 0"})
		return
	}

	producto, err := h.productoRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar producto"})
		return
	}
	if producto == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
		return
	}

	if err := h.productoRepo.AjustarStock(id, req.Cantidad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al ajustar stock"})
		return
	}

	mov := &models.MovimientoStock{
		ProductoID: id,
		Tipo:       "entrada",
		Cantidad:   req.Cantidad,
		Notas:      req.Notas,
		Fecha:      time.Now(),
	}
	_ = h.movimientoRepo.Registrar(mov)

	c.JSON(http.StatusOK, gin.H{
		"producto_id":      id,
		"nombre":           producto.Nombre,
		"stock_anterior":   producto.StockActual,
		"cantidad_agregada": req.Cantidad,
		"stock_actual":     producto.StockActual + req.Cantidad,
		"unidad":           producto.Unidad,
	})
}
