package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"restaurantos/internal/models"
)

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("error conectando a la base de datos: %v", err)
	}
	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Cliente{},
		&models.Producto{},
		&models.Pedido{},
		&models.PedidoItem{},
		&models.Pago{},
		&models.MovimientoStock{},
		&models.MenuPublico{},
		&models.Reserva{},
	)
	if err != nil {
		log.Fatalf("error en migraciones: %v", err)
	}
}
