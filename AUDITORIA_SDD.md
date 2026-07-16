# Auditoría de Ciclo de Vida SDD — RestaurantOS

**Fecha de la auditoría:** 2026-07-14
**Repositorio auditado:** `santiago-noe/RestaurantOS` (rama `main`)
**Metodología:** Inspección directa de archivos, `git log`, `stat` (timestamps NTFS reales del sistema de archivos) y ejecución real de `go test ./... -coverprofile`. Ningún dato de este documento fue inventado; donde la evidencia es insuficiente o ambigua, se indica explícitamente.

---

## ⚠️ ADVERTENCIA METODOLÓGICA CRÍTICA — leer antes de citar cualquier fecha de `git log`

Este repositorio de Git **se inicializó el 2026-07-11**, no al inicio del proyecto. Antes de esa fecha, el proyecto existía como carpeta local sin control de versiones. Evidencia:

```
$ git log --oneline --all
4286193 Corrige el sidebar del dashboard forzado siempre visible en movil
e5a37a3 Reemplaza iconos de Material Symbols (fuente externa) por lucide-react
9e2c646 Saca los capitulos del informe/tesis del repo de codigo
ae0da3d Elimina workflow deploy.yml redundante
0eb43f4 Agrega middleware CORS para permitir consumo desde otro dominio (Vercel)
df68924 Elimina imports no usados que rompian el build de produccion (tsc -b)
58a3451 Hace configurable la URL base de la API del frontend via VITE_API_URL
4b54d50 Commit inicial: RestaurantOS (SDD) - backend Go+Gin+GORM, frontend React+Vite, feature cards y tests
```

**Solo hay 8 commits, todos del 2026-07-11 en adelante.** El commit `4b54d50` ("Commit inicial") agregó de una sola vez prácticamente todo el proyecto (backend, frontend, feature cards, `PLAN_SDD.md`). Por lo tanto:

- **`git log --follow` para cualquier archivo anterior al 2026-07-11 devolverá la fecha `2026-07-11` y el hash `4b54d50`**, sin importar cuándo se escribió realmente ese archivo. Usar esa fecha como "fecha de creación" en la tesis sería **incorrecto y no verificable**.
- La única fuente de verdad disponible para la cronología **real** anterior al 11 de julio es el **timestamp de creación (`Birth`) del sistema de archivos NTFS**, obtenido con `stat <archivo>`. Estos timestamps se preservaron porque los archivos no fueron reescritos después de su creación original — excepto los que se indican explícitamente como "modificados en esta sesión" (en cuyo caso su `Birth` se reseteó a la fecha de la edición, y por tanto **no sirve como evidencia de su fecha original**).
- A partir del 2026-07-11, sí hay commits reales incrementales (CORS, fix de íconos, fix de sidebar móvil, etc.) que **sí reflejan cronología real de cambios**, porque ocurrieron dentro de la ventana de vida de este repositorio Git.

En las tablas siguientes, la columna **"Commit"** se completa solo cuando aporta información real (p. ej. cuándo se *eliminó* un archivo, o cambios posteriores al 11 de julio). Para archivos anteriores a esa fecha, se usa la columna **"Fecha real (FS)"** basada en `stat`, y se marca `git: 4b54d50 (irrelevante para cronología)` para dejar constancia de que el commit no aporta información temporal útil.

---

## 1. FASE DE ANÁLISIS / ESPECIFICACIÓN

### 1.1 Documentos de especificación general

| Fecha real (FS) | Archivo | Descripción | Commit (git) |
|---|---|---|---|
| 2026-05-29 (fecha declarada en el propio documento: `**Fecha:** 2026-05-29`; el `Birth` en disco marca 2026-07-11 porque el archivo fue editado varias veces en esta sesión para incorporar FC-09 y las nuevas relaciones de FK) | `PLAN_SDD.md` | Documento SDD maestro: metodología, visión general, SRS completo (RF-01 a RF-08, RNF-01 a RNF-10), historias de usuario, diseño de BD normalizada, arquitectura, plan de pruebas, planificación de sprints | `4b54d50` (versión original) + `4286193`…(ediciones posteriores de esta sesión, ver sección 6 de este mismo documento para el detalle) |
| 2026-05-29 12:55 aprox. (Birth ahora 2026-07-11 por edición en esta sesión) | `README.md` | Índice del proyecto, flujo SDD resumido, credenciales de seed | `4b54d50` |
| No existe un archivo SRS separado | — | Los requisitos funcionales y no funcionales viven **dentro** de `PLAN_SDD.md` sección 3 (no hay un `SRS.md` independiente) | — |

### 1.2 Feature Cards (FC-01 a FC-09)

| Fecha real (FS, `Birth`) | Archivo | RF/RNF que cubre | Commit |
|---|---|---|---|
| **2026-05-29 15:59:52** | `feature-cards/FC-01-autenticacion.md` | RF-01, RNF-03, RNF-04 | `4b54d50` (sin cambios desde el original) |
| **2026-05-29 16:00:12** | `feature-cards/FC-02-clientes.md` | RF-02 | `4b54d50` (sin cambios desde el original) |
| **2026-05-29 16:00:35** | `feature-cards/FC-03-pedidos.md` | RF-03 | `4b54d50` (sin cambios desde el original) |
| **2026-05-29 16:01:22** | `feature-cards/FC-04-creditos-pagos.md` | RF-04 | `4b54d50` (sin cambios desde el original) |
| **2026-05-29 16:02:39** | `feature-cards/FC-05-inventario.md` | RF-05 | `4b54d50` (sin cambios desde el original) |
| **2026-05-29 16:03:21** | `feature-cards/FC-06-reportes.md` | RF-06 | `4b54d50` (sin cambios desde el original) |
| **2026-05-29 16:05:04** | `feature-cards/FC-07-ia-basica.md` | RF-07 | `4b54d50` (sin cambios desde el original) |
| Original ~2026-05-29 / editado 2026-07-11 13:10:21 (esta sesión: se agregó nota de `producto_id` y referencia cruzada a FC-09) | `feature-cards/FC-08-landing-page.md` | RF-08 | `4b54d50` |
| **2026-07-11 13:09:53** (creado íntegramente en esta sesión colaborativa, NO forma parte de la especificación original del curso) | `feature-cards/FC-09-reservas.md` | Requisito nuevo, no mapea a ningún RF-01..RF-08 original — ver nota abajo | `4b54d50` |

> **Nota importante para la tesis:** FC-09 (Reservas) **no estaba en el alcance original documentado en `PLAN_SDD.md`**. Surgió de una conversación de diseño posterior ("el cliente no puede hacer su pedido desde el landing, solo reserva") y se especificó, implementó y probó dentro de esta misma sesión de trabajo, el 2026-07-11. Si la tesis presenta un ciclo SDD "ideal" de 8 features, FC-09 debe presentarse explícitamente como una **ampliación posterior (Sprint 6)**, no como parte del SRS inicial — y así quedó documentado en `PLAN_SDD.md` sección 8.

### 1.3 RF/RNF cubiertos y criterios de aceptación por Feature Card

| FC | RF cubierto (según `PLAN_SDD.md` §3.1) | Criterios de aceptación (DoD tal como aparecen en el archivo) |
|---|---|---|
| FC-01 | RF-01 — Autenticación y Roles | Todos los tests pasan · Cobertura módulo auth ≥90% · Login funciona con admin y empleado de seed · Ruta protegida rechaza requests sin token · El token expira correctamente |
| FC-02 | RF-02 — Gestión de Clientes | Todos los tests pasan · Cobertura ≥90% · Empleado NO puede crear/editar/eliminar clientes (403) · Paginación funciona · Soft delete no aparece en listado normal |
| FC-03 | RF-03 — Gestión de Pedidos | Todos los tests pasan · Cobertura ≥90% · Pedido a crédito actualiza `deuda_total` · Anular pedido devuelve stock y reduce deuda · Todo ocurre en una transacción |
| FC-04 | RF-04 — Créditos y Pagos | Todos los tests pasan · Cobertura ≥90% · No se puede pagar más de la deuda · Estado de cuenta refleja pagos en tiempo real · Alerta de deuda alta visible en dashboard |
| FC-05 | RF-05 — Inventario | Todos los tests pasan · Cobertura ≥90% · Stock nunca queda negativo · Cada movimiento queda registrado con fecha y referencia · Alertas de stock bajo visibles |
| FC-06 | RF-06 — Reportes | Todos los tests pasan · Cobertura ≥90% · Fechas inválidas retornan error descriptivo · PDF/Excel se generan y son descargables · Resumen numérico matemáticamente correcto |
| FC-07 | RF-07 — Dashboard IA | Todos los tests pasan · Cobertura ≥90% · Sin historial → predicción retorna 0, no error 500 · Predicción solo usa el día de la semana correcto · Alertas visibles en tiempo real |
| FC-08 | RF-08 — Página Pública | Todos los tests pasan (nota: tests especificados en **Vitest/RTL, frontend** — no existen actualmente, ver sección 5) · Página responsive · Menú se actualiza sin redeploy · Formulario de contacto valida campos vacíos · Lighthouse ≥80 |
| FC-09 | (no mapea a RF original — ampliación) | Todos los tests unitarios/integración pasan · POST público sin JWT · GET/PUT requieren JWT · Reserva nunca descuenta stock ni crea Pedido · Frontend público conectado al backend real · Dashboard permite confirmar/cancelar · Cobertura ≥90% marcada como **pendiente de verificar** en el propio archivo |

---

## 2. FASE DE DISEÑO (ARQUITECTURA)

**No existe documentación formal de arquitectura separada (no hay ADRs, no hay carpeta `docs/adr/`, no hay diagrama en herramienta externa).** Búsqueda realizada:

```
$ find . -maxdepth 4 -iname "*adr*" -o -iname "*architecture*" -o -iname "*arquitectura*"
(sin resultados)
```

La arquitectura de 3 capas está documentada **únicamente** como texto/ASCII dentro de `PLAN_SDD.md`, sección 6 ("ARQUITECTURA DEL SISTEMA"), y se **infiere adicionalmente** de la estructura real de carpetas, verificada con `find`:

### Árbol real `backend/` (3 niveles)
```
backend
backend/cmd
backend/cmd/genhash
backend/cmd/migrate
backend/cmd/seed
backend/cmd/server
backend/internal
backend/internal/auth
backend/internal/config
backend/internal/database
backend/internal/export
backend/internal/handlers
backend/internal/middleware
backend/internal/models
backend/internal/repository
backend/internal/services
```

### Árbol real `frontend/src/` (3 niveles)
```
frontend/src
frontend/src/assets
frontend/src/components
frontend/src/components/layout
frontend/src/components/ui
frontend/src/context
frontend/src/hooks
frontend/src/pages
frontend/src/pages/dashboard
frontend/src/pages/public
frontend/src/services
frontend/src/types
```

Esta estructura confirma el patrón declarado: **handlers → services/repository → models**, con `middleware` como capa transversal (JWT, CORS), consistente con una arquitectura en capas (Presentación / Negocio / Datos), aunque **nunca se formalizó en un documento de diseño dedicado** más allá del ASCII de `PLAN_SDD.md`.

---

## 3. FASE DE DISEÑO DE BASE DE DATOS

### 3.1 Migraciones

**No existen migraciones SQL versionadas ni carpeta `migrations/`.** Verificado:

```
$ ls -la database/
seed_dev.sql   (único archivo, 2026-05-31 11:01, no es una migración — es un script de datos semilla)

$ find . -iname "migrations" -type d
(sin resultados)
```

El esquema de base de datos se gestiona **exclusivamente vía GORM `AutoMigrate`**, definido en código Go (`backend/internal/database/database.go`), no en archivos `.sql` versionados. Esto contradice lo que describe `PLAN_SDD.md` sección 5.3 (que muestra `CREATE TABLE` como si fueran migraciones formales) — esos `CREATE TABLE` son **documentación descriptiva del esquema, no migraciones ejecutables reales**. Si la tesis afirma que existen "migraciones numeradas", **eso no es correcto**; debe describirse como "esquema gestionado por ORM (AutoMigrate)".

### 3.2 Modelo de datos actual (fuente de verdad real: `backend/internal/models/models.go`)

| Tabla (struct Go) | PK | FKs | Notas de normalización |
|---|---|---|---|
| `User` | `ID` | — | Email con `uniqueIndex` |
| `Cliente` | `ID` | — | `DeudaTotal` es campo calculado, no tiene trigger de BD real (ver brecha en sección 5) |
| `Producto` | `ID` | — | — |
| `Pedido` | `ID` | `ClienteID→Cliente`, `UserID→User` | — |
| `PedidoItem` | `ID` | `PedidoID→Pedido`, `ProductoID→Producto` | Resuelve 2FN (detalle separado de cabecera) |
| `Pago` | `ID` | `ClienteID→Cliente`, `PedidoID→Pedido` (nullable) | — |
| `MovimientoStock` | `ID` | `ProductoID→Producto` | `ReferenciaID` es un `*int` sin FK real (referencia lógica, no constraint) |
| `Reserva` | `ID` | `PedidoID→Pedido` (nullable, agregado en esta sesión 2026-07-11) | Sin FK a `Cliente` — diseño intencional, ver `FC-09` |
| `MenuPublico` | `ID` | `ProductoID→Producto` (nullable, agregado en esta sesión 2026-07-11) | — |

**Verificado en BD real** (Postgres, vía `\d`): las FK `fk_reservas_pedido` y `fk_menu_publicos_producto` existen físicamente como constraints, confirmado el 2026-07-11 tras correr `AutoMigrate` contra la base de datos de Railway.

> ⚠️ **Discrepancia detectada:** `PLAN_SDD.md` sección 5.4 documenta un **trigger SQL** (`trg_deuda_pedido`, `trg_deuda_pago`) para mantener `clientes.deuda_total` sincronizado automáticamente. **Ese trigger no existe en la base de datos real** — no hay ningún archivo `.sql` que lo cree, y `AutoMigrate` de GORM no ejecuta triggers. La actualización de `deuda_total` en el código real se hace **manualmente en Go** (`services/pedido_service.go` acumula el total en memoria; el seed también actualiza `deuda_total` con un `UPDATE` explícito). Es decir: el documento de diseño describe un mecanismo (trigger de BD) que el código implementado **no usa**.

---

## 4. FASE DE IMPLEMENTACIÓN (TDD)

### 4.1 Inventario código ↔ test por módulo, con conteo real (`grep -c "^func Test"`)

| Módulo | Archivo implementación | Birth (FS) | Archivo test | Birth (FS) | N° tests |
|---|---|---|---|---|---|
| auth | `internal/auth/auth.go` | 2026-05-29 16:27:21 | `internal/auth/auth_test.go` | 2026-05-29 16:26:48 | 14 |
| middleware | `internal/middleware/jwt.go` | 2026-05-29 16:27:35 | `internal/middleware/jwt_test.go` | 2026-05-29 17:13:08 | 8 |
| handlers/auth | `internal/handlers/auth_handler.go` | 2026-05-29 17:20:14 | `internal/handlers/auth_handler_test.go` | 2026-05-29 17:20:38 | 8 |
| handlers/clientes | `internal/handlers/clientes_handler.go` | 2026-05-29 17:16:20 | `internal/handlers/clientes_handler_test.go` | 2026-05-29 17:22:08 | 23 |
| handlers/creditos | `internal/handlers/creditos_handler.go` | 2026-05-29 17:50:41 | `internal/handlers/creditos_handler_test.go` | 2026-05-29 18:00:14 | 12 |
| handlers/inventario | `internal/handlers/inventario_handler.go` | 2026-05-29 | `internal/handlers/inventario_handler_test.go` | 2026-07-07 (reescrito luego) | 24 |
| handlers/menu | `internal/handlers/menu_handler.go` | 2026-07-11 (editado esta sesión: campo `producto_id`) | `internal/handlers/menu_handler_test.go` | 2026-07-07 | 13 |
| handlers/pedidos | `internal/handlers/pedidos_handler.go` | 2026-07-07 02:12:11 | `internal/handlers/pedidos_handler_test.go` | 2026-07-07 02:14:59 | 16 |
| handlers/reportes | `internal/handlers/reportes_handler.go` | 2026-07-07 01:33:58 | `internal/handlers/reportes_handler_test.go` | 2026-07-07 01:34:35 | 11 |
| handlers/reservas | `internal/handlers/reservas_handler.go` | 2026-07-11 13:00:38 | `internal/handlers/reservas_handler_test.go` | 2026-07-11 13:02:57 | 17 |
| repository/clientes | `internal/repository/clientes_repo.go` | 2026-05-29 | `internal/repository/clientes_repo_test.go` | 2026-07-07 (reescrito luego) | 8 |
| repository/menu | `internal/repository/menu_repo.go` | 2026-07-07 01:01:00 | `internal/repository/menu_repo_test.go` | 2026-07-07 01:01:19 | 8 |
| repository/reservas | `internal/repository/reservas_repo.go` | 2026-07-11 13:00:02 | `internal/repository/reservas_repo_test.go` | 2026-07-11 13:03:22 | 8 |
| services/ia | `internal/services/ia_service.go` | 2026-05-29 17:51:30 | `internal/services/ia_service_test.go` | 2026-05-29 17:51:11 | 8 |
| services/pedido | `internal/services/pedido_service.go` | 2026-07-07 02:10:38 | `internal/services/pedido_service_test.go` | 2026-07-07 02:10:59 | 24 |
| services/reportes | `internal/services/reportes_service.go` | 2026-07-07 01:30:40 | `internal/services/reportes_service_test.go` | 2026-07-07 01:31:04 | 11 |
| export | `internal/export/export.go` | 2026-07-07 01:32:08 | `internal/export/export_test.go` | 2026-07-07 01:32:25 | 7 |
| **repository/pedidos** | `internal/repository/pedidos_repo.go` | 2026-07-07 | **NO EXISTE archivo de test** | — | **0** |
| **repository/user** | `internal/repository/user_repo.go` | 2026-05-29 | **NO EXISTE archivo de test** | — | **0** |
| **middleware/cors** | `internal/middleware/cors.go` | 2026-07-11 (nuevo, esta sesión) | **NO EXISTE archivo de test** | — | **0** |

**Total de funciones de test reales en el repositorio: 208** (suma de la columna "N° tests" arriba). Este número reemplaza cualquier cifra distinta mencionada en versiones previas de `PLAN_SDD.md` dentro de esta misma conversación — se recalculó de cero para esta auditoría con `grep -c` directo sobre cada archivo, sin arrastrar cifras anteriores.

### 4.2 Orden cronológico real de implementación por commit (reconstrucción de "sprints")

Como se explicó en la advertencia metodológica, `git log` **no** puede reconstruir el orden real antes del 2026-07-11 (todo entra en un solo commit). La única evidencia de orden cronológico real y granular **antes** de esa fecha es el timestamp `Birth` de cada archivo (tabla 4.1), que permite reconstruir agrupaciones por día:

- **2026-05-29** (día 1 de implementación): `auth`, `middleware/jwt`, `handlers/auth`, `handlers/clientes`, `handlers/creditos`, `handlers/inventario` (parcial), `repository/clientes` (parcial), `repository/user`, `services/ia` → corresponde a **Sprint 1 y parte de Sprint 2** según `PLAN_SDD.md` §8.
- **2026-05-30 a 2026-06-02**: archivos de frontend base (`PrivateRoute.tsx`, `main.tsx`, `IAPage.tsx`, `AuthContext.tsx`, `InventarioPage.tsx`, `DashboardHome.tsx`) y `database/seed_dev.sql` (2026-05-31) → **Sprint 4 (frontend) iniciado en paralelo, antes de terminar todo el backend**, lo cual **contradice el orden secuencial "Sprint 1 → 2 → 3 → 4" declarado en `PLAN_SDD.md` §8** (el frontend arrancó antes de que terminara Sprint 3).
- **2026-07-07**: `handlers/pedidos`, `handlers/reportes`, `services/pedido`, `services/reportes`, `export`, `repository/menu`, y reescritura de varios tests existentes → **Sprint 2/3 completado**.
- **2026-07-11 en adelante**: Reservas (FC-09), CORS, fixes de despliegue, fix de íconos, fix de sidebar móvil → **Sprint 6 (ampliación), documentado con commits reales e incrementales**.

### 4.3 ¿Hay evidencia de que las pruebas se escribieron antes que el código (TDD estricto)?

**Respuesta honesta: la evidencia de timestamps NO respalda TDD estricto ("red-green") en la mayoría de los módulos.** Comparando el `Birth` exacto (con hora) del archivo de implementación contra su test, en los 14 pares donde ambos archivos comparten la misma fecha:

| Módulo | ¿Test nace antes que la implementación? |
|---|---|
| auth | **Sí** (33 segundos antes) |
| services/ia | **Sí** (19 segundos antes) |
| middleware | No (45 min después) |
| handlers/auth | No (24 s después) |
| handlers/clientes | No (~6 min después) |
| handlers/creditos | No (~10 min después) |
| handlers/pedidos | No (~3 min después) |
| handlers/reportes | No (~37 s después) |
| handlers/reservas | No (~2 min 19 s después) |
| repository/menu | No (~19 s después) |
| repository/reservas | No (~3 min 20 s después) |
| services/pedido | No (~21 s después) |
| services/reportes | No (~24 s después) |
| export | No (~17 s después) |

**Solo 2 de 14 pares comparables (auth, services/ia) muestran el archivo de test naciendo antes que el de implementación.** En los otros 12, el archivo de implementación se creó primero, y el test se creó segundos o minutos después — un patrón de **"implementar y luego probar de inmediato"**, no de TDD estricto en el sentido de "escribir el test, verlo fallar, luego implementar".

**Limitación honesta de este método:** el timestamp `Birth` de NTFS marca cuándo el *archivo* se creó en disco, no necesariamente cuándo se escribió el *contenido* significativo. Es posible que un editor haya creado ambos archivos vacíos casi simultáneamente y el desarrollador haya escrito el contenido en un orden distinto al de creación del archivo. Por tanto, esta tabla debe presentarse en la tesis como **"evidencia de creación de archivos", no como prueba definitiva del proceso mental de TDD**.

---

## 5. FASE DE VERIFICACIÓN / CI-CD

### 5.1 Contenido completo de `.github/workflows/test.yml` (vigente)

```yaml
name: Tests y Cobertura

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: restaurantos_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache-dependency-path: backend/go.sum

      - name: Instalar dependencias
        working-directory: ./backend
        run: go mod download

      - name: Correr tests con cobertura
        working-directory: ./backend
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/restaurantos_test?sslmode=disable
          JWT_SECRET: test-secret-ci
        run: go test ./... -coverprofile=coverage.out -covermode=atomic -v

      - name: Verificar cobertura minima 90% (paquetes de logica de negocio)
        working-directory: ./backend
        run: |
          grep -v "restaurantos/cmd/" coverage.out > coverage_negocio.out
          COVERAGE=$(go tool cover -func=coverage_negocio.out | grep "^total:" | awk '{print $3}' | tr -d '%')
          echo "Cobertura (internal/*): ${COVERAGE}%"
          if awk "BEGIN {exit !($COVERAGE < 90)}"; then
            echo "FALLO: Cobertura ${COVERAGE}% menor al 90% requerido"
            exit 1
          fi
          echo "OK: Cobertura ${COVERAGE}% cumple el minimo"

      - name: Subir reporte de cobertura
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: coverage-report
          path: backend/coverage.out
```

### 5.2 Contenido de `.github/workflows/deploy.yml` — **ELIMINADO, ya no existe en el repo**

Recuperado del historial de git (commit `4b54d50`, borrado en el commit `ae0da3d` "Elimina workflow deploy.yml redundante"):

```yaml
name: Deploy a Railway

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    needs: []

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Deploy via Railway CLI
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
        run: |
          npm install -g @railway/cli
          railway up --service backend
```

**Motivo real de la eliminación (no especulativo, documentado en el mensaje de commit):** este workflow nunca llegó a ejecutarse correctamente — GitHub Actions estuvo bloqueado a nivel de cuenta por un problema de facturación ("account is locked due to a billing issue", confirmado vía API de GitHub), y además el secreto `RAILWAY_TOKEN` nunca se configuró. El despliegue real a Railway ocurre por la integración nativa de Railway con GitHub (webhook directo), **no por GitHub Actions**. Esto es una **brecha real respecto a RNF-07** (ver 5.4).

### 5.3 Salida real de `go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out`

Ejecutado el 2026-07-14 contra una base de datos Postgres local real (`restaurantos_test_audit`, creada específicamente para esta auditoría, mismo motor que usa el proyecto). Salida completa, sin resumir:

```
	restaurantos/cmd/genhash		coverage: 0.0% of statements
	restaurantos/cmd/migrate		coverage: 0.0% of statements
	restaurantos/cmd/seed		coverage: 0.0% of statements
	restaurantos/cmd/server		coverage: 0.0% of statements
ok  	restaurantos/internal/auth	5.931s	coverage: 91.3% of statements
	restaurantos/internal/config		coverage: 0.0% of statements
	restaurantos/internal/database		coverage: 0.0% of statements
ok  	restaurantos/internal/export	4.499s	coverage: 100.0% of statements
ok  	restaurantos/internal/handlers	6.542s	coverage: 85.8% of statements
ok  	restaurantos/internal/middleware	1.325s	coverage: 76.5% of statements
?   	restaurantos/internal/models	[no test files]
ok  	restaurantos/internal/repository	3.002s	coverage: 49.1% of statements
ok  	restaurantos/internal/services	1.608s	coverage: 96.4% of statements
```

`go tool cover -func` completo (todas las funciones, sin omitir ninguna):

```
restaurantos/cmd/genhash/main.go:8:              main                    0.0%
restaurantos/cmd/migrate/main.go:10:              main                    0.0%
restaurantos/cmd/seed/main.go:39:                 pesoTipoComida          0.0%
restaurantos/cmd/seed/main.go:51:                 cantidadPara            0.0%
restaurantos/cmd/seed/main.go:59:                 main                    0.0%
restaurantos/cmd/seed/main.go:293:                round2                  0.0%
restaurantos/cmd/server/main.go:16:               main                    0.0%
restaurantos/internal/auth/auth.go:20:            HashPassword            100.0%
restaurantos/internal/auth/auth.go:28:            CheckPassword           100.0%
restaurantos/internal/auth/auth.go:36:            GenerateJWT             100.0%
restaurantos/internal/auth/auth.go:50:            ValidateJWT             83.3%
restaurantos/internal/config/config.go:15:        Load                    0.0%
restaurantos/internal/config/config.go:25:        getEnv                  0.0%
restaurantos/internal/database/database.go:13:    Connect                 0.0%
restaurantos/internal/database/database.go:23:    Migrate                 0.0%
restaurantos/internal/export/export.go:16:        VentasExcel             100.0%
restaurantos/internal/export/export.go:43:        VentasPDF               100.0%
restaurantos/internal/export/export.go:74:        DeudoresExcel           100.0%
restaurantos/internal/export/export.go:103:       DeudoresPDF             100.0%
restaurantos/internal/export/export.go:135:       MovimientosExcel        100.0%
restaurantos/internal/export/export.go:166:       MovimientosPDF          100.0%
restaurantos/internal/handlers/auth_handler.go:18:        NewAuthHandler          100.0%
restaurantos/internal/handlers/auth_handler.go:40:        Login                   87.5%
restaurantos/internal/handlers/auth_handler.go:76:        Me                      77.8%
restaurantos/internal/handlers/clientes_handler.go:18:    NewClienteHandler       100.0%
restaurantos/internal/handlers/clientes_handler.go:39:    Crear                   100.0%
restaurantos/internal/handlers/clientes_handler.go:74:    ObtenerPorID            100.0%
restaurantos/internal/handlers/clientes_handler.go:94:    Listar                  83.3%
restaurantos/internal/handlers/clientes_handler.go:120:   Actualizar              85.3%
restaurantos/internal/handlers/clientes_handler.go:174:   Desactivar              100.0%
restaurantos/internal/handlers/creditos_handler.go:18:    NewCreditosHandler      100.0%
restaurantos/internal/handlers/creditos_handler.go:22:    ListarDeudores          100.0%
restaurantos/internal/handlers/creditos_handler.go:39:    EstadoCuenta            93.3%
restaurantos/internal/handlers/creditos_handler.go:75:    RegistrarPago           90.9%
restaurantos/internal/handlers/inventario_handler.go:23:  NewInventarioHandler    100.0%
restaurantos/internal/handlers/inventario_handler.go:27:  Listar                  100.0%
restaurantos/internal/handlers/inventario_handler.go:36:  Alertas                 100.0%
restaurantos/internal/handlers/inventario_handler.go:69:  ObtenerPorID            93.3%
restaurantos/internal/handlers/inventario_handler.go:101: Crear                   100.0%
restaurantos/internal/handlers/inventario_handler.go:132: Actualizar              88.9%
restaurantos/internal/handlers/inventario_handler.go:179: Restock                 100.0%
restaurantos/internal/handlers/menu_handler.go:18:        NewMenuHandler          100.0%
restaurantos/internal/handlers/menu_handler.go:45:        Publico                 100.0%
restaurantos/internal/handlers/menu_handler.go:55:        Listar                  60.0%
restaurantos/internal/handlers/menu_handler.go:64:        Crear                   75.0%
restaurantos/internal/handlers/menu_handler.go:95:        Actualizar              62.2%
restaurantos/internal/handlers/menu_handler.go:153:       Eliminar                63.6%
restaurantos/internal/handlers/pedidos_handler.go:28:     NewPedidoHandler        100.0%
restaurantos/internal/handlers/pedidos_handler.go:47:     Crear                   88.2%
restaurantos/internal/handlers/pedidos_handler.go:88:     Listar                  90.9%
restaurantos/internal/handlers/pedidos_handler.go:112:    ObtenerPorID            83.3%
restaurantos/internal/handlers/pedidos_handler.go:132:    MarcarEntregado         100.0%
restaurantos/internal/handlers/pedidos_handler.go:147:    Anular                  100.0%
restaurantos/internal/handlers/reportes_handler.go:21:    NewReportesHandler      100.0%
restaurantos/internal/handlers/reportes_handler.go:25:    parseFecha              33.3%
restaurantos/internal/handlers/reportes_handler.go:37:    Ventas                  90.9%
restaurantos/internal/handlers/reportes_handler.go:76:    Deudores                53.3%
restaurantos/internal/handlers/reportes_handler.go:103:   Inventario              75.0%
restaurantos/internal/handlers/reservas_handler.go:19:    NewReservaHandler       100.0%
restaurantos/internal/handlers/reservas_handler.go:40:    Crear                   100.0%
restaurantos/internal/handlers/reservas_handler.go:71:    Listar                  83.3%
restaurantos/internal/handlers/reservas_handler.go:98:    ActualizarEstado        77.8%
restaurantos/internal/handlers/reservas_handler.go:131:   VincularPedido          90.9%
restaurantos/internal/middleware/cors.go:11:              CORSMiddleware          0.0%
restaurantos/internal/middleware/jwt.go:14:               JWTMiddleware           100.0%
restaurantos/internal/middleware/jwt.go:32:               RequireRole             100.0%
restaurantos/internal/middleware/jwt.go:50:               GetClaims               100.0%
restaurantos/internal/repository/clientes_repo.go:24:     NewClienteRepo          100.0%
restaurantos/internal/repository/clientes_repo.go:28:     Create                  100.0%
restaurantos/internal/repository/clientes_repo.go:32:     FindByID                100.0%
restaurantos/internal/repository/clientes_repo.go:41:     FindAll                 80.0%
restaurantos/internal/repository/clientes_repo.go:59:     Update                  66.7%
restaurantos/internal/repository/clientes_repo.go:66:     Deactivate              66.7%
restaurantos/internal/repository/clientes_repo.go:77:     EmailExists             80.0%
restaurantos/internal/repository/menu_repo.go:22:         NewMenuRepo             100.0%
restaurantos/internal/repository/menu_repo.go:26:         FindPublico             100.0%
restaurantos/internal/repository/menu_repo.go:32:         FindAll                 100.0%
restaurantos/internal/repository/menu_repo.go:38:         FindByID                100.0%
restaurantos/internal/repository/menu_repo.go:47:         Create                  100.0%
restaurantos/internal/repository/menu_repo.go:51:         Update                  66.7%
restaurantos/internal/repository/menu_repo.go:58:         Delete                  83.3%
restaurantos/internal/repository/pedidos_repo.go:24:      NewPedidoRepo           100.0%
restaurantos/internal/repository/pedidos_repo.go:26:      Create                  0.0%
restaurantos/internal/repository/pedidos_repo.go:38:      FindByID                0.0%
restaurantos/internal/repository/pedidos_repo.go:48:      FindAll                 0.0%
restaurantos/internal/repository/pedidos_repo.go:69:      FindEntreFechas         100.0%
restaurantos/internal/repository/pedidos_repo.go:76:      UpdateEstado            0.0%
restaurantos/internal/repository/pedidos_repo.go:100:     NewProductoRepo         0.0%
restaurantos/internal/repository/pedidos_repo.go:102:     FindByID                0.0%
restaurantos/internal/repository/pedidos_repo.go:111:     AjustarStock            0.0%
restaurantos/internal/repository/pedidos_repo.go:116:     FindAll                 0.0%
restaurantos/internal/repository/pedidos_repo.go:126:     Create                  0.0%
restaurantos/internal/repository/pedidos_repo.go:130:     Update                  0.0%
restaurantos/internal/repository/pedidos_repo.go:147:     NewMovimientoRepo       100.0%
restaurantos/internal/repository/pedidos_repo.go:149:     Registrar               0.0%
restaurantos/internal/repository/pedidos_repo.go:153:     FindByProducto          0.0%
restaurantos/internal/repository/pedidos_repo.go:159:     FindEntreFechas         100.0%
restaurantos/internal/repository/pedidos_repo.go:168:     ErrNotFound             0.0%
restaurantos/internal/repository/pedidos_repo.go:182:     NewPagoRepo             0.0%
restaurantos/internal/repository/pedidos_repo.go:184:     Create                  0.0%
restaurantos/internal/repository/pedidos_repo.go:189:     FindByCliente           0.0%
restaurantos/internal/repository/pedidos_repo.go:195:     SumaPagadoPorCliente    0.0%
restaurantos/internal/repository/reservas_repo.go:23:     NewReservaRepo          100.0%
restaurantos/internal/repository/reservas_repo.go:27:     Create                  100.0%
restaurantos/internal/repository/reservas_repo.go:31:     FindByID                100.0%
restaurantos/internal/repository/reservas_repo.go:40:     FindAll                 90.0%
restaurantos/internal/repository/reservas_repo.go:58:     UpdateEstado            83.3%
restaurantos/internal/repository/reservas_repo.go:69:     VincularPedido          83.3%
restaurantos/internal/repository/user_repo.go:20:         NewUserRepo             0.0%
restaurantos/internal/repository/user_repo.go:24:         FindByEmail             0.0%
restaurantos/internal/repository/user_repo.go:33:         FindByID                0.0%
restaurantos/internal/services/ia_service.go:29:          PredecirDemanda         100.0%
restaurantos/internal/services/ia_service.go:75:          GenerarAlertas          100.0%
restaurantos/internal/services/ia_service.go:108:         severidadOrden          75.0%
restaurantos/internal/services/ia_service.go:119:         nombreDia               100.0%
restaurantos/internal/services/pedido_service.go:29:      CalcularTotal           100.0%
restaurantos/internal/services/pedido_service.go:45:      NewPedidoService        100.0%
restaurantos/internal/services/pedido_service.go:49:      Crear                   92.6%
restaurantos/internal/services/pedido_service.go:117:     MarcarEntregado         90.0%
restaurantos/internal/services/pedido_service.go:135:     Anular                  92.3%
restaurantos/internal/services/pedido_service.go:162:     FindByID                100.0%
restaurantos/internal/services/pedido_service.go:166:     FindAll                 100.0%
restaurantos/internal/services/reportes_service.go:29:    RangoPeriodo            100.0%
restaurantos/internal/services/reportes_service.go:50:    GenerarReporteVentas    100.0%
restaurantos/internal/services/reportes_service.go:92:    FiltrarClientesConDeuda 100.0%
restaurantos/internal/services/reportes_service.go:112:   GenerarReporteMovimientos 100.0%
total:                                                    (statements)          67.8%
```

**Cobertura total (todos los paquetes, incluido `cmd/*`): 67.8%.**
**Cobertura excluyendo `cmd/*`** (misma regla exacta que usa `.github/workflows/test.yml` para su gate del 90%): **82.1%.**

### 5.4 🚨 Brechas reales detectadas — RF/RNF sin implementación correspondiente, o con cobertura insuficiente

Esta es la sección más importante para no repetir en la tesis solo lo que "debería" existir según la especificación:

1. **RNF-02 (Cobertura ≥90%) — NO SE CUMPLE actualmente.** La cobertura real medida hoy es **82.1%** (excluyendo `cmd/*`, mismo criterio del propio CI) o **67.8%** (global). El gate de `test.yml` **fallaría** si corriera ahora mismo. Principales responsables:
   - `internal/repository/pedidos_repo.go`: **0% de cobertura en casi todas sus funciones** (`Create`, `FindByID`, `FindAll`, `UpdateEstado` de `PedidoRepo`; `NewProductoRepo`, `FindByID`, `AjustarStock`, `FindAll`, `Create`, `Update` de `ProductoRepo`; `Registrar`, `FindByProducto` de `MovimientoRepo`; `NewPagoRepo`, `Create`, `FindByCliente`, `SumaPagadoPorCliente` de `PagoRepo`). **No existe `pedidos_repo_test.go`.** Esto afecta directamente a RF-03, RF-04 y RF-05, cuya capa de acceso a datos está completamente sin probar a nivel de integración (solo se prueba indirectamente vía mocks en los tests de `handlers`/`services`, que no ejecutan SQL real).
   - `internal/repository/user_repo.go`: **0% de cobertura**, no existe `user_repo_test.go`. Afecta a RF-01 (login) en su capa de datos.
   - `internal/middleware/cors.go`: **0% de cobertura**, sin test. Este archivo se agregó en esta misma sesión (2026-07-11) para resolver un problema real de despliegue (bloqueo CORS entre Vercel y Railway) y **no se le escribió ningún test** — es una brecha introducida por el propio trabajo de esta sesión, no heredada del curso.
   - `internal/config/config.go` y `internal/database/database.go`: 0%, sin test — son wrappers delgados sobre variables de entorno y `gorm.Open`, de bajo riesgo, pero formalmente están en 0%.

2. **RF-07 (Dashboard IA) — la lógica existe pero NO está expuesta por la API.** `internal/services/ia_service.go` implementa `PredecirDemanda` y `GenerarAlertas` con 100% de cobertura de tests, **pero no existe ningún `ia_handler.go`, y no hay ninguna ruta `/api/admin/ia/prediccion` ni `/api/admin/ia/alertas` registrada en `backend/cmd/server/main.go`.** Verificado explícitamente:
   ```
   $ grep -n "ia" backend/cmd/server/main.go
   (sin resultados relevantes — ninguna ruta de IA registrada)
   ```
   El frontend (`frontend/src/services/api.ts`, objeto `iaApi`) y la página `IAPage.tsx` **sí llaman a esos endpoints**, por lo que en producción esas llamadas devuelven `404 Not Found`. **Esta es una brecha funcional real: RF-07 está probado en aislamiento pero no es alcanzable por ningún usuario del sistema.**

3. **RF-08 (Página Pública) — el formulario de contacto especificado en FC-08 no existe.** FC-08 especifica `POST /api/public/contacto`. Búsqueda exhaustiva:
   ```
   $ grep -rn "contacto" --include="*.go" .
   (sin resultados)
   ```
   No hay handler, no hay ruta, y `LandingPage.tsx` no contiene ningún formulario de contacto (solo teléfono/email/WhatsApp como enlaces estáticos en el footer). **El requisito "formulario de contacto que envía notificación al admin" de RF-08 no está implementado.**

4. **RF-04 (Créditos y Pagos) — falta el endpoint de historial.** FC-04 especifica `GET /api/admin/pagos` ("Historial de pagos con filtros"). El handler real (`creditos_handler.go`) solo expone `ListarDeudores`, `EstadoCuenta` y `RegistrarPago` — **no existe un método para listar todos los pagos**, y no hay ruta `GET /api/admin/pagos` en `main.go`. El historial de pagos solo es visible indirectamente dentro de `EstadoCuenta` (por cliente individual), no como listado general.

5. **RNF-07 ("Deploy automático vía GitHub Actions al hacer push a main") — ya NO se cumple literalmente.** Como se documentó en 5.2, `deploy.yml` fue eliminado. El despliegue sigue siendo automático, pero **vía la integración nativa de Railway/Vercel con GitHub, no vía GitHub Actions** como establece textualmente el RNF. Si la tesis cita este RNF como cumplido, debe matizarse.

6. **RNF-01 (latencia <500ms en 95% de requests) y RNF-06 (50 usuarios concurrentes) — sin evidencia de verificación.** No existe ningún test de carga (no hay `k6`, `locust`, `artillery`, ni script similar en el repositorio). Estos RNF están **declarados pero no verificados** por ninguna prueba automatizada.

7. **Discrepancia de diseño de BD (ver sección 3.2):** el trigger SQL documentado en `PLAN_SDD.md` §5.4 para mantener `deuda_total` sincronizada **no existe en la base de datos real**; la sincronización ocurre manualmente en código Go. Si la tesis presenta el trigger como parte del diseño implementado, es un dato **no verificable en el sistema real** y debe corregirse o aclararse como "diseño alternativo, no trigger de BD".

8. **FC-08 especifica pruebas en Vitest + React Testing Library para el frontend, que no existen.** Búsqueda:
   ```
   $ find frontend/src -iname "*.test.*" -o -iname "*.spec.*"
   (sin resultados)
   ```
   El proyecto **no tiene ninguna prueba automatizada de frontend**, pese a que `PLAN_SDD.md` (tabla de tecnologías, sección 9) declara "Pruebas (React): Vitest + RTL" y FC-08 lista 6 tests específicos (`TestMenuSection_...`, `TestContactoForm_...`) que nunca se escribieron. **Los 208 tests reales del proyecto son 100% backend (Go); no hay cobertura de pruebas en el frontend.**

---

## 6. RESUMEN FINAL — Ciclo completo por Feature Card

| FC | Migración/tabla de BD asociada | Archivos de código principales | Archivo(s) de test | Estado real (verificado, no declarado) |
|---|---|---|---|---|
| **FC-01** Autenticación | `User` (AutoMigrate, sin migración SQL) | `internal/auth/auth.go`, `internal/middleware/jwt.go`, `internal/handlers/auth_handler.go`, `internal/repository/user_repo.go` | `auth_test.go` (14), `jwt_test.go` (8), `auth_handler_test.go` (8) — **`user_repo.go` sin test (0%)** | **Parcial** — lógica y HTTP probados; capa de datos (repository) sin cobertura |
| **FC-02** Clientes | `Cliente` | `internal/handlers/clientes_handler.go`, `internal/repository/clientes_repo.go` | `clientes_handler_test.go` (23), `clientes_repo_test.go` (8) | **Implementado**, con test en las 3 capas |
| **FC-03** Pedidos | `Pedido`, `PedidoItem` | `internal/services/pedido_service.go`, `internal/handlers/pedidos_handler.go`, `internal/repository/pedidos_repo.go` | `pedido_service_test.go` (24), `pedidos_handler_test.go` (16) — **`pedidos_repo.go` sin test (0%)** | **Parcial** — lógica de negocio y HTTP muy bien probados (mocks); capa de datos real (SQL) sin ningún test de integración |
| **FC-04** Créditos y Pagos | `Pago` (en `pedidos_repo.go`, struct `PagoRepo`) | `internal/handlers/creditos_handler.go` | `creditos_handler_test.go` (12) — **`PagoRepo` sin test (0%)** | **Parcial** — falta endpoint `GET /api/admin/pagos` (historial) y test de repositorio |
| **FC-05** Inventario | `Producto`, `MovimientoStock` (en `pedidos_repo.go`, structs `ProductoRepo`/`MovimientoRepo`) | `internal/handlers/inventario_handler.go` | `inventario_handler_test.go` (24) — **`ProductoRepo`/`MovimientoRepo` sin test (0%)** | **Parcial** — handler bien probado; capa de datos real sin ningún test |
| **FC-06** Reportes | Consultas sobre `Pedido`/`Pago`/`MovimientoStock` | `internal/services/reportes_service.go`, `internal/handlers/reportes_handler.go`, `internal/export/export.go` | `reportes_service_test.go` (11), `reportes_handler_test.go` (11), `export_test.go` (7), `reportes_repo_test.go` (3) | **Implementado y bien probado** |
| **FC-07** IA Básica | Lee de `Pedido`, `Cliente`, `Producto` (sin tabla propia) | `internal/services/ia_service.go` | `ia_service_test.go` (8) | **Lógica implementada y probada, pero SIN endpoint HTTP — inalcanzable desde el frontend** |
| **FC-08** Landing Page | `MenuPublico` | `internal/handlers/menu_handler.go`, `internal/repository/menu_repo.go`, `frontend/src/pages/public/LandingPage.tsx` | `menu_handler_test.go` (13), `menu_repo_test.go` (8) — **sin tests de frontend (Vitest/RTL especificados, nunca escritos)** | **Parcial** — menú público funcional y probado en backend; falta formulario de contacto (RF-08) y toda prueba de frontend |
| **FC-09** Reservas *(ampliación, no en el SRS original)* | `Reserva` | `internal/handlers/reservas_handler.go`, `internal/repository/reservas_repo.go`, `frontend/src/pages/public/ReservaPage.tsx`, `frontend/src/pages/dashboard/ReservasPage.tsx` | `reservas_handler_test.go` (17), `reservas_repo_test.go` (8) | **Implementado y probado en las 3 capas** (el módulo más nuevo, pero el más completo en tests) |

**Total de pruebas reales verificadas por `grep -c "^func Test"` en esta auditoría: 208.**
**Cobertura real medida hoy: 82.1% (excluyendo `cmd/*`) — por debajo del 90% que exige RNF-02 y el DoD de cada Feature Card.**

---

## Conclusión para uso en tesis

Este proyecto demuestra un ciclo SDD real y verificable — especificación en Markdown, Feature Cards con criterios de aceptación, TDD parcial, CI configurado — pero **no es un ciclo perfecto**, y presentarlo como tal en la tesis sería inexacto. Los puntos honestos a incluir explícitamente son:
1. La cronología pre-11-julio solo es reconstruible por timestamps de archivo, no por git.
2. FC-09 es una ampliación posterior, no parte del SRS original.
3. La cobertura real hoy (82.1%) no alcanza el 90% objetivo.
4. RF-07 y parte de RF-04/RF-08 tienen brechas de implementación reales y verificables.
5. El trigger de BD documentado no existe; la sincronización es manual en Go.
6. No hay pruebas de frontend pese a estar especificadas.

Estos seis puntos, presentados con la evidencia de este documento, son más valiosos académicamente que afirmar un cumplimiento 100% que el propio repositorio no sostiene.
