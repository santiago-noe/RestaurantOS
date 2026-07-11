package services

import (
	"sort"
	"time"

	"restaurantos/internal/models"
)

// ─── Reporte de ventas ────────────────────────────────────────────────────────

type VentaDia struct {
	Fecha           string  `json:"fecha"`
	Total           float64 `json:"total"`
	CantidadPedidos int     `json:"cantidad_pedidos"`
}

type ReporteVentas struct {
	Periodo      string     `json:"periodo"`
	Desde        string     `json:"desde"`
	Hasta        string     `json:"hasta"`
	TotalVentas  float64    `json:"total_ventas"`
	TotalPedidos int        `json:"total_pedidos"`
	PorDia       []VentaDia `json:"por_dia"`
}

// RangoPeriodo calcula el rango [desde, hasta] para "diario", "semanal" o "mensual"
// tomando como referencia la fecha dada.
func RangoPeriodo(periodo string, referencia time.Time) (time.Time, time.Time) {
	inicioDia := time.Date(referencia.Year(), referencia.Month(), referencia.Day(), 0, 0, 0, 0, referencia.Location())

	switch periodo {
	case "semanal":
		// Lunes como inicio de semana
		offset := (int(inicioDia.Weekday()) + 6) % 7
		desde := inicioDia.AddDate(0, 0, -offset)
		hasta := desde.AddDate(0, 0, 6)
		return desde, hasta
	case "mensual":
		desde := time.Date(referencia.Year(), referencia.Month(), 1, 0, 0, 0, 0, referencia.Location())
		hasta := desde.AddDate(0, 1, -1)
		return desde, hasta
	default: // "diario"
		return inicioDia, inicioDia
	}
}

// GenerarReporteVentas agrupa los pedidos (ya filtrados por rango) por día,
// excluyendo los anulados.
func GenerarReporteVentas(pedidos []models.Pedido, periodo string, desde, hasta time.Time) ReporteVentas {
	porDia := map[string]*VentaDia{}
	totalVentas := 0.0
	totalPedidos := 0

	for _, p := range pedidos {
		if p.Estado == "anulado" {
			continue
		}
		key := p.Fecha.Format("2006-01-02")
		if _, ok := porDia[key]; !ok {
			porDia[key] = &VentaDia{Fecha: key}
		}
		porDia[key].Total += p.Total
		porDia[key].CantidadPedidos++
		totalVentas += p.Total
		totalPedidos++
	}

	claves := make([]string, 0, len(porDia))
	for k := range porDia {
		claves = append(claves, k)
	}
	sort.Strings(claves)

	porDiaLista := make([]VentaDia, 0, len(claves))
	for _, k := range claves {
		porDiaLista = append(porDiaLista, *porDia[k])
	}

	return ReporteVentas{
		Periodo:      periodo,
		Desde:        desde.Format("2006-01-02"),
		Hasta:        hasta.Format("2006-01-02"),
		TotalVentas:  totalVentas,
		TotalPedidos: totalPedidos,
		PorDia:       porDiaLista,
	}
}

// ─── Reporte de clientes con deuda ────────────────────────────────────────────

func FiltrarClientesConDeuda(clientes []models.Cliente) []models.Cliente {
	deudores := []models.Cliente{}
	for _, c := range clientes {
		if c.DeudaTotal > 0 {
			deudores = append(deudores, c)
		}
	}
	return deudores
}

// ─── Reporte de movimientos de inventario ─────────────────────────────────────

type ReporteMovimientos struct {
	Desde         string                    `json:"desde"`
	Hasta         string                    `json:"hasta"`
	TotalEntradas float64                   `json:"total_entradas"`
	TotalSalidas  float64                   `json:"total_salidas"`
	Movimientos   []models.MovimientoStock  `json:"movimientos"`
}

func GenerarReporteMovimientos(movs []models.MovimientoStock, desde, hasta time.Time) ReporteMovimientos {
	var entradas, salidas float64
	for _, m := range movs {
		if m.Tipo == "entrada" {
			entradas += m.Cantidad
		} else {
			salidas += m.Cantidad
		}
	}
	return ReporteMovimientos{
		Desde:         desde.Format("2006-01-02"),
		Hasta:         hasta.Format("2006-01-02"),
		TotalEntradas: entradas,
		TotalSalidas:  salidas,
		Movimientos:   movs,
	}
}
