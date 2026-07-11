# RestaurantOS — Sistema de Gestión para Restaurante

Proyecto Final — Pruebas y Aseguramiento de Calidad de Software  
Metodología: **SDD (Spec Driven Development)**

---

## Flujo SDD de este proyecto

```
1. PLAN_SDD.md           ← Especificación general (ya hecho)
2. feature-cards/FC-*.md ← Una Feature Card por módulo (ya hecho)
3. Tests unitarios       ← Escribir ANTES del código (siguiente paso)
4. Implementación        ← Código que hace pasar los tests
5. Cobertura ≥ 90%       ← Verificar con go test -coverprofile
6. CI/CD GitHub Actions  ← Tests automáticos en cada push
7. Deploy                ← Railway
```

---

## Estructura del proyecto

```
PROYECTO FINAL/
├── PLAN_SDD.md              ← Plan completo: requisitos, BD, arquitectura
├── README.md                ← Este archivo
├── feature-cards/
│   ├── FC-01-autenticacion.md
│   ├── FC-02-clientes.md
│   ├── FC-03-pedidos.md
│   ├── FC-04-creditos-pagos.md
│   ├── FC-05-inventario.md
│   ├── FC-06-reportes.md
│   ├── FC-07-ia-basica.md
│   ├── FC-08-landing-page.md
│   └── FC-09-reservas.md
├── database/
│   └── seed_dev.sql
├── backend/                 ← Go + Gin + GORM (implementado)
└── frontend/                ← React + Vite + Tailwind (implementado)
```

---

## Estado actual (2026-07-11)

Todos los módulos de FC-01 a FC-09 están implementados (backend + frontend), con 208 pruebas entre unitarias (mocks) y de integración (Postgres real). Ver el detalle en [PLAN_SDD.md, sección 8](PLAN_SDD.md#8-planificación-de-entregables) y la sección 7.1.1 (inventario real de pruebas).

Para añadir un módulo nuevo, sigue el mismo flujo SDD:

1. Escribir/actualizar la Feature Card en `feature-cards/FC-XX-nombre.md`.
2. Escribir los tests ANTES del código (`*_test.go`, con mocks para lo unitario y `setupTestDB` para integración).
3. Ejecutar `go test ./...` → debe fallar mientras no exista el código.
4. Implementar hasta que todos los tests pasen.
5. Verificar cobertura ≥ 90% (`go tool cover -func=coverage.out`).

### Verificar cobertura

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
```

---

## Tecnologías

| Capa       | Tech                    |
|------------|-------------------------|
| Backend    | Go 1.22 + Gin + GORM    |
| Base datos | PostgreSQL 15           |
| Frontend   | React + TypeScript + Vite + Tailwind |
| Auth       | JWT (HS256)             |
| Tests Go   | testing + testify       |
| Tests React| Vitest + RTL            |
| CI/CD      | GitHub Actions          |
| Deploy     | Railway                 |

---

## Credenciales de desarrollo (seed)

| Email                       | Password     | Rol      |
|-----------------------------|--------------|----------|
| admin@restaurante.com       | password123  | admin    |
| carlos@restaurante.com      | password123  | empleado |
