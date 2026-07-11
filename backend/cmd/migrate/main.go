package main

import (
    "log"
    "os"
    "restaurantos/internal/config"
    "restaurantos/internal/database"
)

func main() {
    if os.Getenv("DATABASE_URL") == "" {
        os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/restaurantos?sslmode=disable")
    }
    cfg := config.Load()
    db := database.Connect(cfg.DatabaseURL)
    database.Migrate(db)
    log.Println("Migraciones completadas OK")
}
