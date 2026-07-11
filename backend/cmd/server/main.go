package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"restaurantos/internal/config"
	"restaurantos/internal/database"
	"restaurantos/internal/handlers"
	"restaurantos/internal/middleware"
	"restaurantos/internal/repository"
	"restaurantos/internal/services"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DatabaseURL)
	database.Migrate(db)

	r := gin.Default()

	// ─── Rutas públicas ────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authH := handlers.NewAuthHandler(repository.NewUserRepo(db), cfg.JWTSecret)
	r.POST("/api/auth/login", authH.Login)

	menuRepo := repository.NewMenuRepo(db)
	menuH := handlers.NewMenuHandler(menuRepo)
	r.GET("/api/public/menu", menuH.Publico)

	reservaRepo := repository.NewReservaRepo(db)

	// ─── Rutas protegidas (cualquier rol autenticado) ─────────────────────────
	protected := r.Group("/api", middleware.JWTMiddleware(cfg.JWTSecret))
	{
		protected.GET("/auth/me", authH.Me)
	}

	// ─── Rutas empleado (autenticado, cualquier rol) ───────────────────────────
	clienteRepo := repository.NewClienteRepo(db)
	clienteH := handlers.NewClienteHandler(clienteRepo)

	empleado := r.Group("/api/empleado", middleware.JWTMiddleware(cfg.JWTSecret))
	{
		empleado.GET("/clientes", clienteH.Listar)
		empleado.GET("/clientes/:id", clienteH.ObtenerPorID)
	}

	// ─── Pedidos ───────────────────────────────────────────────────────────────
	pedidoRepo := repository.NewPedidoRepo(db)
	pedidoSvc := services.NewPedidoService(
		pedidoRepo,
		repository.NewProductoRepo(db),
		clienteRepo,
	)
	pedidoH := handlers.NewPedidoHandler(pedidoSvc)

	empleado.POST("/pedidos", pedidoH.Crear)
	empleado.GET("/pedidos", pedidoH.Listar)
	empleado.GET("/pedidos/:id", pedidoH.ObtenerPorID)
	empleado.PUT("/pedidos/:id/entregar", pedidoH.MarcarEntregado)

	// ─── Inventario ────────────────────────────────────────────────────────────
	productoRepo := repository.NewProductoRepo(db)
	movimientoRepo := repository.NewMovimientoRepo(db)
	inventarioH := handlers.NewInventarioHandler(productoRepo, movimientoRepo)
	empleado.GET("/productos", inventarioH.Listar)

	// ─── Reservas ──────────────────────────────────────────────────────────────
	reservaH := handlers.NewReservaHandler(reservaRepo, pedidoRepo)
	r.POST("/api/public/reservas", reservaH.Crear)
	empleado.GET("/reservas", reservaH.Listar)
	empleado.PUT("/reservas/:id/estado", reservaH.ActualizarEstado)
	empleado.PUT("/reservas/:id/pedido", reservaH.VincularPedido)

	// ─── Rutas solo admin ──────────────────────────────────────────────────────
	admin := r.Group("/api/admin", middleware.JWTMiddleware(cfg.JWTSecret), middleware.RequireRole("admin"))
	{
		admin.POST("/clientes", clienteH.Crear)
		admin.PUT("/clientes/:id", clienteH.Actualizar)
		admin.DELETE("/clientes/:id", clienteH.Desactivar)
		admin.DELETE("/pedidos/:id", pedidoH.Anular)

		// Inventario admin
		admin.GET("/productos/alertas", inventarioH.Alertas)
		admin.GET("/productos/:id", inventarioH.ObtenerPorID)
		admin.POST("/productos", inventarioH.Crear)
		admin.PUT("/productos/:id", inventarioH.Actualizar)
		admin.POST("/productos/:id/restock", inventarioH.Restock)

		// Créditos y pagos
		creditosH := handlers.NewCreditosHandler(repository.NewPagoRepo(db), clienteRepo)
		admin.GET("/creditos", creditosH.ListarDeudores)
		admin.GET("/creditos/:cliente_id", creditosH.EstadoCuenta)
		admin.POST("/pagos", creditosH.RegistrarPago)

		// Menú público (gestión desde el dashboard)
		admin.GET("/menu", menuH.Listar)
		admin.POST("/menu", menuH.Crear)
		admin.PUT("/menu/:id", menuH.Actualizar)
		admin.DELETE("/menu/:id", menuH.Eliminar)

		// Reportes
		reportesH := handlers.NewReportesHandler(pedidoRepo, movimientoRepo, clienteRepo)
		admin.GET("/reportes/ventas", reportesH.Ventas)
		admin.GET("/reportes/deudores", reportesH.Deudores)
		admin.GET("/reportes/inventario", reportesH.Inventario)
	}

	log.Printf("Servidor iniciando en puerto %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
