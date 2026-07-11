package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restaurantos/internal/export"
	"restaurantos/internal/repository"
	"restaurantos/internal/services"
)

type ReportesHandler struct {
	pedidoRepo     repository.PedidoRepo
	movimientoRepo repository.MovimientoRepo
	clienteRepo    repository.ClienteRepo
}

func NewReportesHandler(pr repository.PedidoRepo, mr repository.MovimientoRepo, cr repository.ClienteRepo) *ReportesHandler {
	return &ReportesHandler{pedidoRepo: pr, movimientoRepo: mr, clienteRepo: cr}
}

func parseFecha(valor string, porDefecto time.Time) time.Time {
	if valor == "" {
		return porDefecto
	}
	t, err := time.Parse("2006-01-02", valor)
	if err != nil {
		return porDefecto
	}
	return t
}

// Ventas: GET /api/admin/reportes/ventas?periodo=diario|semanal|mensual&fecha=YYYY-MM-DD&formato=pdf|excel
func (h *ReportesHandler) Ventas(c *gin.Context) {
	periodo := c.DefaultQuery("periodo", "diario")
	if periodo != "diario" && periodo != "semanal" && periodo != "mensual" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "periodo debe ser 'diario', 'semanal' o 'mensual'"})
		return
	}

	referencia := parseFecha(c.Query("fecha"), time.Now())
	desde, hasta := services.RangoPeriodo(periodo, referencia)
	// hasta es el inicio del último día del rango; se extiende hasta el final de ese día para incluir todos los pedidos
	hastaFinDia := time.Date(hasta.Year(), hasta.Month(), hasta.Day(), 23, 59, 59, 0, hasta.Location())

	pedidos, err := h.pedidoRepo.FindEntreFechas(desde, hastaFinDia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener las ventas"})
		return
	}

	reporte := services.GenerarReporteVentas(pedidos, periodo, desde, hasta)

	switch c.Query("formato") {
	case "excel":
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=reporte-ventas-%s.xlsx", reporte.Desde))
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		if err := export.VentasExcel(c.Writer, reporte); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar el Excel"})
		}
	case "pdf":
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=reporte-ventas-%s.pdf", reporte.Desde))
		c.Header("Content-Type", "application/pdf")
		if err := export.VentasPDF(c.Writer, reporte); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar el PDF"})
		}
	default:
		c.JSON(http.StatusOK, reporte)
	}
}

// Deudores: GET /api/admin/reportes/deudores?formato=pdf|excel
func (h *ReportesHandler) Deudores(c *gin.Context) {
	clientes, _, err := h.clienteRepo.FindAll(1, 1000, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener los clientes"})
		return
	}
	deudores := services.FiltrarClientesConDeuda(clientes)

	switch c.Query("formato") {
	case "excel":
		c.Header("Content-Disposition", "attachment; filename=reporte-deudores.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		if err := export.DeudoresExcel(c.Writer, deudores); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar el Excel"})
		}
	case "pdf":
		c.Header("Content-Disposition", "attachment; filename=reporte-deudores.pdf")
		c.Header("Content-Type", "application/pdf")
		if err := export.DeudoresPDF(c.Writer, deudores); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar el PDF"})
		}
	default:
		c.JSON(http.StatusOK, gin.H{"deudores": deudores, "total": len(deudores)})
	}
}

// Inventario: GET /api/admin/reportes/inventario?desde=YYYY-MM-DD&hasta=YYYY-MM-DD&formato=pdf|excel
func (h *ReportesHandler) Inventario(c *gin.Context) {
	hastaDefecto := time.Now()
	desdeDefecto := hastaDefecto.AddDate(0, 0, -30)

	desde := parseFecha(c.Query("desde"), desdeDefecto)
	hasta := parseFecha(c.Query("hasta"), hastaDefecto)
	hastaFinDia := time.Date(hasta.Year(), hasta.Month(), hasta.Day(), 23, 59, 59, 0, hasta.Location())

	movs, err := h.movimientoRepo.FindEntreFechas(desde, hastaFinDia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener los movimientos"})
		return
	}

	reporte := services.GenerarReporteMovimientos(movs, desde, hasta)

	switch c.Query("formato") {
	case "excel":
		c.Header("Content-Disposition", "attachment; filename=reporte-inventario.xlsx")
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		if err := export.MovimientosExcel(c.Writer, reporte); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar el Excel"})
		}
	case "pdf":
		c.Header("Content-Disposition", "attachment; filename=reporte-inventario.pdf")
		c.Header("Content-Type", "application/pdf")
		if err := export.MovimientosPDF(c.Writer, reporte); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar el PDF"})
		}
	default:
		c.JSON(http.StatusOK, reporte)
	}
}
