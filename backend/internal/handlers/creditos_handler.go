package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

type CreditosHandler struct {
	pagoRepo    repository.PagoRepo
	clienteRepo repository.ClienteRepo
}

func NewCreditosHandler(pr repository.PagoRepo, cr repository.ClienteRepo) *CreditosHandler {
	return &CreditosHandler{pagoRepo: pr, clienteRepo: cr}
}

func (h *CreditosHandler) ListarDeudores(c *gin.Context) {
	clientes, _, err := h.clienteRepo.FindAll(1, 100, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener deudores"})
		return
	}

	deudores := []models.Cliente{}
	for _, cl := range clientes {
		if cl.DeudaTotal > 0 {
			deudores = append(deudores, cl)
		}
	}

	c.JSON(http.StatusOK, gin.H{"deudores": deudores, "total": len(deudores)})
}

func (h *CreditosHandler) EstadoCuenta(c *gin.Context) {
	clienteID, err := strconv.Atoi(c.Param("cliente_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de cliente inválido"})
		return
	}

	cliente, err := h.clienteRepo.FindByID(clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar cliente"})
		return
	}
	if cliente == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cliente no encontrado"})
		return
	}

	pagos, err := h.pagoRepo.FindByCliente(clienteID)
	if err != nil {
		pagos = []models.Pago{}
	}

	c.JSON(http.StatusOK, gin.H{
		"cliente":          cliente,
		"deuda_total":      cliente.DeudaTotal,
		"pagos_realizados": pagos,
	})
}

type registrarPagoReq struct {
	ClienteID int     `json:"cliente_id" binding:"required"`
	Monto     float64 `json:"monto" binding:"required,gt=0"`
	Metodo    string  `json:"metodo" binding:"required"`
	PedidoID  *int    `json:"pedido_id"`
}

func (h *CreditosHandler) RegistrarPago(c *gin.Context) {
	var req registrarPagoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cliente, err := h.clienteRepo.FindByID(req.ClienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar cliente"})
		return
	}
	if cliente == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cliente no encontrado"})
		return
	}
	if cliente.DeudaTotal <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "el cliente no tiene deuda pendiente"})
		return
	}
	if req.Monto > cliente.DeudaTotal {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "el monto excede la deuda del cliente",
			"deuda": cliente.DeudaTotal,
		})
		return
	}

	pago := &models.Pago{
		ClienteID: req.ClienteID,
		Monto:     req.Monto,
		Metodo:    req.Metodo,
		PedidoID:  req.PedidoID,
	}
	if err := h.pagoRepo.Create(pago); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al registrar pago"})
		return
	}

	c.JSON(http.StatusCreated, pago)
}
