package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/models"
)

func TestClienteRepo_Create_GuardaClienteEnBD(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClienteRepo(db)

	cliente := &models.Cliente{Nombre: "Juan Perez", Tipo: "individual", Email: "juan@test.com"}

	err := repo.Create(cliente)

	require.NoError(t, err)
	assert.NotZero(t, cliente.ID)
}

func TestClienteRepo_FindByID_ClienteExisteRetornaDatos(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClienteRepo(db)

	cliente := &models.Cliente{Nombre: "Ana Lopez", Tipo: "individual"}
	require.NoError(t, repo.Create(cliente))

	encontrado, err := repo.FindByID(cliente.ID)

	require.NoError(t, err)
	require.NotNil(t, encontrado)
	assert.Equal(t, "Ana Lopez", encontrado.Nombre)
}

func TestClienteRepo_FindByID_InexistenteRetornaNil(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClienteRepo(db)

	encontrado, err := repo.FindByID(999999)

	require.NoError(t, err)
	assert.Nil(t, encontrado)
}

func TestClienteRepo_FindAll_RetornaSoloActivosPaginados(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClienteRepo(db)

	require.NoError(t, repo.Create(&models.Cliente{Nombre: "Cliente A", Tipo: "individual"}))
	require.NoError(t, repo.Create(&models.Cliente{Nombre: "Cliente B", Tipo: "empresa"}))

	clientes, total, err := repo.FindAll(1, 10, "")

	require.NoError(t, err)
	assert.EqualValues(t, 2, total)
	assert.Len(t, clientes, 2)
}

func TestClienteRepo_Update_ActualizaCamposCorrectamente(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClienteRepo(db)

	cliente := &models.Cliente{Nombre: "Pedro", Tipo: "individual"}
	require.NoError(t, repo.Create(cliente))

	actualizado, err := repo.Update(cliente.ID, map[string]interface{}{"nombre": "Pedro Editado"})

	require.NoError(t, err)
	require.NotNil(t, actualizado)
	assert.Equal(t, "Pedro Editado", actualizado.Nombre)
}

func TestClienteRepo_Deactivate_HaceSoftDeleteYaNoAparece(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClienteRepo(db)

	cliente := &models.Cliente{Nombre: "Carlos", Tipo: "individual"}
	require.NoError(t, repo.Create(cliente))

	err := repo.Deactivate(cliente.ID)
	require.NoError(t, err)

	encontrado, err := repo.FindByID(cliente.ID)
	require.NoError(t, err)
	assert.Nil(t, encontrado, "un cliente desactivado no debe encontrarse por FindByID")
}

func TestClienteRepo_EmailExists_EmailYaRegistradoRetornaTrue(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClienteRepo(db)

	require.NoError(t, repo.Create(&models.Cliente{Nombre: "Luis", Tipo: "individual", Email: "luis@test.com"}))

	assert.True(t, repo.EmailExists("luis@test.com", 0))
}

func TestClienteRepo_EmailExists_EmailNoRegistradoRetornaFalse(t *testing.T) {
	db := setupTestDB(t)
	repo := NewClienteRepo(db)

	assert.False(t, repo.EmailExists("noexiste@test.com", 0))
}
