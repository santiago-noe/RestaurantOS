# FC-07 — IA Básica (Predicción y Alertas)
**Estado:** Especificado  
**Sprint:** 3  
**Prioridad:** Media  
**Depende de:** FC-03 (historial de pedidos)

---

## Descripción
Módulo de inteligencia básica sin modelos externos. Usa promedios históricos para predecir demanda y detectar patrones de alerta. No requiere ML externo — es estadística pura.

---

## Endpoints

| Método | Ruta                      | Rol   | Descripción                              |
|--------|---------------------------|-------|------------------------------------------|
| GET    | /api/admin/ia/prediccion  | admin | Predicción de demanda para mañana        |
| GET    | /api/admin/ia/alertas     | admin | Alertas activas del sistema              |

---

## Contrato GET /api/admin/ia/prediccion

```json
{
  "fecha_prediccion": "2026-05-30",
  "dia_semana": "Sábado",
  "predicciones": [
    {
      "tipo_comida": "almuerzo",
      "porciones_estimadas": 48,
      "base_calculo": "promedio de los últimos 8 sábados"
    },
    {
      "tipo_comida": "desayuno",
      "porciones_estimadas": 22,
      "base_calculo": "promedio de los últimos 8 sábados"
    }
  ],
  "ingredientes_criticos": [
    {
      "nombre": "Pollo",
      "cantidad_necesaria": 12.0,
      "stock_actual": 3.5,
      "deficit": 8.5,
      "unidad": "kg"
    }
  ]
}
```

## Contrato GET /api/admin/ia/alertas

```json
{
  "alertas": [
    {
      "tipo": "deuda_alta",
      "mensaje": "Juan Pérez tiene una deuda de S/ 245.00 (umbral: S/ 200)",
      "cliente_id": 5,
      "severidad": "alta"
    },
    {
      "tipo": "stock_bajo",
      "mensaje": "Aceite: 0.5 litros restantes (mínimo: 2 litros)",
      "producto_id": 7,
      "severidad": "media"
    }
  ]
}
```

---

## Algoritmo de predicción

```
Para el día D de la semana:
  1. Obtener todos los pedidos de los últimos 30 días donde day_of_week(fecha) = D
  2. Agrupar por tipo_comida
  3. Para cada tipo: prediccion = promedio(cantidad_pedidos_ese_dia)
  4. Redondear hacia arriba (siempre es mejor tener de más)
  5. Cruzar prediccion con stock actual para identificar déficits
```

---

## Pruebas Unitarias

```
TestPredecirDemanda_SinHistorialRetornaPrediccionCero
TestPredecirDemanda_HistorialCompletoPredicePromedio
TestPredecirDemanda_RedondeoHaciaArriba
TestPredecirDemanda_SoloUsaDiaSemanaCorrespondiente
TestCalcularDeficitIngredientes_StockSuficienteNoAlert
TestCalcularDeficitIngredientes_StockInsuficienteAlerta
TestGenerarAlertas_DeudaAltaIncluida
TestGenerarAlertas_DeudaBajaNoIncluida
TestGenerarAlertas_StockBajoIncluido
TestGenerarAlertas_OrdenadaPorSeveridad
```

---

## Definición de "Hecho" (DoD)
- [ ] Todos los tests pasan.
- [ ] Cobertura ≥ 90%.
- [ ] Sin datos históricos → predicción retorna 0, no error 500.
- [ ] La predicción solo usa el día de la semana correcto.
- [ ] Las alertas aparecen en el dashboard en tiempo real.
