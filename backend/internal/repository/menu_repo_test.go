package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/models"
)

func TestMenuRepo_Create_GuardaItemEnBD(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMenuRepo(db)

	item := &models.MenuPublico{Categoria: "Entradas", Nombre: "Puca Picante", Precio: 64.0, Disponible: true}

	err := repo.Create(item)

	require.NoError(t, err)
	assert.NotZero(t, item.ID)
}

func TestMenuRepo_FindPublico_SoloRetornaDisponibles(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMenuRepo(db)

	require.NoError(t, repo.Create(&models.MenuPublico{Categoria: "Fondos", Nombre: "Cuy Chactado", Disponible: true}))
	require.NoError(t, repo.Create(&models.MenuPublico{Categoria: "Fondos", Nombre: "Plato Agotado", Disponible: false}))

	items, err := repo.FindPublico()

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "Cuy Chactado", items[0].Nombre)
}

func TestMenuRepo_FindPublico_OrdenaPorCampoOrden(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMenuRepo(db)

	require.NoError(t, repo.Create(&models.MenuPublico{Nombre: "Segundo", Disponible: true, Orden: 2}))
	require.NoError(t, repo.Create(&models.MenuPublico{Nombre: "Primero", Disponible: true, Orden: 1}))

	items, err := repo.FindPublico()

	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, "Primero", items[0].Nombre)
	assert.Equal(t, "Segundo", items[1].Nombre)
}

func TestMenuRepo_FindAll_RetornaDisponiblesYNoDisponibles(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMenuRepo(db)

	require.NoError(t, repo.Create(&models.MenuPublico{Nombre: "A", Disponible: true}))
	require.NoError(t, repo.Create(&models.MenuPublico{Nombre: "B", Disponible: false}))

	items, err := repo.FindAll()

	require.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestMenuRepo_FindByID_InexistenteRetornaNil(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMenuRepo(db)

	item, err := repo.FindByID(999999)

	require.NoError(t, err)
	assert.Nil(t, item)
}

func TestMenuRepo_Update_ActualizaCamposCorrectamente(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMenuRepo(db)

	item := &models.MenuPublico{Nombre: "Original", Precio: 10, Disponible: true}
	require.NoError(t, repo.Create(item))

	actualizado, err := repo.Update(item.ID, map[string]interface{}{"nombre": "Editado", "precio": 20.5})

	require.NoError(t, err)
	require.NotNil(t, actualizado)
	assert.Equal(t, "Editado", actualizado.Nombre)
	assert.Equal(t, 20.5, actualizado.Precio)
}

func TestMenuRepo_Delete_EliminaElItem(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMenuRepo(db)

	item := &models.MenuPublico{Nombre: "A eliminar", Disponible: true}
	require.NoError(t, repo.Create(item))

	err := repo.Delete(item.ID)
	require.NoError(t, err)

	encontrado, err := repo.FindByID(item.ID)
	require.NoError(t, err)
	assert.Nil(t, encontrado)
}

func TestMenuRepo_Delete_IDInexistenteRetornaError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMenuRepo(db)

	err := repo.Delete(999999)

	assert.Error(t, err)
}
