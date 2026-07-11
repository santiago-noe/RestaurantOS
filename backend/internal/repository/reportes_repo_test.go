package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/models"
)

// ─── PedidoRepo.FindEntreFechas ───────────────────────────────────────────────

func TestPedidoRepo_FindEntreFechas_RetornaSoloLosDelRango(t *testing.T) {
	db := setupTestDB(t)

	cliente := &models.Cliente{Nombre: "Cliente Reporte", Tipo: "individual"}
	require.NoError(t, db.Create(cliente).Error)
	user := &models.User{Nombre: "Ana", Apellido: "Lopez", Email: "ana.reporte@test.com", Password: "hash", Rol: "empleado"}
	require.NoError(t, db.Create(user).Error)

	dentro := time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC)
	fuera := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	require.NoError(t, db.Create(&models.Pedido{
		ClienteID: cliente.ID, UserID: user.ID, Fecha: dentro,
		TipoComida: "almuerzo", Estado: "entregado", FormaPago: "contado", Total: 50,
	}).Error)
	require.NoError(t, db.Create(&models.Pedido{
		ClienteID: cliente.ID, UserID: user.ID, Fecha: fuera,
		TipoComida: "almuerzo", Estado: "entregado", FormaPago: "contado", Total: 30,
	}).Error)

	repo := NewPedidoRepo(db)
	desde := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	hasta := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)

	pedidos, err := repo.FindEntreFechas(desde, hasta)

	require.NoError(t, err)
	require.Len(t, pedidos, 1)
	assert.Equal(t, 50.0, pedidos[0].Total)
}

func TestPedidoRepo_FindEntreFechas_SinPedidosRetornaVacio(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPedidoRepo(db)

	pedidos, err := repo.FindEntreFechas(
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC),
	)

	require.NoError(t, err)
	assert.Empty(t, pedidos)
}

// ─── MovimientoRepo.FindEntreFechas ───────────────────────────────────────────

func TestMovimientoRepo_FindEntreFechas_RetornaSoloLosDelRangoConProducto(t *testing.T) {
	db := setupTestDB(t)

	producto := &models.Producto{Nombre: "Papa Andina", Unidad: "kg"}
	require.NoError(t, db.Create(producto).Error)

	dentro := time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC)
	fuera := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	require.NoError(t, db.Create(&models.MovimientoStock{
		ProductoID: producto.ID, Tipo: "entrada", Cantidad: 10, Fecha: dentro,
	}).Error)
	require.NoError(t, db.Create(&models.MovimientoStock{
		ProductoID: producto.ID, Tipo: "salida", Cantidad: 3, Fecha: fuera,
	}).Error)

	repo := NewMovimientoRepo(db)
	desde := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	hasta := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)

	movs, err := repo.FindEntreFechas(desde, hasta)

	require.NoError(t, err)
	require.Len(t, movs, 1)
	assert.Equal(t, "entrada", movs[0].Tipo)
	assert.Equal(t, "Papa Andina", movs[0].Producto.Nombre)
}
