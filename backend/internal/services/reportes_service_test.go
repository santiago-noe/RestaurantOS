package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"restaurantos/internal/models"
)

// ─── RangoPeriodo ──────────────────────────────────────────────────────────────

func TestRangoPeriodo_Diario_RetornaMismoDia(t *testing.T) {
	ref := time.Date(2026, 7, 15, 14, 30, 0, 0, time.UTC)
	desde, hasta := RangoPeriodo("diario", ref)

	assert.Equal(t, "2026-07-15", desde.Format("2006-01-02"))
	assert.Equal(t, "2026-07-15", hasta.Format("2006-01-02"))
}

func TestRangoPeriodo_Semanal_RetornaLunesADomingo(t *testing.T) {
	// 2026-07-15 es un miércoles
	ref := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	desde, hasta := RangoPeriodo("semanal", ref)

	assert.Equal(t, "2026-07-13", desde.Format("2006-01-02")) // lunes
	assert.Equal(t, "2026-07-19", hasta.Format("2006-01-02")) // domingo
}

func TestRangoPeriodo_Mensual_RetornaMesCompleto(t *testing.T) {
	ref := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
	desde, hasta := RangoPeriodo("mensual", ref)

	assert.Equal(t, "2026-02-01", desde.Format("2006-01-02"))
	assert.Equal(t, "2026-02-28", hasta.Format("2006-01-02"))
}

// ─── GenerarReporteVentas ──────────────────────────────────────────────────────

func TestGenerarReporteVentas_SumaTotalesPorDia(t *testing.T) {
	dia1 := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	dia2 := time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC)
	pedidos := []models.Pedido{
		{Fecha: dia1, Total: 50, Estado: "entregado"},
		{Fecha: dia1, Total: 30, Estado: "entregado"},
		{Fecha: dia2, Total: 20, Estado: "pendiente"},
	}

	r := GenerarReporteVentas(pedidos, "semanal", dia1, dia2)

	assert.Equal(t, 100.0, r.TotalVentas)
	assert.Equal(t, 3, r.TotalPedidos)
	assert.Len(t, r.PorDia, 2)
	assert.Equal(t, 80.0, r.PorDia[0].Total)
	assert.Equal(t, 2, r.PorDia[0].CantidadPedidos)
}

func TestGenerarReporteVentas_ExcluyeAnulados(t *testing.T) {
	dia := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	pedidos := []models.Pedido{
		{Fecha: dia, Total: 50, Estado: "entregado"},
		{Fecha: dia, Total: 999, Estado: "anulado"},
	}

	r := GenerarReporteVentas(pedidos, "diario", dia, dia)

	assert.Equal(t, 50.0, r.TotalVentas)
	assert.Equal(t, 1, r.TotalPedidos)
}

func TestGenerarReporteVentas_SinPedidosRetornaCeros(t *testing.T) {
	dia := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	r := GenerarReporteVentas([]models.Pedido{}, "diario", dia, dia)

	assert.Equal(t, 0.0, r.TotalVentas)
	assert.Equal(t, 0, r.TotalPedidos)
	assert.Empty(t, r.PorDia)
}

func TestGenerarReporteVentas_OrdenaPorDiaAscendente(t *testing.T) {
	dia1 := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	dia2 := time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC)
	pedidos := []models.Pedido{
		{Fecha: dia2, Total: 10, Estado: "entregado"},
		{Fecha: dia1, Total: 20, Estado: "entregado"},
	}

	r := GenerarReporteVentas(pedidos, "semanal", dia1, dia2)

	assert.Equal(t, "2026-07-01", r.PorDia[0].Fecha)
	assert.Equal(t, "2026-07-02", r.PorDia[1].Fecha)
}

// ─── FiltrarClientesConDeuda ───────────────────────────────────────────────────

func TestFiltrarClientesConDeuda_RetornaSoloConDeudaMayorACero(t *testing.T) {
	clientes := []models.Cliente{
		{Nombre: "Con deuda", DeudaTotal: 100},
		{Nombre: "Sin deuda", DeudaTotal: 0},
	}

	deudores := FiltrarClientesConDeuda(clientes)

	assert.Len(t, deudores, 1)
	assert.Equal(t, "Con deuda", deudores[0].Nombre)
}

func TestFiltrarClientesConDeuda_SinDeudoresRetornaVacio(t *testing.T) {
	clientes := []models.Cliente{{Nombre: "Sin deuda", DeudaTotal: 0}}

	deudores := FiltrarClientesConDeuda(clientes)

	assert.Empty(t, deudores)
}

// ─── GenerarReporteMovimientos ─────────────────────────────────────────────────

func TestGenerarReporteMovimientos_SumaEntradasYSalidasPorSeparado(t *testing.T) {
	desde := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	hasta := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)
	movs := []models.MovimientoStock{
		{Tipo: "entrada", Cantidad: 10},
		{Tipo: "entrada", Cantidad: 5},
		{Tipo: "salida", Cantidad: 3},
	}

	r := GenerarReporteMovimientos(movs, desde, hasta)

	assert.Equal(t, 15.0, r.TotalEntradas)
	assert.Equal(t, 3.0, r.TotalSalidas)
	assert.Len(t, r.Movimientos, 3)
}

func TestGenerarReporteMovimientos_SinMovimientosRetornaCeros(t *testing.T) {
	desde := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	hasta := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)

	r := GenerarReporteMovimientos([]models.MovimientoStock{}, desde, hasta)

	assert.Equal(t, 0.0, r.TotalEntradas)
	assert.Equal(t, 0.0, r.TotalSalidas)
}
