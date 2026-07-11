# FC-09 — Reservas (Landing Page → Agenda del Restaurante)
**Estado:** Implementado
**Sprint:** 6 (ampliación posterior al Sprint 4)
**Prioridad:** Media
**Depende de:** FC-08 (Landing Page pública)

---

## Descripción
Permite que un cliente reserve una mesa desde la landing page pública, **sin crear un pedido real**. Una reserva solo registra la intención de visita (nombre, contacto, fecha, número de personas); no descuenta inventario ni genera un pedido. Cuando el cliente llega y empieza a consumir, el mesero registra el `Pedido` real desde el dashboard (FC-03), que es donde sí se descuenta stock y se calcula el total.

Se modela como una entidad independiente de `Pedido` — sin relación por FK — porque una reserva puede cancelarse o no concretarse nunca en una visita real.

---

## Endpoints

| Método | Ruta                              | Rol       | Descripción                                  |
|--------|------------------------------------|-----------|-----------------------------------------------|
| POST   | /api/public/reservas               | sin auth  | Crear una reserva desde la landing page       |
| GET    | /api/empleado/reservas             | any       | Listar reservas (paginado, filtro por estado) |
| PUT    | /api/empleado/reservas/:id/estado  | any       | Confirmar o cancelar una reserva              |
| PUT    | /api/empleado/reservas/:id/pedido  | any       | Vincular la reserva a un `Pedido` real ya creado |

---

## Contrato de la API

### POST /api/public/reservas — body:
```json
{
  "nombre": "Juan Pérez",
  "whatsapp": "987654321",
  "fecha": "2026-08-01",
  "personas": "4",
  "ocasion": "cumpleaños"
}
```
Respuesta `201`:
```json
{
  "id": 12,
  "nombre": "Juan Pérez",
  "whatsapp": "987654321",
  "fecha": "2026-08-01T00:00:00Z",
  "personas": "4",
  "ocasion": "cumpleaños",
  "estado": "pendiente",
  "created_at": "2026-07-11T10:00:00Z"
}
```

### GET /api/empleado/reservas?estado=pendiente
```json
{
  "data": [ { "id": 12, "nombre": "Juan Pérez", "estado": "pendiente", "...": "..." } ],
  "total": 1,
  "page": 1,
  "per_page": 20
}
```

### PUT /api/empleado/reservas/12/estado — body:
```json
{ "estado": "confirmada" }
```

### PUT /api/empleado/reservas/12/pedido — body:
```json
{ "pedido_id": 45 }
```
Se usa cuando el cliente llega y el mesero ya registró el `Pedido` real (FC-03) — este endpoint solo enlaza los dos registros para trazabilidad; valida que el `pedido_id` exista (`404` si no) antes de guardar la relación.

---

## Validaciones
- `nombre`, `whatsapp`, `fecha`, `personas`: requeridos.
- `fecha`: debe venir en formato `YYYY-MM-DD`; si no, `400`.
- `estado` (al actualizar): debe ser `pendiente`, `confirmada` o `cancelada`; cualquier otro valor → `400`.
- Toda reserva nace en estado `pendiente` — el cambio a `confirmada`/`cancelada` lo hace el staff desde el dashboard, no el cliente.
- `pedido_id` (al vincular): requerido y debe corresponder a un `Pedido` existente; si no existe → `404`.

---

## Decisiones de diseño (por qué no es parte de `Pedido`)
- Una reserva pública no debe poder descontar stock ni generar un total por error.
- El flujo real es: **Reserva (agenda)** → el cliente llega → **Pedido (venta real, con productos y stock)** → vinculado a la reserva vía `reservas.pedido_id` (FK opcional, `NULL` hasta que se concreta) solo para saber que la visita ocurrió.
- `reservas.pedido_id` es una FK real hacia `pedidos.id` (antes esta tabla no tenía ninguna relación; ver PLAN_SDD.md sección 5.2), pero se mantiene **opcional y de solo lectura desde el lado de Reserva** — nunca se crea un `Pedido` automáticamente a partir de una reserva, ni se descuenta stock por su causa.
- Como filtro anti-broma de bajo costo (sin integrar pasarela de pago), la confirmación de la reserva es un paso manual del staff (llamar/escribir al WhatsApp dejado en el formulario) antes de pasarla a `confirmada`.

---

## Pruebas

**Unitarias — handler con mock de repositorio** (`reservas_handler_test.go`, 17 tests):
```
TestCrearReserva_DatosValidosRetorna201
TestCrearReserva_SinNombreRetorna400
TestCrearReserva_FechaInvalidaRetorna400
TestCrearReserva_ErrorBDRetorna500
TestListarReservas_RetornaPaginacion
TestListarReservas_FiltrarPorEstado
TestListarReservas_ErrorBDRetorna500
TestActualizarEstadoReserva_EstadoValidoRetorna200
TestActualizarEstadoReserva_EstadoInvalidoRetorna400
TestActualizarEstadoReserva_IDInexistenteRetorna404
TestActualizarEstadoReserva_IDInvalidoRetorna400
TestVincularPedidoReserva_PedidoExisteRetorna200
TestVincularPedidoReserva_PedidoInexistenteRetorna404
TestVincularPedidoReserva_ReservaInexistenteRetorna404
TestVincularPedidoReserva_SinPedidoIDRetorna400
TestVincularPedidoReserva_IDInvalidoRetorna400
TestVincularPedidoReserva_ErrorBuscarPedidoRetorna500
```

**Integración — repositorio contra Postgres real** (`reservas_repo_test.go`, 8 tests):
```
TestReservaRepo_Create_GuardaReservaEnBD
TestReservaRepo_FindByID_ExisteRetornaDatos
TestReservaRepo_FindByID_InexistenteRetornaNil
TestReservaRepo_FindAll_RetornaPaginadoYFiltraPorEstado
TestReservaRepo_UpdateEstado_ActualizaCorrectamente
TestReservaRepo_UpdateEstado_InexistenteRetornaError
TestReservaRepo_VincularPedido_ActualizaCorrectamente
TestReservaRepo_VincularPedido_InexistenteRetornaError
```

---

## Definición de "Hecho" (DoD)
- [x] Todos los tests unitarios y de integración pasan.
- [x] `POST /api/public/reservas` no requiere JWT.
- [x] `GET` y `PUT` de reservas requieren JWT (cualquier rol).
- [x] Una reserva creada nunca descuenta stock ni crea un `Pedido`.
- [x] El frontend público (`ReservaPage.tsx`) llama al endpoint real (ya no simula el envío).
- [x] El dashboard (`ReservasPage.tsx`) permite confirmar/cancelar reservas pendientes.
- [ ] Cobertura ≥ 90% verificada en el próximo run de CI (pendiente de re-ejecutar el pipeline completo).
