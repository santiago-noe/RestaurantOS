# FC-06 — Reportes
**Estado:** Especificado  
**Sprint:** 3  
**Prioridad:** Media  
**Depende de:** FC-03, FC-04, FC-05

---

## Descripción
Generación de reportes para tomar decisiones de negocio. El admin puede ver ventas, deudas y movimientos por rango de fecha, y exportarlos.

---

## Endpoints

| Método | Ruta                          | Rol   | Descripción                        |
|--------|-------------------------------|-------|------------------------------------|
| GET    | /api/admin/reportes/ventas    | admin | Reporte de ventas por rango        |
| GET    | /api/admin/reportes/deudas    | admin | Reporte de clientes con deuda      |
| GET    | /api/admin/reportes/inventario| admin | Reporte de movimientos de stock    |
| GET    | /api/admin/reportes/ventas/pdf| admin | Exportar reporte ventas como PDF   |
| GET    | /api/admin/reportes/ventas/excel| admin| Exportar reporte ventas como Excel|

---

## Query params para /reportes/ventas

```
GET /api/admin/reportes/ventas?desde=2026-05-01&hasta=2026-05-31
```

**Response:**
```json
{
  "periodo": {
    "desde": "2026-05-01",
    "hasta": "2026-05-31"
  },
  "resumen": {
    "total_ventas": 4520.00,
    "ventas_contado": 2100.00,
    "ventas_credito": 2420.00,
    "total_pedidos": 185,
    "ticket_promedio": 24.43
  },
  "por_tipo_comida": [
    { "tipo": "almuerzo", "cantidad": 120, "total": 3000.00 },
    { "tipo": "desayuno", "cantidad": 50, "total": 1000.00 },
    { "tipo": "bebida", "cantidad": 15, "total": 520.00 }
  ],
  "por_dia": [
    { "fecha": "2026-05-01", "total": 145.00 },
    { "fecha": "2026-05-02", "total": 210.00 }
  ]
}
```

---

## Pruebas Unitarias

```
TestReporteVentas_RangoValidoRetornaDatos
TestReporteVentas_FechaDesdeHastaInvertidaRetorna400
TestReporteVentas_SinPedidosEnRangoRetornaCeros
TestCalcularTicketPromedio_PedidosConTotalRetornaPromedio
TestCalcularTicketPromedio_SinPedidosRetornaCero
TestReporteDeudas_ClientesConDeudaAparecen
TestReporteDeudas_ClientesSinDeudaNoAparecen
TestExportarPDF_GeneraArchivoNoVacio
TestExportarExcel_GeneraArchivoNoVacio
```

---

## Definición de "Hecho" (DoD)
- [ ] Todos los tests pasan.
- [ ] Cobertura ≥ 90%.
- [ ] Fechas inválidas retornan error descriptivo.
- [ ] PDF y Excel se generan correctamente y son descargables.
- [ ] El resumen numérico es matemáticamente correcto (verificar con datos de seed).
