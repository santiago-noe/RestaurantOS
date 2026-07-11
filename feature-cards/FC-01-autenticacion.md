# FC-01 — Autenticación y Roles
**Estado:** Especificado  
**Sprint:** 1  
**Prioridad:** Alta (bloqueante para todo lo demás)

---

## Descripción
El sistema necesita identificar quién está usando el dashboard y qué puede hacer. Se implementa login con email/contraseña, generación de JWT y verificación de rol en cada ruta protegida.

---

## Endpoints

| Método | Ruta            | Auth | Rol   | Descripción               |
|--------|-----------------|------|-------|---------------------------|
| POST   | /api/auth/login | No   | —     | Login, devuelve JWT       |
| GET    | /api/auth/me    | Sí   | any   | Datos del usuario actual  |

---

## Contrato de la API

### POST /api/auth/login

**Request body:**
```json
{
  "email": "admin@restaurante.com",
  "password": "miContraseña123"
}
```

**Response 200 OK:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "nombre": "María",
    "apellido": "García",
    "email": "admin@restaurante.com",
    "rol": "admin"
  }
}
```

**Response 401 Unauthorized:**
```json
{
  "error": "Credenciales inválidas"
}
```

---

## Pruebas Unitarias (escribir ANTES del código)

### auth_test.go

```
TestHashPassword_GeneraHashDiferenteAlOriginal
TestHashPassword_MismaContraseñaGeneraHashesDistintos
TestCheckPassword_ContraseñaCorrectaRetornaTrue
TestCheckPassword_ContraseñaIncorrectaRetornaFalse
TestGenerateJWT_TokenContieneUserID
TestGenerateJWT_TokenContieneRol
TestGenerateJWT_TokenExpiraEn8Horas
TestValidateJWT_TokenValidoDevuelveClaimsCorrectos
TestValidateJWT_TokenExpiradoDevuelveError
TestValidateJWT_TokenMalformadoDevuelveError
TestValidateJWT_TokenConFirmaFalsaDevuelveError
```

### middleware_test.go

```
TestJWTMiddleware_SinHeaderDevuelve401
TestJWTMiddleware_TokenInvalidoDevuelve401
TestJWTMiddleware_TokenValidoPasaAlSiguiente
TestJWTMiddleware_RolAdminAccedeARutaAdmin
TestJWTMiddleware_RolEmpleadoNoAccedeARutaAdmin
```

---

## Reglas de negocio
1. La contraseña se almacena con bcrypt costo 12. Nunca en texto plano.
2. El JWT contiene: `user_id`, `email`, `rol`, `exp` (8h).
3. El JWT se firma con `JWT_SECRET` desde variable de entorno.
4. Ruta `/api/admin/*` solo accesible con rol `admin`.
5. Ruta `/api/empleado/*` accesible con rol `admin` o `empleado`.

---

## Definición de "Hecho" (DoD)
- [ ] Todos los tests pasan.
- [ ] Cobertura del módulo auth ≥ 90%.
- [ ] Login funciona con usuario admin y empleado de seed.
- [ ] Ruta protegida rechaza requests sin token.
- [ ] El token expira correctamente.
