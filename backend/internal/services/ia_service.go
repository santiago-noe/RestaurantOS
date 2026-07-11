package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"restaurantos/internal/models"
)

var tiposComida = []string{"desayuno", "almuerzo", "cena", "bebida"}

type Prediccion struct {
	TipoComida          string `json:"tipo_comida"`
	PorcionesEstimadas  int    `json:"porciones_estimadas"`
	BaseCalculo         string `json:"base_calculo"`
}

type Alerta struct {
	Tipo      string `json:"tipo"`
	Mensaje   string `json:"mensaje"`
	Severidad string `json:"severidad"`
	RefID     int    `json:"ref_id,omitempty"`
}

// PredecirDemanda calcula el promedio de pedidos por tipo de comida
// para el mismo día de la semana que `fechaObjetivo`, basado en historial.
func PredecirDemanda(pedidos []models.Pedido, fechaObjetivo time.Time) []Prediccion {
	diaSemana := fechaObjetivo.Weekday()

	// contar pedidos por tipo agrupados por fecha (para calcular promedio por día)
	type clave struct {
		fecha      string
		tipoComida string
	}
	contadorPorDia := map[clave]int{}
	diasVistos := map[string]bool{}

	for _, p := range pedidos {
		if p.Fecha.Weekday() != diaSemana {
			continue
		}
		fechaStr := p.Fecha.Format("2006-01-02")
		diasVistos[fechaStr] = true
		contadorPorDia[clave{fechaStr, p.TipoComida}]++
	}

	totalDias := len(diasVistos)

	predicciones := make([]Prediccion, 0, len(tiposComida))
	for _, tipo := range tiposComida {
		var suma int
		for k, v := range contadorPorDia {
			if k.tipoComida == tipo {
				suma += v
			}
		}

		estimado := 0
		if totalDias > 0 {
			estimado = int(math.Ceil(float64(suma) / float64(totalDias)))
		}

		predicciones = append(predicciones, Prediccion{
			TipoComida:         tipo,
			PorcionesEstimadas: estimado,
			BaseCalculo:        fmt.Sprintf("promedio de los últimos %d %s", totalDias, nombreDia(diaSemana)),
		})
	}
	return predicciones
}

// GenerarAlertas produce alertas de deuda alta y stock bajo.
func GenerarAlertas(clientes []models.Cliente, productos []models.Producto, umbralDeuda float64) []Alerta {
	var alertas []Alerta

	for _, cl := range clientes {
		if cl.DeudaTotal > umbralDeuda {
			alertas = append(alertas, Alerta{
				Tipo:      "deuda_alta",
				Mensaje:   fmt.Sprintf("%s %s tiene una deuda de S/ %.2f (umbral: S/ %.2f)", cl.Nombre, cl.Apellido, cl.DeudaTotal, umbralDeuda),
				Severidad: "alta",
				RefID:     cl.ID,
			})
		}
	}

	for _, p := range productos {
		if p.StockActual < p.StockMinimo {
			alertas = append(alertas, Alerta{
				Tipo:      "stock_bajo",
				Mensaje:   fmt.Sprintf("%s: %.1f %s restantes (mínimo: %.1f)", p.Nombre, p.StockActual, p.Unidad, p.StockMinimo),
				Severidad: "media",
				RefID:     p.ID,
			})
		}
	}

	// Ordenar: alta > media
	sort.Slice(alertas, func(i, j int) bool {
		return severidadOrden(alertas[i].Severidad) > severidadOrden(alertas[j].Severidad)
	})

	return alertas
}

func severidadOrden(s string) int {
	switch s {
	case "alta":
		return 2
	case "media":
		return 1
	default:
		return 0
	}
}

func nombreDia(d time.Weekday) string {
	nombres := map[time.Weekday]string{
		time.Sunday: "domingos", time.Monday: "lunes", time.Tuesday: "martes",
		time.Wednesday: "miércoles", time.Thursday: "jueves",
		time.Friday: "viernes", time.Saturday: "sábados",
	}
	return nombres[d]
}
