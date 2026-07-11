# FC-03 — Registro de Pedidos
**Estado:** Especificado  
**Sprint:** 2  
**Prioridad:** Alta  
**Depende de:** FC-01, FC-02, FC-05 (productos deben existir)

---

## Descripción
Módulo central del sistema. Registra lo que consume cada cliente, calcula el total, descuenta stock y —si es a crédito— genera una deuda.

---

## Endpoints

| Método | Ruta                    | Rol   | Descripción                  |
|--------|-------------------------|-------|------------------------------|
| GET    | /api/empleado/pedidos   | any   | Listar pedidos (con filtros) |
| GET    | /api/empleado/pedidos/:id | any | Ver detalle del pedido       |
| POST   | /api/empleado/pedidos   | any   | Registrar nuevo pedido       |
| PUT    | /api/admin/pedidos/:id/estado | admin | Cambiar estado        |
| DELETE | /api/admin/pedidos/:id  | admin | Anular pedido                |

---

## Contrato POST /api/empleado/pedidos

**Request:**
```json
{
  "cliente_id": 5,
  "fecha": "2026-05-29",
  "tipo_comida": "almuerzo",
  "forma_pago": "credito",
  "notas": "Sin ají",
  "items": [
    {
      "producto_id": 3,
      "cantidad": 2,
      "precio_unitario": 12.50
    },
    {
      "producto_id": 7,
      "cantidad": 1,
      "precio_unitario": 3.00
    }
  ]
}
```

**Response 201 Created:**
```json
{
  "id": 42,
  "cliente_id": 5,
  "total": 28.00,
  "estado": "pendiente",
  "forma_pago": "credito",
  "items": [...]
}
```

---

## Lógica de negocio (service layer)

1. Validar que el cliente existe y está activo.
2. Validar que cada `producto_id` existe.
3. Calcular `subtotal = cantidad × precio_unitario` por ítem.
4. Calcular `total = SUM(subtotales)`.
5. Insertar pedido + pedido_items en una **transacción**.
6. Descontar stock: `UPDATE productos SET stock_actual = stock_actual - cantidad`.
7. Si `forma_pago = 'credito'` → el trigger de BD actualiza `clientes.deuda_total`.
8. Si stock de algún producto queda bajo el mínimo → registrar alerta.

---

## Pruebas Unitarias

```
TestCalcularTotal_ItemsMultiplesRetornaSumaCorrecta
TestCalcularTotal_ListaVaciaRetornaCero
TestCrearPedido_DatosValidosRetorna201
TestCrearPedido_ClienteInexistenteRetorna404
TestCrearPedido_ProductoInexistenteRetorna404
TestCrearPedido_StockInsuficienteRetorna422
TestCrearPedido_CreditoActualizaDeudaCliente
TestCrearPedido_ContadoNoModificaDeuda
TestCrearPedido_DescuentaStockCorrectamente
TestAnularPedido_DevuelveStockAlInventario
TestAnularPedido_CreditoReduceDeudaCliente
TestAnularPedido_EstadoAnuladoNoPermiteOtraAnulacion
```

---

## Definición de "Hecho" (DoD)
- [ ] Todos los tests pasan.
- [ ] Cobertura ≥ 90%.
- [ ] Pedido a crédito actualiza `deuda_total` del cliente.
- [ ] Anular un pedido devuelve stock y reduce deuda.
- [ ] Todo ocurre en una transacción (si falla algo, no queda nada a medias).
