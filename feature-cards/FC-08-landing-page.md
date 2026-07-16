# FC-08 — Landing Page Pública
**Estado:** Implementado (con diferencias respecto a la especificación original — ver "Notas de concordancia" al final)
**Sprint:** 4 (según plan) — implementación real dispersa entre el 29-may y el 11-jul de 2026 (ver `AUDITORIA_SDD.md` sección 4.2)
**Prioridad:** Media
**Depende de:** FC-01 (mismo backend), FC-09 (botón "Reservar" enlaza al flujo de reservas)

---

## Descripción
Página web pública del restaurante (`frontend/src/pages/public/LandingPage.tsx`). Sin autenticación. Muestra el menú (desde el backend), información del local, testimonios y datos de contacto estáticos. **No tiene formulario de contacto** — ver Notas de concordancia. El contenido del menú se gestiona desde el dashboard admin (FC-08 → `MenuPage.tsx` / `menu_handler.go`).
Nombre comercial mostrado en la UI: "Vilca Suyo" (aparece en: nav, hero, footer, dashboard, login — ver LandingPage.tsx, LoginPage.tsx, DashboardLayout.tsx).
---

## Secciones reales (verificadas contra `LandingPage.tsx`)

### 1. Nav (barra fija)
- Logo/nombre del restaurante, links de scroll: Heritage · Menú · Experiencia · Ubicación.
- Botón **"Reservar Mesa"** → `/reservar` (FC-09).
- Link **"Admin"** → `/login`.
- Menú hamburguesa en mobile.

### 2. Hero
- Imagen de fondo con parallax, headline, subtítulo, botones **"Reservar Mesa"** y **"Ver Menú"** (scroll a `#menu`).

### 3. Sobre el lugar ("Nuestra Herencia")
- Texto fijo (no editable desde el dashboard) sobre la historia del lugar. Foto + tarjeta "Siglo XV".

### 4. Menú destacado (`id="menu"`)
- Se consume desde `GET /api/public/menu` vía `menuApi.obtener()`.
- **Si el backend no devuelve ítems, usa un menú de respaldo hardcodeado (`FALLBACK_MENU`, 3 platos)** — no está en la especificación original.
- Se muestra como **grilla de tarjetas simple, sin tabs por categoría** (la especificación original decía "Dividido en tabs: Desayunos | Almuerzos | Bebidas | Promociones" — eso **no** está implementado; la categoría solo se muestra como una etiqueta/badge dentro de cada tarjeta).
- `MenuPublico.producto_id` (opcional, `NULL` permitido): vincula un ítem del menú con su `Producto` real de inventario. Se gestiona desde `POST/PUT /api/admin/menu`.

### 5. Reservar mesa (ver FC-09)
- Botón/enlace "Reservar" que lleva a `/reservar` — formulario público independiente del menú, especificado en [FC-09-reservas.md](FC-09-reservas.md).

### 6. Experiencia turística (`id="experiencia"`) — **sección no documentada en la especificación original**
- "Ubicación Privilegiada", "Atención Bilingüe", "Logística de Viaje" — contenido fijo en código, sin backend.

### 7. Testimonios — **sección no documentada en la especificación original**
- 3 testimonios fijos en código (nombre, país, texto, estrellas), sin backend ni gestión desde el dashboard.

### 8. Ubicación (`id="ubicacion"`)
- Desde Ayacucho / Dirección / Horarios (texto fijo) + link a Google Maps.

### 9. CTA final / "Contacto" (`id="reservas"`)
- Botón "Hacer una Reserva" → `/reservar`.
- Teléfono, correo (`hola@vilcasuyo.com`) y WhatsApp como **texto/links estáticos**, no como formulario.

### 10. Footer
- Links decorativos (Heritage, Inca Trail, Private Dining, Sostenibilidad) que **no navegan a ninguna parte** (`href="#"`).
- Íconos sociales, también decorativos (`href="#"`).
- Copyright + link "Admin" → `/login`.

---

## Endpoints públicos reales (sin JWT)

| Método | Ruta                   | Estado | Descripción                          |
|--------|------------------------|--------|---------------------------------------|
| GET    | `/api/public/menu`     | ✅ Implementado | Obtener ítems de menú activos (array plano) |
| POST   | `/api/public/reservas` | ✅ Implementado (ver FC-09) | Crear una reserva |
| POST   | `/api/public/contacto` | ❌ **NO implementado** | Especificado originalmente, nunca se construyó (sin handler, sin ruta, sin UI) |

### GET /api/public/menu — response real (array plano, NO agrupado por categoría)
```json
[
  {
    "id": 1,
    "categoria": "Desayunos",
    "nombre": "Desayuno Completo",
    "descripcion": "Pan, huevo frito, jugo natural y café",
    "precio": 8.00,
    "imagen_url": "",
    "disponible": true,
    "orden": 1,
    "producto_id": null
  }
]
```
*(La especificación original mostraba la respuesta agrupada por categoría (`{"Desayunos": [...], "Almuerzos": [...]}`) — eso no coincide con la implementación real de `MenuRepo.FindPublico()`, que devuelve un array plano ordenado por `orden, nombre`.)*

---

## Stack Frontend real

```
React 19 + TypeScript + Vite
Estilos: mayormente inline (style={{...}}) + clases responsive de TailwindCSS (hidden/md:flex/lg:hidden)
Iconos: lucide-react (SVG) — antes usaba una fuente de Google (Material Symbols),
        reemplazada por fallar en redes inestables (ver commit "Reemplaza iconos de
        Material Symbols... por lucide-react")
React Router v6 (SPA)
```

---

## Pruebas

**No existen pruebas automatizadas de frontend.** La especificación original listaba 6 tests en Vitest + React Testing Library (`TestMenuSection_CargaItemsDelAPI`, `TestContactoForm_...`, etc.) que **nunca se escribieron** — confirmado con `find frontend/src -iname "*.test.*"` sin resultados (ver `AUDITORIA_SDD.md` hallazgo 8). Además, dos de esos tests (`TestContactoForm_*`) prueban un formulario que ni siquiera existe. Esta sección de la especificación original se considera **no ejecutada**, no "pendiente".

---

## Definición de "Hecho" (DoD) — reevaluado contra la implementación real

- [ ] Tests — **no aplica / no hay ninguno que corra** (ver sección Pruebas).
- [x] La página es responsive (usa clases `hidden`/`md:`/`lg:` de Tailwind en nav, footer y layout de secciones).
- [x] El menú se actualiza sin redeployar (viene del backend vía `GET /api/public/menu`), con fallback local si el backend no responde.
- [ ] Formulario de contacto valida campos vacíos — **no aplica, el formulario no existe**.
- [ ] Lighthouse score ≥ 80 en Performance — **sin verificar, no hay evidencia de que se haya corrido nunca**.

---

## Notas de concordancia (spec vs. código real)

Esta sección existe para que la Feature Card sirva como fuente de verdad **actualizada**, no como el plan original sin verificar (ver discusión de trazabilidad en la auditoría SDD):

1. **Contacto:** la especificación pedía un formulario de contacto con `POST /api/public/contacto`. Se implementó en su lugar contacto **estático** (teléfono, correo, WhatsApp como texto/links). Si se requiere el formulario real, es trabajo pendiente, no algo ya hecho parcialmente.
2. **Menú por tabs:** la especificación pedía tabs por categoría; la implementación usa una grilla simple con la categoría como badge. Si se necesita el filtro por tabs, es una mejora pendiente sobre `LandingPage.tsx`.
3. **Secciones no especificadas originalmente:** "Experiencia turística" y "Testimonios" se agregaron durante la implementación sin pasar primero por una actualización de esta Feature Card — quedan documentadas aquí retroactivamente para que la especificación no quede desactualizada.
4. **Formato de respuesta del menú:** la especificación mostraba un JSON agrupado por categoría; el backend real devuelve un array plano. El agrupamiento visual (si lo hay) lo haría el frontend, y actualmente no lo hace.
