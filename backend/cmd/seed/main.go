// Genera datos históricos simulados (clientes, pedidos, inventario, créditos)
// desde el 2025-11-01 hasta la fecha actual, usando los mismos repositorios y
// servicios de la app para mantener la lógica de negocio consistente.
package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"restaurantos/internal/config"
	"restaurantos/internal/database"
	"restaurantos/internal/models"
	"restaurantos/internal/repository"
	"restaurantos/internal/services"
)

type nuevoCliente struct {
	nombre, apellido, tipo, telefono string
}

var clientesNuevos = []nuevoCliente{
	{"Rosa", "Quispe Mamani", "individual", "987002001"},
	{"Miguel", "Flores Rojas", "individual", "987002002"},
	{"Carmen", "Vargas Huamán", "individual", "987002003"},
	{"Luis", "Mamani Ccora", "individual", "987002004"},
	{"Elena", "Torres Quispe", "individual", "987002005"},
	{"Jorge", "Huamán Rivera", "individual", "987002006"},
	{"Lucía", "Rojas Palomino", "individual", "987002007"},
	{"Minera Vilcas", "EIRL", "empresa", "01-345678"},
	{"Transportes Ayacucho", "SAC", "empresa", "01-456789"},
	{"Hotel Wari", "SAC", "empresa", "01-567890"},
}

var tiposComida = []string{"desayuno", "almuerzo", "cena"}
var metodosPago = []string{"efectivo", "yape", "transferencia"}

func pesoTipoComida() string {
	x := rand.Float64()
	switch {
	case x < 0.20:
		return "desayuno"
	case x < 0.75:
		return "almuerzo"
	default:
		return "cena"
	}
}

func cantidadPara(unidad string) float64 {
	if unidad == "unidad" {
		return float64(1 + rand.Intn(5))
	}
	return float64(1+rand.Intn(5)) * 0.5 // 0.5 a 3.0
}


func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DatabaseURL)
	database.Migrate(db)

	clienteRepo := repository.NewClienteRepo(db)
	productoRepo := repository.NewProductoRepo(db)
	pedidoRepo := repository.NewPedidoRepo(db)
	movimientoRepo := repository.NewMovimientoRepo(db)
	pagoRepo := repository.NewPagoRepo(db)
	pedidoSvc := services.NewPedidoService(pedidoRepo, productoRepo, clienteRepo)

	// ─── 1. Asegurar usuarios base ────────────────────────────────────────────
	var userIDs []int
	var users []models.User
	db.Find(&users)
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	if len(userIDs) == 0 {
		log.Fatal("no hay usuarios en la BD; corre el seed_dev.sql primero para crear el admin/empleado")
	}

	// ─── 2. Agregar clientes nuevos (si no existen ya por nombre) ─────────────
	var clientesExistentes []models.Cliente
	db.Find(&clientesExistentes)
	nombresExistentes := map[string]bool{}
	for _, c := range clientesExistentes {
		nombresExistentes[c.Nombre+"|"+c.Apellido] = true
	}

	for _, nc := range clientesNuevos {
		if nombresExistentes[nc.nombre+"|"+nc.apellido] {
			continue
		}
		cliente := &models.Cliente{
			Nombre: nc.nombre, Apellido: nc.apellido, Tipo: nc.tipo,
			Telefono: nc.telefono, Activo: true,
		}
		if err := clienteRepo.Create(cliente); err != nil {
			log.Fatalf("error creando cliente %s: %v", nc.nombre, err)
		}
	}

	var todosClientes []models.Cliente
	db.Where("activo = true").Find(&todosClientes)
	var clienteIDs []int
	deudaAcumulada := map[int]float64{}
	for _, c := range todosClientes {
		clienteIDs = append(clienteIDs, c.ID)
		deudaAcumulada[c.ID] = c.DeudaTotal
	}

	// ─── 3. Productos existentes ───────────────────────────────────────────────
	var productos []models.Producto
	db.Where("activo = true").Find(&productos)
	if len(productos) == 0 {
		log.Fatal("no hay productos en la BD; corre el seed_dev.sql primero")
	}

	// Stock en memoria para decidir cuándo reabastecer (punto de reorden),
	// evitando acumular stock indefinidamente semana tras semana.
	stockActual := map[int]float64{}
	stockMinimo := map[int]float64{}
	for _, p := range productos {
		stockActual[p.ID] = p.StockActual
		stockMinimo[p.ID] = p.StockMinimo
	}

	// ─── 4. Simulación día a día ────────────────────────────────────────────────
	inicio := time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local)
	hoy := time.Now()
	hoyTruncado := time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 0, 0, 0, 0, time.Local)

	totalPedidos, totalAnulados, totalPagos, totalRestocks := 0, 0, 0, 0

	for d := inicio; !d.After(hoyTruncado); d = d.AddDate(0, 0, 1) {
		// Reabastecimiento semanal (lunes), solo si el stock está por debajo
		// del punto de reorden — así no se acumula indefinidamente.
		if d.Weekday() == time.Monday {
			for _, p := range productos {
				puntoReorden := stockMinimo[p.ID] * 3
				objetivo := stockMinimo[p.ID] * 5
				if stockActual[p.ID] >= puntoReorden {
					continue
				}
				qty := round2(objetivo - stockActual[p.ID])
				if qty <= 0 {
					continue
				}
				if err := productoRepo.AjustarStock(p.ID, qty); err != nil {
					continue
				}
				stockActual[p.ID] += qty
				movimientoRepo.Registrar(&models.MovimientoStock{
					ProductoID: p.ID, Tipo: "entrada", Cantidad: qty,
					Notas: "Reabastecimiento semanal", Fecha: d,
				})
				totalRestocks++
			}
		}

		// Volumen de pedidos según día de la semana y temporada
		base := 3
		switch d.Weekday() {
		case time.Friday, time.Saturday:
			base = 6
		case time.Sunday:
			base = 5
		}
		if d.Month() == time.December && d.Day() >= 15 {
			base = int(float64(base) * 1.6)
		}
		ordersHoy := base + rand.Intn(3)

		diasAntiguedad := hoyTruncado.Sub(d).Hours() / 24

		for i := 0; i < ordersHoy; i++ {
			clienteID := clienteIDs[rand.Intn(len(clienteIDs))]
			userID := userIDs[rand.Intn(len(userIDs))]

			nItems := 1 + rand.Intn(4)
			items := make([]services.ItemInput, 0, nItems)
			usados := map[int]bool{}
			for len(items) < nItems {
				p := productos[rand.Intn(len(productos))]
				if usados[p.ID] {
					continue
				}
				usados[p.ID] = true
				items = append(items, services.ItemInput{
					ProductoID: p.ID, Cantidad: cantidadPara(p.Unidad), PrecioUnitario: p.PrecioVenta,
				})
			}

			hora := 8 + rand.Intn(13)
			fecha := time.Date(d.Year(), d.Month(), d.Day(), hora, rand.Intn(60), 0, 0, time.Local)

			formaPago := "contado"
			if rand.Float64() < 0.35 {
				formaPago = "credito"
			}

			pedido, err := pedidoSvc.Crear(services.CrearPedidoInput{
				ClienteID: clienteID, UserID: userID, Fecha: fecha,
				TipoComida: pesoTipoComida(), FormaPago: formaPago, Items: items,
			})
			if err != nil {
				continue // stock insuficiente u otro; se omite este pedido
			}
			totalPedidos++

			// Movimiento de salida por el consumo de stock del pedido
			for _, it := range items {
				stockActual[it.ProductoID] -= it.Cantidad
				movimientoRepo.Registrar(&models.MovimientoStock{
					ProductoID: it.ProductoID, Tipo: "salida", Cantidad: it.Cantidad,
					Notas: "Consumo por pedido", Fecha: fecha,
				})
			}

			// Estado final según antigüedad
			var estado string
			x := rand.Float64()
			if diasAntiguedad <= 1 {
				switch {
				case x < 0.55:
					estado = "pendiente"
				case x < 0.95:
					estado = "entregado"
				default:
					estado = "anulado"
				}
			} else {
				if x < 0.93 {
					estado = "entregado"
				} else {
					estado = "anulado"
				}
			}

			switch estado {
			case "entregado":
				pedidoRepo.UpdateEstado(pedido.ID, "entregado")
			case "anulado":
				if err := pedidoSvc.Anular(pedido.ID); err == nil {
					for _, it := range items {
						stockActual[it.ProductoID] += it.Cantidad
						movimientoRepo.Registrar(&models.MovimientoStock{
							ProductoID: it.ProductoID, Tipo: "entrada", Cantidad: it.Cantidad,
							Notas: "Devolución por pedido anulado", Fecha: fecha,
						})
					}
					totalAnulados++
				}
			}

			if formaPago == "credito" && estado != "anulado" {
				deudaAcumulada[clienteID] += pedido.Total
			}
		}

		// Pagos de crédito los domingos
		if d.Weekday() == time.Sunday {
			for _, clienteID := range clienteIDs {
				deuda := deudaAcumulada[clienteID]
				if deuda <= 0 || rand.Float64() > 0.4 {
					continue
				}
				porcentaje := 0.4 + rand.Float64()*0.6
				monto := round2(deuda * porcentaje)
				if monto <= 0 {
					continue
				}
				if err := pagoRepo.Create(&models.Pago{
					ClienteID: clienteID, Monto: monto,
					Metodo: metodosPago[rand.Intn(len(metodosPago))], Fecha: d,
				}); err == nil {
					deudaAcumulada[clienteID] = math.Max(0, round2(deuda-monto))
					totalPagos++
				}
			}
		}
	}

	// ─── 5. Persistir deuda_total final por cliente ───────────────────────────
	for clienteID, deuda := range deudaAcumulada {
		clienteRepo.Update(clienteID, map[string]interface{}{"deuda_total": round2(deuda)})
	}

	log.Printf("Listo. Pedidos creados: %d (anulados: %d) | Pagos: %d | Restocks: %d | Clientes: %d",
		totalPedidos, totalAnulados, totalPagos, totalRestocks, len(clienteIDs))
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
