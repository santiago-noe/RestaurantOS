# FC-04 — Créditos y Pagos
**Estado:** Especificado  
**Sprint:** 2  
**Prioridad:** Alta  
**Depende de:** FC-02, FC-03

---

## Descripción
Gestión financiera del restaurante. Permite registrar abonos de clientes con deuda, ver el estado de cuenta detallado y alertar cuando una deuda supera el umbral.

---

## Endpoints

| Método | Ruta                              | Rol   | Descripción                        |
|--------|-----------------------------------|-------|------------------------------------|
| GET    | /api/admin/creditos               | admin | Listado de clientes con deuda      |
| GET    | /api/admin/creditos/:cliente_id   | admin | Estado de cuenta detallado         |
| POST   | /api/admin/pagos                  | admin | Registrar abono                    |
| GET    | /api/admin/pagos                  | admin | Historial de pagos (con filtros)   |

---

## Contrato GET /api/admin/creditos/:cliente_id

```json
{
  "cliente": {
    "id": 5,
    "nombre": "Juan Pérez",
    "telefono": "987654321"
  },
  "deuda_total": 120.00,
  "pedidos_pendientes": [
    {
      "id": 42,
      "fecha": "2026-05-20",
      "total": 28.00,
      "abonado": 0.00,
      "saldo": 28.00
    },
    {
      "id": 38,
      "fecha": "2026-05-18",
      "total": 92.00,
      "abonado": 0.00,
      "saldo": 92.00
    }
  ],
  "pagos_realizados": [
    {
      "id": 10,
      "fecha": "2026-05-25",
      "monto": 50.00,
      "metodo": "yape"
    }
  ]
}
```

## Contrato POST /api/admin/pagos

**Request:**
```json
{
  "cliente_id": 5,
  "monto": 50.00,
  "metodo": "yape",
  "fecha": "2026-05-29",
  "pedido_id": 42
}
```

---

## Reglas de negocio
1. Un pago puede asociarse a un pedido específico (`pedido_id`) o ser un abono general (`pedido_id = null`).
2. El monto del pago no puede ser mayor a la deuda total del cliente.
3. Si la deuda supera S/ 200 (configurable), el sistema genera una alerta visible en el dashboard.
4. Al registrar un pago, el trigger actualiza `clientes.deuda_total` automáticamente.

---

## Pruebas Unitarias

```
TestRegistrarPago_MontoValidoActualizaDeuda
TestRegistrarPago_MontoMayorADeudaRetorna422
TestRegistrarPago_ClienteSinDeudaRetorna422
TestRegistrarPago_PedidoIDInexistenteRetorna404
TestObtenerEstadoCuenta_MuestraPedidosPendientes
TestObtenerEstadoCuenta_MuestraPagosRealizados
TestListarClientesConDeuda_FiltroDeudaAlta
TestAlertaDeudaAlta_UmbralSuperadoRetornaAlerta
TestAlertaDeudaAlta_UmbralNoSuperadoNoAlerta
```

---

## Definición de "Hecho" (DoD)
- [ ] Todos los tests pasan.
- [ ] Cobertura ≥ 90%.
- [ ] No es posible registrar un pago mayor a la deuda.
- [ ] El estado de cuenta refleja pagos en tiempo real.
- [ ] La alerta de deuda alta aparece en el dashboard cuando corresponde.
