package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/models"
)

// ─── PredecirDemanda ─────────────────────────────────────────────────────────

func TestPredecirDemanda_SinHistorialRetornaCero(t *testing.T) {
	pred := PredecirDemanda([]models.Pedido{}, time.Now().AddDate(0, 0, 1))
	for _, p := range pred {
		assert.Equal(t, 0, p.PorcionesEstimadas)
	}
}

func TestPredecirDemanda_HistorialRetornaPromedioPorDia(t *testing.T) {
	// 3 lunes con 10 almuerzos cada uno → promedio 10
	lunes1 := time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC)  // lunes
	lunes2 := time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC) // lunes
	lunes3 := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC) // lunes

	pedidos := []models.Pedido{
		{Fecha: lunes1, TipoComida: "almuerzo"},
		{Fecha: lunes1, TipoComida: "almuerzo"},
		{Fecha: lunes2, TipoComida: "almuerzo"},
		{Fecha: lunes2, TipoComida: "almuerzo"},
		{Fecha: lunes3, TipoComida: "almuerzo"},
		{Fecha: lunes3, TipoComida: "almuerzo"},
	}

	// Predicción para el próximo lunes
	proximoLunes := time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)
	pred := PredecirDemanda(pedidos, proximoLunes)

	var almuerzo *Prediccion
	for i := range pred {
		if pred[i].TipoComida == "almuerzo" {
			almuerzo = &pred[i]
			break
		}
	}

	require.NotNil(t, almuerzo)
	assert.Equal(t, 2, almuerzo.PorcionesEstimadas) // 6 pedidos / 3 lunes = 2
}

func TestPredecirDemanda_SoloUsaDiaSemanaCorrespondiente(t *testing.T) {
	lunes := time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC)
	martes := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)

	pedidos := []models.Pedido{
		{Fecha: lunes, TipoComida: "almuerzo"},
		{Fecha: martes, TipoComida: "almuerzo"},
		{Fecha: martes, TipoComida: "almuerzo"},
		{Fecha: martes, TipoComida: "almuerzo"},
	}

	// Predecir para el próximo lunes
	proximoLunes := time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)
	pred := PredecirDemanda(pedidos, proximoLunes)

	var almuerzo *Prediccion
	for i := range pred {
		if pred[i].TipoComida == "almuerzo" {
			almuerzo = &pred[i]
		}
	}
	require.NotNil(t, almuerzo)
	// Solo cuenta los lunes: 1 pedido / 1 lunes = 1 (no los martes)
	assert.Equal(t, 1, almuerzo.PorcionesEstimadas)
}

// ─── GenerarAlertas ──────────────────────────────────────────────────────────

func TestGenerarAlertas_DeudaAltaApareceEnAlertas(t *testing.T) {
	clientes := []models.Cliente{
		{ID: 1, Nombre: "Juan", DeudaTotal: 250.00},
		{ID: 2, Nombre: "Ana", DeudaTotal: 50.00},
	}
	alertas := GenerarAlertas(clientes, []models.Producto{}, 200.00)

	deudas := filtrarTipo(alertas, "deuda_alta")
	assert.Len(t, deudas, 1)
	assert.Contains(t, deudas[0].Mensaje, "Juan")
}

func TestGenerarAlertas_DeudaBajaNoApareceEnAlertas(t *testing.T) {
	clientes := []models.Cliente{{ID: 1, Nombre: "Ana", DeudaTotal: 50.00}}
	alertas := GenerarAlertas(clientes, []models.Producto{}, 200.00)

	deudas := filtrarTipo(alertas, "deuda_alta")
	assert.Len(t, deudas, 0)
}

func TestGenerarAlertas_StockBajoApareceEnAlertas(t *testing.T) {
	productos := []models.Producto{
		{ID: 1, Nombre: "Aceite", StockActual: 0.5, StockMinimo: 2.0},
		{ID: 2, Nombre: "Arroz", StockActual: 10.0, StockMinimo: 5.0},
	}
	alertas := GenerarAlertas([]models.Cliente{}, productos, 200.00)

	stock := filtrarTipo(alertas, "stock_bajo")
	assert.Len(t, stock, 1)
	assert.Contains(t, stock[0].Mensaje, "Aceite")
}

func TestGenerarAlertas_OrdenadaPorSeveridad(t *testing.T) {
	clientes := []models.Cliente{{ID: 1, DeudaTotal: 500.00}} // alta severidad
	productos := []models.Producto{{ID: 1, StockActual: 0, StockMinimo: 5}} // media

	alertas := GenerarAlertas(clientes, productos, 200.00)
	require.GreaterOrEqual(t, len(alertas), 2)
	assert.Equal(t, "alta", alertas[0].Severidad)
}

func TestGenerarAlertas_SinProblemasRetornaVacio(t *testing.T) {
	clientes := []models.Cliente{{ID: 1, DeudaTotal: 0}}
	productos := []models.Producto{{ID: 1, StockActual: 10, StockMinimo: 5}}

	alertas := GenerarAlertas(clientes, productos, 200.00)
	assert.Empty(t, alertas)
}

// helper de test
func filtrarTipo(alertas []Alerta, tipo string) []Alerta {
	var r []Alerta
	for _, a := range alertas {
		if a.Tipo == tipo {
			r = append(r, a)
		}
	}
	return r
}
