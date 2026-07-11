package export

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/models"
	"restaurantos/internal/services"
)

var reporteVentasEjemplo = services.ReporteVentas{
	Periodo: "semanal", Desde: "2026-07-01", Hasta: "2026-07-07",
	TotalVentas: 150, TotalPedidos: 3,
	PorDia: []services.VentaDia{{Fecha: "2026-07-01", Total: 150, CantidadPedidos: 3}},
}

var clientesDeudoresEjemplo = []models.Cliente{
	{Nombre: "Juan", Apellido: "Perez", Tipo: "individual", Telefono: "999999999", DeudaTotal: 120},
}

var reporteMovimientosEjemplo = services.ReporteMovimientos{
	Desde: "2026-07-01", Hasta: "2026-07-31",
	TotalEntradas: 10, TotalSalidas: 3,
	Movimientos: []models.MovimientoStock{
		{Tipo: "entrada", Cantidad: 10, Fecha: time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC), Producto: models.Producto{Nombre: "Papa"}},
	},
}

func TestVentasExcel_GeneraArchivoValido(t *testing.T) {
	var buf bytes.Buffer
	err := VentasExcel(&buf, reporteVentasEjemplo)

	require.NoError(t, err)
	assert.NotEmpty(t, buf.Bytes())
	assert.Equal(t, "PK", string(buf.Bytes()[:2])) // firma de archivo .xlsx (zip)
}

func TestVentasPDF_GeneraArchivoValido(t *testing.T) {
	var buf bytes.Buffer
	err := VentasPDF(&buf, reporteVentasEjemplo)

	require.NoError(t, err)
	assert.NotEmpty(t, buf.Bytes())
	assert.Equal(t, "%PDF", string(buf.Bytes()[:4]))
}

func TestDeudoresExcel_GeneraArchivoValido(t *testing.T) {
	var buf bytes.Buffer
	err := DeudoresExcel(&buf, clientesDeudoresEjemplo)

	require.NoError(t, err)
	assert.Equal(t, "PK", string(buf.Bytes()[:2]))
}

func TestDeudoresPDF_GeneraArchivoValido(t *testing.T) {
	var buf bytes.Buffer
	err := DeudoresPDF(&buf, clientesDeudoresEjemplo)

	require.NoError(t, err)
	assert.Equal(t, "%PDF", string(buf.Bytes()[:4]))
}

func TestMovimientosExcel_GeneraArchivoValido(t *testing.T) {
	var buf bytes.Buffer
	err := MovimientosExcel(&buf, reporteMovimientosEjemplo)

	require.NoError(t, err)
	assert.Equal(t, "PK", string(buf.Bytes()[:2]))
}

func TestMovimientosPDF_GeneraArchivoValido(t *testing.T) {
	var buf bytes.Buffer
	err := MovimientosPDF(&buf, reporteMovimientosEjemplo)

	require.NoError(t, err)
	assert.Equal(t, "%PDF", string(buf.Bytes()[:4]))
}

func TestVentasExcel_SinDatosNoFalla(t *testing.T) {
	var buf bytes.Buffer
	err := VentasExcel(&buf, services.ReporteVentas{})

	require.NoError(t, err)
	assert.NotEmpty(t, buf.Bytes())
}
