# FC-08 — Landing Page Pública
**Estado:** Especificado  
**Sprint:** 4  
**Prioridad:** Media  
**Depende de:** FC-01 (usa el mismo backend, endpoint público)

---

## Descripción
Página web pública del restaurante. Sin autenticación. Muestra el menú, información del local y formulario de contacto. El contenido del menú se gestiona desde el dashboard admin.

---

## Secciones

### 1. Hero / Inicio
- Nombre del restaurante, slogan, imagen principal.
- Botón "Ver Menú" (scroll a sección menú).

### 2. Menú
- Se consume desde `GET /api/public/menu`.
- Dividido en tabs: Desayunos | Almuerzos | Bebidas | Promociones.
- Cada ítem muestra: nombre, descripción, precio, imagen.
- `MenuPublico.producto_id` (opcional, `NULL` permitido): vincula un ítem del menú con su `Producto` real de inventario, cuando corresponde uno a uno. Se gestiona desde `POST/PUT /api/admin/menu` (campo `producto_id`); los ítems puramente promocionales pueden dejarlo vacío.

### 2.1 Reservar mesa (ver FC-09)
- Botón/enlace "Reservar" que lleva a `/reservar` — formulario público independiente del menú, especificado en detalle en [FC-09-reservas.md](FC-09-reservas.md).

### 3. Nosotros
- Historia y misión del restaurante (texto fijo, editable en código).

### 4. Contacto
- WhatsApp (link directo `wa.me/...`).
- Dirección con mapa (Google Maps embed o enlace).
- Horarios de atención.
- Formulario simple → `POST /api/public/contacto`.

---

## Endpoints públicos (sin JWT)

| Método | Ruta                   | Descripción                          |
|--------|------------------------|--------------------------------------|
| GET    | /api/public/menu       | Obtener menú activo por categoría    |
| POST   | /api/public/contacto   | Enviar mensaje de contacto al admin  |

### GET /api/public/menu response:
```json
{
  "Desayunos": [
    {
      "nombre": "Desayuno completo",
      "descripcion": "Pan, huevo, jugo natural",
      "precio": 8.00,
      "imagen_url": "/images/desayuno.jpg"
    }
  ],
  "Almuerzos": [...],
  "Bebidas": [...],
  "Promociones": [...]
}
```

### POST /api/public/contacto body:
```json
{
  "nombre": "María López",
  "telefono": "987654321",
  "mensaje": "¿Hacen delivery?"
}
```

---

## Stack Frontend

```
React + TypeScript + Vite
TailwindCSS (responsive: mobile first)
React Router v6 (SPA)
```

---

## Pruebas

Las pruebas de la landing son en el **frontend** (Vitest + React Testing Library):

```
TestMenuSection_CargaItemsDelAPI
TestMenuSection_TabsFiltraPorCategoria
TestMenuSection_MuestraSkeletonMientrasCargar
TestContactoForm_EnviaDatosCorrectos
TestContactoForm_SinNombreNoEnvia
TestContactoForm_MuestraMensajeExitoAlEnviar
```

---

## Definición de "Hecho" (DoD)
- [ ] Todos los tests pasan.
- [ ] La página es responsive (mobile, tablet, desktop).
- [ ] El menú se actualiza sin redeployar (viene del backend).
- [ ] El formulario de contacto valida campos vacíos.
- [ ] Lighthouse score ≥ 80 en Performance.
