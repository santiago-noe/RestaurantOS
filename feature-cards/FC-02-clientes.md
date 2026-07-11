# FC-02 — Gestión de Clientes
**Estado:** Especificado  
**Sprint:** 2  
**Prioridad:** Alta  
**Depende de:** FC-01 (JWT)

---

## Descripción
CRUD completo de clientes. Un cliente puede ser persona individual o empresa. Se lleva registro de su deuda total (calculada automáticamente por trigger en la BD).

---

## Endpoints

| Método | Ruta                   | Rol      | Descripción                        |
|--------|------------------------|----------|------------------------------------|
| GET    | /api/empleado/clientes | any      | Listar todos (con paginación)      |
| GET    | /api/empleado/clientes/:id | any  | Ver cliente + deuda + historial    |
| POST   | /api/admin/clientes    | admin    | Crear cliente                      |
| PUT    | /api/admin/clientes/:id | admin   | Editar cliente                     |
| DELETE | /api/admin/clientes/:id | admin   | Desactivar (soft delete)           |

---

## Contrato de la API

### GET /api/empleado/clientes
```json
{
  "data": [
    {
      "id": 1,
      "nombre": "Juan",
      "apellido": "Pérez",
      "tipo": "individual",
      "telefono": "987654321",
      "deuda_total": 120.50
    }
  ],
  "total": 45,
  "page": 1,
  "per_page": 20
}
```

### POST /api/admin/clientes — body:
```json
{
  "nombre": "Constructora Norte",
  "apellido": "",
  "tipo": "empresa",
  "telefono": "01-234567",
  "direccion": "Av. Principal 123",
  "email": "pagos@constructoranorte.com"
}
```

---

## Validaciones
- `nombre`: requerido, max 100 caracteres.
- `tipo`: requerido, debe ser `individual` o `empresa`.
- `telefono`: opcional, formato libre.
- `email`: opcional, formato email válido.

---

## Pruebas Unitarias

```
TestCrearCliente_DatosValidosRetorna201
TestCrearCliente_SinNombreRetorna400
TestCrearCliente_TipoInvalidoRetorna400
TestCrearCliente_EmailDuplicadoRetorna409
TestObtenerCliente_IDExisteRetornaDatos
TestObtenerCliente_IDInexistenteRetorna404
TestListarClientes_RetornaPaginacion
TestListarClientes_FiltrarPorTipo
TestEditarCliente_DatosValidosActualiza
TestDesactivarCliente_SoftDeleteNoElimina
TestDesactivarCliente_ClienteConDeudaAlerta  ← edge case importante
```

---

## Definición de "Hecho" (DoD)
- [ ] Todos los tests pasan.
- [ ] Cobertura ≥ 90%.
- [ ] Un empleado NO puede crear/editar/eliminar clientes (403).
- [ ] La paginación funciona correctamente.
- [ ] Soft delete: cliente desactivado no aparece en listado normal.
