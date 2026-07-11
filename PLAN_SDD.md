# PROYECTO FINAL — Sistema de Gestión para Restaurante
**Curso:** Pruebas y Aseguramiento de Calidad de Software  
**Alumno:** Santiago Noe  
**Metodología:** SDD (Spec Driven Development / Design)  
**Fecha:** 2026-05-29  
**Última actualización:** 2026-07-11 — agregado módulo de Reservas (FC-09), inventario real de pruebas (208 tests) y sincronización de la especificación con el código implementado.

---

## ESPECIFICACIÓN DE SDD DE ESTE PROYECTO

De las modalidades posibles para aplicar SDD (Spec Kit con plantillas, OpenSpec, Kiro u otra herramienta dedicada, especificación manual en archivos `.md`, o generación de código asistida por IA), este proyecto combina:

1. **Manual (archivos de configuración `.md`)** — este documento (`PLAN_SDD.md`) y las Feature Cards (`feature-cards/FC-*.md`) se escriben a mano, sin un CLI de spec-driven-development (no se usa Spec Kit, OpenSpec ni Kiro).
2. **IA — Generar código** — una vez escrita la especificación, se usa un asistente de IA para implementar el código y las pruebas que la satisfacen, iterando hasta que los tests definidos en cada Feature Card pasan.

---

## ÍNDICE

1. [Metodología SDD Aplicada](#1-metodología-sdd-aplicada)
2. [Visión General del Sistema](#2-visión-general-del-sistema)
3. [Especificación de Requisitos (SRS)](#3-especificación-de-requisitos-srs)
4. [Historias de Usuario](#4-historias-de-usuario)
5. [Diseño de Base de Datos Normalizada](#5-diseño-de-base-de-datos-normalizada)
6. [Arquitectura del Sistema](#6-arquitectura-del-sistema)
7. [Plan de Pruebas Unitarias (≥90% cobertura)](#7-plan-de-pruebas-unitarias)
8. [Planificación de Entregables](#8-planificación-de-entregables)
9. [Tecnologías (libres)](#9-tecnologías)

---

## 1. METODOLOGÍA SDD APLICADA

Según el pizarrón de clase, el flujo SDD es:

```
Roles:  A (Analista) → D (Diseñador) → I (Implementador)
             ↓
        Artefactos
             ↓
   Planificación / Especificación
             ↓
        .MD  +  F.C. (Feature Cards)
             ↓
        IDE  →  Código
             ↓
        P.U. (Pruebas Unitarias)
             ↓
     Cobertura de Código ≥ 90%
             ↓
  P (Calidad) → D (Desplegar) → Servidor
             ↓
          GITHUB
```

### Qué significa SDD aquí
- **Spec First:** Antes de escribir UNA SOLA línea de código, se define TODO en archivos `.md`.
- **Feature Cards:** Cada funcionalidad se describe como una tarjeta de especificación.
- **Test Driven:** Las pruebas unitarias se diseñan ANTES de implementar (TDD).
- **Cobertura 90%:** Meta mínima de cobertura de código.
- **CI/CD:** GitHub Actions despliega automáticamente al servidor.

---

## 2. VISIÓN GENERAL DEL SISTEMA

### Nombre del proyecto
**RestaurantOS** — Sistema de Gestión Integral para Restaurante

### Descripción
Plataforma web de dos capas:
- **Capa pública:** Landing page del restaurante para clientes.
- **Capa privada:** Dashboard administrativo con autenticación por roles.

### Problema que resuelve
Los restaurantes pequeños y medianos manejan pedidos, créditos e inventario de forma manual (cuadernos, Excel). Esto genera:
- Pérdida de datos
- Errores en cobros
- Sin visibilidad de stock
- Sin reportes de ventas

### Solución propuesta
Un sistema web accesible desde cualquier dispositivo que centralice pedidos, pagos, créditos, inventario y reportes, con alertas IA básicas.

---

## 3. ESPECIFICACIÓN DE REQUISITOS (SRS)

### 3.1 Requisitos Funcionales

#### RF-01 — Autenticación y Roles
- El sistema debe permitir login con email y contraseña.
- Debe generar un JWT con expiración de 8 horas.
- Debe soportar dos roles: `admin` y `empleado`.
- El rol `admin` tiene acceso total.
- El rol `empleado` puede registrar pedidos y ver clientes; NO puede modificar configuraciones ni ver reportes financieros.

#### RF-02 — Gestión de Clientes
- Registrar cliente con: nombre, apellido, tipo (individual/empresa), teléfono, dirección, email.
- Listar, buscar, editar y desactivar clientes.
- Ver historial de pedidos y estado de cuenta por cliente.

#### RF-03 — Gestión de Pedidos
- Registrar pedido con: cliente, fecha, tipo (desayuno/almuerzo), lista de ítems, cantidad, precio unitario.
- Calcular total automáticamente.
- Estados del pedido: `pendiente`, `entregado`, `anulado`.
- Un pedido puede pagarse al contado o a crédito.

#### RF-04 — Créditos y Pagos
- Si un pedido es a crédito, se crea una deuda en la cuenta del cliente.
- Registrar abonos parciales o totales.
- Ver saldo pendiente por cliente en tiempo real.
- Alertar cuando un cliente supere S/ 200 de deuda.

#### RF-05 — Inventario
- Registrar productos con: nombre, unidad (kg, litro, unidad), stock actual, stock mínimo.
- Descontar stock automáticamente al registrar un pedido.
- Alertar cuando el stock de un producto baje del mínimo.

#### RF-06 — Reportes
- Reporte de ventas diario, semanal y mensual.
- Reporte de clientes con deuda.
- Reporte de movimientos de inventario.
- Exportar en PDF y Excel.

#### RF-07 — Dashboard IA (básico)
- Mostrar predicción de demanda del día siguiente basada en promedio histórico de los últimos 30 días.
- Alerta de cliente con deuda alta (> umbral configurable).
- Sugerencia de reabastecimiento para productos con tendencia baja.

#### RF-08 — Página Pública
- Mostrar información del restaurante: presentación, menú, nosotros, contacto.
- El menú debe actualizarse desde el dashboard (mismo backend).
- Formulario de contacto que envía notificación al admin.

### 3.2 Requisitos No Funcionales

| ID     | Requisito                                                              |
|--------|------------------------------------------------------------------------|
| RNF-01 | Tiempo de respuesta de la API < 500ms en el 95% de las peticiones      |
| RNF-02 | Cobertura de pruebas unitarias ≥ 90%                                   |
| RNF-03 | Autenticación con JWT firmado (HS256)                                  |
| RNF-04 | Contraseñas almacenadas con bcrypt (costo 12)                          |
| RNF-05 | La API debe seguir el estándar REST                                    |
| RNF-06 | El sistema debe soportar al menos 50 usuarios concurrentes             |
| RNF-07 | Deploy automático vía GitHub Actions al hacer push a `main`            |
| RNF-08 | Toda la comunicación por HTTPS                                         |
| RNF-09 | Logs de errores persistentes en servidor                               |
| RNF-10 | Base de datos con backups diarios automáticos                          |

---

## 4. HISTORIAS DE USUARIO

> Formato: **Como** [rol], **quiero** [acción], **para** [beneficio].

### HU-01 — Login
**Como** administrador,  
**quiero** iniciar sesión con email y contraseña,  
**para** acceder al dashboard de forma segura.

**Criterios de aceptación:**
- [ ] Si las credenciales son correctas → recibo un JWT y entro al dashboard.
- [ ] Si las credenciales son incorrectas → veo el mensaje "Credenciales inválidas".
- [ ] El token expira en 8 horas.
- [ ] Las rutas del dashboard son inaccesibles sin token válido.

---

### HU-02 — Registrar pedido
**Como** empleado,  
**quiero** registrar el pedido de un cliente,  
**para** tener un registro formal de lo consumido.

**Criterios de aceptación:**
- [ ] Puedo buscar un cliente por nombre.
- [ ] Puedo agregar múltiples ítems al pedido.
- [ ] El total se calcula automáticamente.
- [ ] Elijo si el pago es contado o crédito.
- [ ] Al guardar, el stock de ingredientes se descuenta.

---

### HU-03 — Ver estado de cuenta
**Como** administrador,  
**quiero** ver cuánto debe cada cliente,  
**para** cobrar correctamente y evitar pérdidas.

**Criterios de aceptación:**
- [ ] Veo la lista de clientes con deuda activa.
- [ ] Al hacer clic en un cliente veo el detalle de cada pedido pendiente.
- [ ] Puedo registrar un abono y el saldo se actualiza en tiempo real.

---

### HU-04 — Control de inventario
**Como** administrador,  
**quiero** ver el stock de cada ingrediente,  
**para** planificar las compras sin quedarse sin insumos.

**Criterios de aceptación:**
- [ ] Veo todos los productos con su stock actual.
- [ ] Los productos con stock bajo aparecen resaltados en rojo.
- [ ] Puedo actualizar manualmente el stock (reabastecimiento).

---

### HU-05 — Predicción IA
**Como** administrador,  
**quiero** ver una predicción de la demanda del día siguiente,  
**para** preparar la cantidad correcta de comida.

**Criterios de aceptación:**
- [ ] El dashboard muestra "Plato más pedido mañana: Almuerzo especial (est. 45 porciones)".
- [ ] La predicción se basa en el promedio de los últimos 30 días para ese día de la semana.

---

### HU-06 — Reporte de ventas
**Como** administrador,  
**quiero** generar un reporte semanal de ventas,  
**para** conocer los ingresos y tomar decisiones.

**Criterios de aceptación:**
- [ ] Elijo rango de fechas y el sistema genera el reporte.
- [ ] Puedo exportarlo como PDF.
- [ ] Muestra: total ventas, ventas contado, ventas crédito, plato más vendido.

---

## 5. DISEÑO DE BASE DE DATOS NORMALIZADA

### 5.1 Análisis de Normalización

**Objetivo:** Llevar el modelo a 3FN (Tercera Forma Normal).

- **1FN:** Todos los atributos son atómicos, sin grupos repetitivos.
- **2FN:** Sin dependencias parciales (cada atributo depende de toda la PK).
- **3FN:** Sin dependencias transitivas (atributos no clave no dependen entre sí).

---

### 5.2 Diagrama de Entidades

```
┌──────────┐     ┌──────────────┐     ┌────────────┐
│  users   │     │   clientes   │     │  pedidos   │
├──────────┤     ├──────────────┤     ├────────────┤
│ id (PK)  │     │ id (PK)      │     │ id (PK)    │
│ nombre   │     │ nombre       │     │ cliente_id │──→ clientes.id
│ apellido │     │ apellido     │     │ user_id    │──→ users.id
│ email    │     │ tipo         │     │ fecha      │
│ password │     │ telefono     │     │ tipo_comida│
│ rol      │     │ direccion    │     │ estado     │
│ activo   │     │ email        │     │ forma_pago │
│ created_at│     │ deuda_total  │     │ total      │
└──────────┘     │ activo       │     │ created_at │
                 │ created_at   │     └────────────┘
                 └──────────────┘           │
                                            │ 1:N
                 ┌──────────────┐     ┌─────▼──────────┐
                 │  productos   │     │  pedido_items  │
                 ├──────────────┤     ├────────────────┤
                 │ id (PK)      │     │ id (PK)        │
                 │ nombre       │     │ pedido_id      │──→ pedidos.id
                 │ unidad       │     │ producto_id    │──→ productos.id
                 │ stock_actual │     │ cantidad       │
                 │ stock_minimo │     │ precio_unitario│
                 │ precio_venta │     │ subtotal       │
                 │ activo       │     └────────────────┘
                 └──────────────┘

┌──────────────────┐     ┌─────────────────────┐
│      pagos       │     │  movimientos_stock  │
├──────────────────┤     ├─────────────────────┤
│ id (PK)          │     │ id (PK)             │
│ cliente_id       │──→  │ producto_id         │──→ productos.id
│ pedido_id (NULL) │     │ tipo (entrada/salida│
│ monto            │     │ cantidad            │
│ metodo           │     │ referencia_id       │
│ fecha            │     │ fecha               │
│ created_at       │     │ created_at          │
└──────────────────┘     └─────────────────────┘

┌───────────────────────┐     ┌───────────────────────┐
│      menu_publico      │     │        reservas        │
├───────────────────────┤     ├───────────────────────┤
│ id (PK)                │     │ id (PK)                │
│ categoria               │     │ nombre                  │
│ nombre                  │     │ whatsapp                │
│ descripcion             │     │ fecha                   │
│ precio                  │     │ personas                │
│ imagen_url              │     │ ocasion                 │
│ disponible              │     │ estado                  │
│ orden                   │     │ created_at              │
│ producto_id (FK, NULL)  │     │ pedido_id (FK, NULL)    │
└──────────┬──────────────┘     └──────────┬──────────────┘
           │                                │
           ▼                                ▼
      productos.id                      pedidos.id
```

> Ambas relaciones son **FK opcionales (NULL permitido)**, no obligatorias:
> - `menu_publico.producto_id` → `productos.id`: vincula un ítem del menú público con el producto real de inventario que representa (por ejemplo, para descontar stock automáticamente si en el futuro se vende directo desde el menú). Puede ser NULL para ítems que son solo promocionales/informativos y no mapean 1:1 a un producto de inventario.
> - `reservas.pedido_id` → `pedidos.id`: se completa cuando el cliente llega y el mesero registra el `Pedido` real — el endpoint `PUT /api/empleado/reservas/:id/pedido` hace ese enlace. Nace en NULL porque una reserva puede cancelarse o no concretarse nunca en una visita real; la reserva **no controla stock ni totales**, eso lo sigue haciendo `Pedido`.
>
> Con esto, las 9 tablas del modelo quedan conectadas entre sí (ninguna tabla aislada).

---

### 5.3 Tablas en SQL (PostgreSQL)

```sql
-- Tabla: users (autenticación y roles)
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(100) NOT NULL,
    apellido    VARCHAR(100) NOT NULL,
    email       VARCHAR(150) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,           -- bcrypt hash
    rol         VARCHAR(20) NOT NULL CHECK (rol IN ('admin', 'empleado')),
    activo      BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Tabla: clientes
CREATE TABLE clientes (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(100) NOT NULL,
    apellido    VARCHAR(100),
    tipo        VARCHAR(20) NOT NULL CHECK (tipo IN ('individual', 'empresa')),
    telefono    VARCHAR(20),
    direccion   TEXT,
    email       VARCHAR(150),
    deuda_total DECIMAL(10,2) DEFAULT 0.00,      -- campo calculado, se actualiza con trigger
    activo      BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Tabla: productos (inventario)
CREATE TABLE productos (
    id           SERIAL PRIMARY KEY,
    nombre       VARCHAR(150) NOT NULL,
    unidad       VARCHAR(20) NOT NULL CHECK (unidad IN ('kg', 'litro', 'unidad', 'porcion')),
    stock_actual DECIMAL(10,3) DEFAULT 0,
    stock_minimo DECIMAL(10,3) DEFAULT 0,
    precio_venta DECIMAL(10,2) DEFAULT 0.00,
    activo       BOOLEAN DEFAULT TRUE,
    created_at   TIMESTAMP DEFAULT NOW()
);

-- Tabla: pedidos
CREATE TABLE pedidos (
    id          SERIAL PRIMARY KEY,
    cliente_id  INT NOT NULL REFERENCES clientes(id),
    user_id     INT NOT NULL REFERENCES users(id),  -- empleado que registró
    fecha       DATE NOT NULL DEFAULT CURRENT_DATE,
    tipo_comida VARCHAR(30) NOT NULL CHECK (tipo_comida IN ('desayuno', 'almuerzo', 'cena', 'bebida', 'otro')),
    estado      VARCHAR(20) NOT NULL DEFAULT 'pendiente' CHECK (estado IN ('pendiente', 'entregado', 'anulado')),
    forma_pago  VARCHAR(20) NOT NULL CHECK (forma_pago IN ('contado', 'credito')),
    total       DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    notas       TEXT,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Tabla: pedido_items (detalle de cada pedido — resuelve 2FN)
CREATE TABLE pedido_items (
    id              SERIAL PRIMARY KEY,
    pedido_id       INT NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    producto_id     INT NOT NULL REFERENCES productos(id),
    cantidad        DECIMAL(10,3) NOT NULL,
    precio_unitario DECIMAL(10,2) NOT NULL,
    subtotal        DECIMAL(10,2) GENERATED ALWAYS AS (cantidad * precio_unitario) STORED
);

-- Tabla: pagos
CREATE TABLE pagos (
    id          SERIAL PRIMARY KEY,
    cliente_id  INT NOT NULL REFERENCES clientes(id),
    pedido_id   INT REFERENCES pedidos(id),    -- NULL si es abono general
    monto       DECIMAL(10,2) NOT NULL,
    metodo      VARCHAR(20) NOT NULL CHECK (metodo IN ('efectivo', 'transferencia', 'yape', 'plin', 'otro')),
    fecha       DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Tabla: movimientos_stock (trazabilidad del inventario)
CREATE TABLE movimientos_stock (
    id            SERIAL PRIMARY KEY,
    producto_id   INT NOT NULL REFERENCES productos(id),
    tipo          VARCHAR(10) NOT NULL CHECK (tipo IN ('entrada', 'salida')),
    cantidad      DECIMAL(10,3) NOT NULL,
    referencia_id INT,                          -- pedido_id si es salida automática
    notas         TEXT,
    fecha         DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at    TIMESTAMP DEFAULT NOW()
);

-- Tabla: menu_publico (datos para la landing page)
CREATE TABLE menu_publico (
    id          SERIAL PRIMARY KEY,
    categoria   VARCHAR(50) NOT NULL CHECK (categoria IN ('Desayunos', 'Almuerzos', 'Bebidas', 'Promociones')),
    nombre      VARCHAR(150) NOT NULL,
    descripcion TEXT,
    precio      DECIMAL(10,2),
    imagen_url  VARCHAR(500),
    disponible  BOOLEAN DEFAULT TRUE,
    orden       INT DEFAULT 0,
    producto_id INT REFERENCES productos(id)   -- opcional: vincula el ítem del menú con el producto real de inventario
);

-- Tabla: reservas (intención de visita desde la landing page, sin login)
CREATE TABLE reservas (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(100) NOT NULL,
    whatsapp    VARCHAR(20) NOT NULL,
    fecha       DATE NOT NULL,
    personas    VARCHAR(10) NOT NULL,
    ocasion     VARCHAR(50),
    estado      VARCHAR(20) NOT NULL DEFAULT 'pendiente' CHECK (estado IN ('pendiente', 'confirmada', 'cancelada')),
    pedido_id   INT REFERENCES pedidos(id),    -- opcional: se completa cuando la reserva se concreta en un pedido real
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Índices para consultas frecuentes
CREATE INDEX idx_pedidos_cliente   ON pedidos(cliente_id);
CREATE INDEX idx_pedidos_fecha     ON pedidos(fecha);
CREATE INDEX idx_pagos_cliente     ON pagos(cliente_id);
CREATE INDEX idx_movimientos_prod  ON movimientos_stock(producto_id);
CREATE INDEX idx_menu_producto     ON menu_publico(producto_id);
CREATE INDEX idx_reservas_pedido   ON reservas(pedido_id);
```

---

### 5.4 Trigger — Actualizar deuda_total del cliente

```sql
-- Función que actualiza deuda_total en clientes
CREATE OR REPLACE FUNCTION actualizar_deuda_cliente()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE clientes
    SET deuda_total = (
        SELECT COALESCE(SUM(p.total), 0)
        FROM pedidos p
        WHERE p.cliente_id = COALESCE(NEW.cliente_id, OLD.cliente_id)
          AND p.forma_pago = 'credito'
          AND p.estado = 'entregado'
    ) - (
        SELECT COALESCE(SUM(pg.monto), 0)
        FROM pagos pg
        WHERE pg.cliente_id = COALESCE(NEW.cliente_id, OLD.cliente_id)
    )
    WHERE id = COALESCE(NEW.cliente_id, OLD.cliente_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_deuda_pedido
AFTER INSERT OR UPDATE ON pedidos
FOR EACH ROW EXECUTE FUNCTION actualizar_deuda_cliente();

CREATE TRIGGER trg_deuda_pago
AFTER INSERT ON pagos
FOR EACH ROW EXECUTE FUNCTION actualizar_deuda_cliente();
```

---

## 6. ARQUITECTURA DEL SISTEMA

```
┌────────────────────────────────────────────────────┐
│                   CLIENTE (Browser)                │
│                                                    │
│  ┌──────────────────┐    ┌─────────────────────┐   │
│  │  Landing Page    │    │  Dashboard Admin    │   │
│  │  React + Tailwind│    │  React + Tailwind   │   │
│  │  (público)       │    │  (privado, JWT)     │   │
│  └────────┬─────────┘    └──────────┬──────────┘   │
└───────────┼──────────────────────── ┼──────────────┘
            │          HTTPS          │
            └────────────┬────────────┘
                         ↓
┌────────────────────────────────────────────────────┐
│                 API REST (Go + Gin)                │
│                                                    │
│  /api/public/*    →  sin auth                     │
│  /api/auth/*      →  login, refresh               │
│  /api/admin/*     →  JWT required (admin)         │
│  /api/empleado/*  →  JWT required (any role)      │
│                                                    │
│  Middleware: JWT → RBAC → Rate Limiting → Logger   │
└───────────────────────┬────────────────────────────┘
                        ↓
┌────────────────────────────────────────────────────┐
│              PostgreSQL (GORM ORM)                 │
│                                                    │
│  users | clientes | pedidos | pedido_items         │
│  pagos | productos | movimientos_stock             │
│  menu_publico | reservas                           │
└────────────────────────────────────────────────────┘
```

### Estructura de carpetas del proyecto (estado actual)

```
restaurantos/
├── backend/                   # Go API
│   ├── cmd/
│   │   ├── server/main.go     # Wiring de rutas Gin
│   │   ├── migrate/main.go    # AutoMigrate standalone
│   │   ├── seed/               # Datos de desarrollo
│   │   └── genhash/            # Utilidad bcrypt
│   ├── internal/
│   │   ├── auth/               # JWT, bcrypt (auth.go + auth_test.go)
│   │   ├── middleware/         # JWT + RBAC (jwt.go + jwt_test.go)
│   │   ├── models/             # models.go — todas las entidades GORM
│   │   ├── config/              # carga de variables de entorno
│   │   ├── database/            # Connect() + Migrate() (AutoMigrate)
│   │   ├── handlers/            # HTTP handlers (Gin) — 1 archivo + test por módulo:
│   │   │   auth_handler(.go/_test.go), clientes_handler, pedidos_handler,
│   │   │   creditos_handler, inventario_handler, menu_handler,
│   │   │   reportes_handler, reservas_handler
│   │   ├── repository/          # Acceso a DB (GORM) — testeado con Postgres real:
│   │   │   clientes_repo, pedidos_repo, menu_repo, reservas_repo,
│   │   │   user_repo, testdb_test.go (helper de transacciones)
│   │   ├── services/            # Lógica de negocio pura/orquestación:
│   │   │   pedido_service, ia_service, reportes_service
│   │   └── export/               # Generación de PDF/Excel (gofpdf, excelize)
│   ├── go.mod / go.sum
│
├── frontend/                   # React + TypeScript + Vite + Tailwind
│   ├── src/
│   │   ├── pages/
│   │   │   ├── public/          # Sin autenticación
│   │   │   │   ├── LandingPage.tsx
│   │   │   │   ├── LoginPage.tsx
│   │   │   │   └── ReservaPage.tsx     # Formulario público → POST /api/public/reservas
│   │   │   └── dashboard/        # Requiere JWT (PrivateRoute)
│   │   │       ├── DashboardHome.tsx
│   │   │       ├── ClientesPage.tsx
│   │   │       ├── PedidosPage.tsx
│   │   │       ├── ReservasPage.tsx     # Gestión: listar/filtrar/confirmar/cancelar
│   │   │       ├── InventarioPage.tsx
│   │   │       ├── CreditosPage.tsx
│   │   │       ├── IAPage.tsx
│   │   │       ├── MenuPage.tsx
│   │   │       └── ReportesPage.tsx
│   │   ├── components/layout/    # DashboardLayout, PrivateRoute
│   │   ├── context/               # AuthContext (JWT en localStorage)
│   │   ├── services/api.ts        # Cliente Axios + un objeto *Api por módulo
│   │   ├── types/index.ts         # Interfaces TS espejo de los modelos Go
│   │   └── App.tsx                # Rutas (react-router-dom)
│   ├── package.json / vite.config.ts
│
├── database/
│   └── seed_dev.sql
│
├── .github/workflows/
│   ├── test.yml               # Tests + cobertura ≥90% en cada push/PR
│   └── deploy.yml             # Deploy a Railway tras merge a main
│
├── feature-cards/              # FC-01 a FC-09 (una por módulo)
├── docker-compose.yml
└── README.md
```

---

## 7. PLAN DE PRUEBAS UNITARIAS

**Meta:** ≥ 90% de cobertura de código en el backend (Go).

### 7.1 Áreas a probar

| Módulo              | Función/Método                    | Tipo de Test           |
|---------------------|-----------------------------------|------------------------|
| auth                | `HashPassword`                    | Unitario               |
| auth                | `CheckPassword`                   | Unitario               |
| auth                | `GenerateJWT`                     | Unitario               |
| auth                | `ValidateJWT`                     | Unitario               |
| middleware          | `JWTMiddleware` — token válido    | Unitario               |
| middleware          | `JWTMiddleware` — token expirado  | Unitario               |
| middleware          | `JWTMiddleware` — sin token       | Unitario               |
| services/pedido     | `CalcularTotal`                   | Unitario               |
| services/pedido     | `CrearPedidoCredito`              | Integración (DB test)  |
| services/pedido     | `AnularPedido` — actualiza stock  | Integración            |
| services/ia         | `PredecirDemanda` — sin historial | Unitario (edge case)   |
| services/ia         | `PredecirDemanda` — 30 días datos | Unitario               |
| repository/clientes | `GetClienteConDeuda`              | Integración (DB test)  |
| handlers/pedidos    | `POST /pedidos` — body inválido   | Unitario (httptest)    |
| handlers/pedidos    | `POST /pedidos` — body válido     | Integración            |
| handlers/reservas   | `POST /public/reservas` — válido/inválido/fecha mala | Unitario (httptest, mock repo) |
| handlers/reservas   | `PUT /reservas/:id/estado` — estado inválido/inexistente | Unitario |
| repository/reservas | `Create`, `FindAll` (filtro por estado), `UpdateEstado` | Integración (DB test) |

### 7.1.1 Inventario real de pruebas (backend)

Conteo de funciones `func Test...` por archivo, a la fecha de esta actualización (2026-07-11):

| Paquete      | Archivo                          | Tests | Tipo                          |
|--------------|-----------------------------------|:-----:|-------------------------------|
| auth         | `auth_test.go`                    |  14   | Unitario                       |
| middleware   | `jwt_test.go`                     |   8   | Unitario                       |
| handlers     | `auth_handler_test.go`            |   8   | Unitario (mock repo)           |
| handlers     | `clientes_handler_test.go`        |  23   | Unitario (mock repo)           |
| handlers     | `pedidos_handler_test.go`         |  16   | Unitario (mock service)        |
| handlers     | `creditos_handler_test.go`        |  12   | Unitario (mock repo)           |
| handlers     | `inventario_handler_test.go`      |  24   | Unitario (mock repo)           |
| handlers     | `menu_handler_test.go`            |  13   | Unitario (mock repo)           |
| handlers     | `reportes_handler_test.go`        |  11   | Unitario (mock repo)           |
| handlers     | `reservas_handler_test.go`        |  17   | Unitario (mock repo)           |
| repository   | `clientes_repo_test.go`           |   8   | Integración (Postgres real)    |
| repository   | `menu_repo_test.go`               |   8   | Integración (Postgres real)    |
| repository   | `reportes_repo_test.go`           |   3   | Integración (Postgres real)    |
| repository   | `reservas_repo_test.go`           |   8   | Integración (Postgres real)    |
| services     | `pedido_service_test.go`          |  24   | Unitario (mock repo)           |
| services     | `ia_service_test.go`               |   8   | Unitario                       |
| services     | `reportes_service_test.go`        |  11   | Unitario                       |
| **Total**    |                                    | **216** |                              |

> Los tests de `repository/*` corren contra un Postgres real (contenedor `restaurantos_test_db`, puerto 5433 en local; servicio `postgres` efímero en CI) usando una transacción por test que se revierte al final ([testdb_test.go](../backend/internal/repository/testdb_test.go)) — no dejan datos residuales. Los de `handlers/*` y `services/*` usan mocks (`testify/mock`) y no requieren base de datos.
>
> Cobertura exacta: correr `go test ./... -coverprofile=coverage.out` (sección 7.2) — el pipeline de CI la calcula y bloquea el merge si baja de 90% (ver 7.3).

### 7.2 Comando para medir cobertura

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out | grep total
```

### 7.3 GitHub Actions — test.yml (contenido real vigente)

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
          # Los paquetes cmd/* son solo composicion/wiring (main.go) y no se
          # prueban unitariamente; se excluyen del calculo para que el umbral
          # del 90% (RNF-02) refleje la cobertura real de internal/*.
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

---

## 8. PLANIFICACIÓN DE ENTREGABLES

> Estado actualizado al 2026-07-11. Los checkboxes reflejan lo verificado en el código actual del repositorio, no solo lo planeado.

### Sprint 0 — Especificación
- [x] PLAN_SDD.md — este documento
- [x] FC-01.md — Feature Card: Autenticación
- [x] FC-02.md — Feature Card: Clientes
- [x] FC-03.md — Feature Card: Pedidos
- [x] FC-04.md — Feature Card: Créditos y Pagos
- [x] FC-05.md — Feature Card: Inventario
- [x] FC-06.md — Feature Card: Reportes
- [x] FC-07.md — Feature Card: IA Básica
- [x] FC-08.md — Feature Card: Landing Page
- [x] FC-09.md — Feature Card: Reservas (landing → agenda del restaurante)

### Sprint 1 — Base del proyecto
- [x] Setup Go + Gin + GORM + PostgreSQL
- [x] Migraciones de base de datos (`database.Migrate` con AutoMigrate)
- [x] Módulo de autenticación con JWT (`internal/auth`, `internal/middleware`)
- [x] Pruebas unitarias de auth (14 tests en `auth_test.go` + 8 en `jwt_test.go`)

### Sprint 2 — Núcleo del negocio
- [x] CRUD de clientes (`clientes_handler.go` + `clientes_repo.go`)
- [x] Registro de pedidos (`pedido_service.go` — cálculo de total, descuento de stock, anulación)
- [x] Sistema de créditos y pagos (`creditos_handler.go`)
- [x] Pruebas unitarias y de integración (ver inventario completo en 7.1.1)

### Sprint 3 — Inventario + IA + Reportes
- [x] Módulo de inventario (`inventario_handler.go`, alertas de stock mínimo)
- [x] Predicción de demanda (IA básica) (`ia_service.go`)
- [x] Generación de reportes PDF/Excel (`internal/export`, `gofpdf` + `excelize`)
- [x] Pruebas unitarias

### Sprint 4 — Frontend
- [x] Landing page pública (React + Tailwind) — `LandingPage.tsx`
- [x] Dashboard administrativo — `DashboardLayout.tsx` + páginas por módulo
- [x] Integración con API — `services/api.ts` (Axios + interceptores JWT)

### Sprint 5 — CI/CD y Deploy
- [x] GitHub Actions — test automático (`.github/workflows/test.yml`)
- [x] GitHub Actions — deploy automático (`.github/workflows/deploy.yml`)
- [x] Deploy en Railway (`Procfile` en la raíz)
- [ ] Revisión final de cobertura ≥ 90% — pendiente de re-ejecutar tras agregar Reservas (correr `go test ./... -coverprofile=coverage.out`)

### Sprint 6 — Ampliación: Reservas desde la landing page
- [x] Modelo `Reserva` + migración (`models.go`, `database.go`)
- [x] Endpoint público `POST /api/public/reservas` (sin login)
- [x] Endpoints internos `GET /api/empleado/reservas` y `PUT /api/empleado/reservas/:id/estado`
- [x] `ReservaPage.tsx` conectada al backend real (antes simulaba el envío)
- [x] `ReservasPage.tsx` en el dashboard: listar, filtrar por estado, confirmar/cancelar
- [x] Pruebas: 11 unitarias (handler, mock repo) + 6 de integración (repo, Postgres real)
- [ ] Confirmar telefónica manual como filtro anti-broma (proceso operativo, no requiere código adicional)

---

## 9. TECNOLOGÍAS

| Capa              | Tecnología       | Justificación                                      |
|-------------------|------------------|----------------------------------------------------|
| Frontend          | React + TypeScript | Componentes reutilizables, tipado fuerte          |
| Estilos           | TailwindCSS      | Desarrollo rápido, responsive por defecto          |
| Backend           | Go + Gin         | Alta performance, compilado, ideal para APIs REST  |
| ORM               | GORM             | Migraciones y queries seguras en Go                |
| Autenticación     | JWT (HS256)      | Estándar sin estado, fácil de escalar              |
| Base de datos     | PostgreSQL       | Robusto, soporte para triggers, transacciones ACID |
| Pruebas (Go)      | testing + testify| Librería estándar + assertions claros              |
| Pruebas (React)   | Vitest + RTL     | Integración nativa con Vite                        |
| CI/CD             | GitHub Actions   | Integrado con GitHub, gratuito                     |
| Deploy            | Railway          | Soporta Go, PostgreSQL, Docker, fácil de configurar|
| Contenedores      | Docker Compose   | Entorno local reproducible                         |

---

## PRÓXIMO PASO

Con esta especificación lista, el siguiente archivo a crear es la primera **Feature Card**:

> `FC-01.md` — Autenticación: Login, JWT, Roles

Eso define exactamente qué pruebas unitarias escribir ANTES de implementar el código.

**Esto es SDD: Spec → Tests → Código. Nunca al revés.**
