package repository

import (
	"os"
	"sync"
	"testing"

	"gorm.io/gorm"

	"restaurantos/internal/database"
)

var (
	testDB     *gorm.DB
	testDBOnce sync.Once
)

// setupTestDB conecta una sola vez a la BD de test (postgres_test, puerto 5433)
// y devuelve una transacción nueva por test que se revierte automáticamente
// al terminar, para que ningún test deje datos residuales para el siguiente.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	testDBOnce.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dsn = "postgres://postgres:santiago09@localhost:5433/restaurantos_test?sslmode=disable"
		}
		testDB = database.Connect(dsn)
		database.Migrate(testDB)
	})

	tx := testDB.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}
