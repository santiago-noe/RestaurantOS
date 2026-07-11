# FC-05 — Inventario
**Estado:** Especificado  
**Sprint:** 3  
**Prioridad:** Media-Alta  
**Depende de:** FC-01

---

## Descripción
Control de stock de ingredientes y productos. El stock se descuenta automáticamente al registrar pedidos. Se alerta cuando el stock baja del mínimo.

---

## Endpoints

| Método | Ruta                          | Rol   | Descripción                         |
|--------|-------------------------------|-------|-------------------------------------|
| GET    | /api/empleado/productos       | any   | Listar productos con stock          |
| GET    | /api/admin/productos/:id      | admin | Ver producto + historial movimientos|
| POST   | /api/admin/productos          | admin | Crear producto                      |
| PUT    | /api/admin/productos/:id      | admin | Editar producto                     |
| POST   | /api/admin/productos/:id/restock | admin | Reabastecimiento (entrada)       |
| GET    | /api/admin/productos/alertas  | admin | Productos bajo stock mínimo         |

---

## Contrato POST /api/admin/productos/:id/restock

**Request:**
```json
{
  "cantidad": 10.5,
  "notas": "Compra del mercado central"
}
```

**Response 200:**
```json
{
  "producto_id": 3,
  "nombre": "Arroz",
  "stock_anterior": 2.5,
  "cantidad_agregada": 10.5,
  "stock_actual": 13.0,
  "unidad": "kg"
}
```

---

## Contrato GET /api/admin/productos/alertas

```json
{
  "alertas": [
    {
      "id": 5,
      "nombre": "Aceite",
      "stock_actual": 0.5,
      "stock_minimo": 2.0,
      "unidad": "litro",
      "deficit": 1.5
    }
  ]
}
```

---

## Reglas de negocio
1. `stock_actual` nunca puede ser negativo. Si un pedido agota el stock → error 422.
2. Todo movimiento (entrada/salida) se registra en `movimientos_stock`.
3. La "salida" automática al crear un pedido se registra con `referencia_id = pedido.id`.
4. El restock registra una "entrada" con nota opcional.

---

## Pruebas Unitarias

```
TestCrearProducto_DatosValidosRetorna201
TestCrearProducto_UnidadInvalidaRetorna400
TestRestock_AumentaStockCorrectamente
TestRestock_RegistraMovimientoEntrada
TestDescontarStock_ReduceStockCorrectamente
TestDescontarStock_StockInsuficienteRetornaError
TestDescontarStock_RegistraMovimientoSalida
TestObtenerAlertas_ProductoBajoMinimoAparece
TestObtenerAlertas_ProductoSobreMinimoNoAparece
TestHistorialMovimientos_OrdenadoPorFecha
```

---

## Definición de "Hecho" (DoD)
- [ ] Todos los tests pasan.
- [ ] Cobertura ≥ 90%.
- [ ] No es posible que el stock quede negativo.
- [ ] Cada movimiento queda registrado con fecha y referencia.
- [ ] Las alertas de stock bajo son visibles en el dashboard.
