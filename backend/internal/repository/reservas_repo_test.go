package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/models"
)

func TestReservaRepo_Create_GuardaReservaEnBD(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReservaRepo(db)

	reserva := &models.Reserva{
		Nombre:   "Juan Perez",
		Whatsapp: "987654321",
		Fecha:    time.Now(),
		Personas: "4",
		Estado:   "pendiente",
	}

	err := repo.Create(reserva)

	require.NoError(t, err)
	assert.NotZero(t, reserva.ID)
}

func TestReservaRepo_FindByID_ExisteRetornaDatos(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReservaRepo(db)

	reserva := &models.Reserva{Nombre: "Ana Lopez", Whatsapp: "999888777", Fecha: time.Now(), Personas: "2", Estado: "pendiente"}
	require.NoError(t, repo.Create(reserva))

	encontrada, err := repo.FindByID(reserva.ID)

	require.NoError(t, err)
	require.NotNil(t, encontrada)
	assert.Equal(t, "Ana Lopez", encontrada.Nombre)
}

func TestReservaRepo_FindByID_InexistenteRetornaNil(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReservaRepo(db)

	encontrada, err := repo.FindByID(999999)

	require.NoError(t, err)
	assert.Nil(t, encontrada)
}

func TestReservaRepo_FindAll_RetornaPaginadoYFiltraPorEstado(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReservaRepo(db)

	require.NoError(t, repo.Create(&models.Reserva{Nombre: "Res A", Whatsapp: "1", Fecha: time.Now(), Personas: "2", Estado: "pendiente"}))
	require.NoError(t, repo.Create(&models.Reserva{Nombre: "Res B", Whatsapp: "2", Fecha: time.Now(), Personas: "3", Estado: "confirmada"}))

	todas, total, err := repo.FindAll(1, 10, "")
	require.NoError(t, err)
	assert.EqualValues(t, 2, total)
	assert.Len(t, todas, 2)

	pendientes, totalPendientes, err := repo.FindAll(1, 10, "pendiente")
	require.NoError(t, err)
	assert.EqualValues(t, 1, totalPendientes)
	assert.Len(t, pendientes, 1)
	assert.Equal(t, "Res A", pendientes[0].Nombre)
}

func TestReservaRepo_UpdateEstado_ActualizaCorrectamente(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReservaRepo(db)

	reserva := &models.Reserva{Nombre: "Carlos", Whatsapp: "3", Fecha: time.Now(), Personas: "5", Estado: "pendiente"}
	require.NoError(t, repo.Create(reserva))

	err := repo.UpdateEstado(reserva.ID, "confirmada")
	require.NoError(t, err)

	actualizada, err := repo.FindByID(reserva.ID)
	require.NoError(t, err)
	assert.Equal(t, "confirmada", actualizada.Estado)
}

func TestReservaRepo_UpdateEstado_InexistenteRetornaError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReservaRepo(db)

	err := repo.UpdateEstado(999999, "confirmada")

	assert.Error(t, err)
}

func TestReservaRepo_VincularPedido_ActualizaCorrectamente(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReservaRepo(db)

	reserva := &models.Reserva{Nombre: "Marta", Whatsapp: "4", Fecha: time.Now(), Personas: "2", Estado: "confirmada"}
	require.NoError(t, repo.Create(reserva))

	user := &models.User{Nombre: "Emp", Apellido: "Leado", Email: "emp-vincular@test.com", Password: "x", Rol: "empleado"}
	require.NoError(t, db.Create(user).Error)
	cliente := &models.Cliente{Nombre: "Cliente Test", Tipo: "individual"}
	require.NoError(t, db.Create(cliente).Error)
	pedido := &models.Pedido{ClienteID: cliente.ID, UserID: user.ID, Fecha: time.Now(), TipoComida: "almuerzo", FormaPago: "contado", Estado: "pendiente"}
	require.NoError(t, db.Create(pedido).Error)

	err := repo.VincularPedido(reserva.ID, pedido.ID)
	require.NoError(t, err)

	actualizada, err := repo.FindByID(reserva.ID)
	require.NoError(t, err)
	require.NotNil(t, actualizada.PedidoID)
	assert.Equal(t, pedido.ID, *actualizada.PedidoID)
}

func TestReservaRepo_VincularPedido_InexistenteRetornaError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReservaRepo(db)

	err := repo.VincularPedido(999999, 1)

	assert.Error(t, err)
}
