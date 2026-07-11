# RestaurantOS — Continuación Capítulo III y Capítulo IV (v4)

> Este documento continúa `RestaurantOS_Capitulos_I_IV_APA7_v3.doc` a partir del apartado 3.6.4. Los datos de cobertura, número de pruebas, hallazgos de auditoría y estructura del sistema que aparecen aquí **son reales**, extraídos directamente del código y de la ejecución de `go test ./... -coverprofile=coverage.out` sobre el proyecto en su estado actual (post cierre de brechas de RF-06 y RF-08). Los datos de la Escala SUS (apartado 4.1.5) son **ilustrativos** — están marcados explícitamente y deben sustituirse por los valores reales una vez aplicada la encuesta a la muestra censal de 5 personas.

---

### 3.6.4. Diseño estadístico

Dado que la presente investigación es de tipo aplicada-tecnológica y de nivel descriptivo, no se contempla un diseño estadístico inferencial para la contrastación de hipótesis poblacionales. El tratamiento de los datos se fundamentó en la estadística descriptiva y en métricas estandarizadas de evaluación de software.

Los resultados del reporte de cobertura de pruebas unitarias, generado automáticamente mediante el comando `go test ./... -coverprofile=coverage.out` y consolidado con `go tool cover -func=coverage.out`, se expresaron como porcentaje de sentencias (*statements*) cubiertas, tanto a nivel global como desagregado por paquete (`auth`, `handlers`, `middleware`, `repository`, `services`, `export`). Los resultados obtenidos mediante el instrumento SUS aplicado a la muestra censal de 5 sujetos se procesan mediante el algoritmo de tabulación estándar del *System Usability Scale* (Brooke, 1996), consolidando las 10 respuestas de cada encuestado en un puntaje individual de 0 a 100, y estos 5 puntajes se promedian para obtener el *Score* global de la muestra. Toda esta información cuantitativa se presenta mediante tablas de evaluación, capturas de pantalla del sistema en ejecución y figuras arquitectónicas.

### 3.6.5. Análisis e interpretación de datos

El procesamiento de la información se ejecutó en tres frentes alineados con las fases del proyecto:

1. **Análisis cualitativo de requisitos.** Los datos recopilados sobre la gestión manual del restaurante (cuadernos de crédito, control empírico de insumos) se interpretaron para derivar los ocho requisitos funcionales (RF-01 a RF-08) y los diez requisitos no funcionales (RNF-01 a RNF-10) documentados en el SRS (`PLAN_SDD.md`).

2. **Validación binaria y porcentual del Funcionamiento (X4).** El reporte automático de cobertura y la Ficha de Análisis Documental se interpretaron bajo un enfoque de validación binaria (Cumple / No Cumple) por cada requisito funcional, y bajo un enfoque porcentual (cobertura de código) contrastado contra el umbral del 90% establecido en el RNF-02. A diferencia de un reporte simulado, el cálculo aquí presentado excluye deliberadamente los paquetes `cmd/*` (puntos de entrada `main.go` sin lógica de negocio, no sujetos a prueba unitaria por convención del lenguaje Go), de modo que el porcentaje refleje exclusivamente los módulos con lógica de negocio.

3. **Tabulación cuantitativa de Usabilidad (X5).** Los datos del instrumento SUS se tabulan aplicando la fórmula de Brooke (1996): para los ítems impares (1, 3, 5, 7, 9) se resta 1 al valor marcado; para los ítems pares (2, 4, 6, 8, 10) se resta el valor marcado a 5; la suma de las diez contribuciones se multiplica por 2.5 para obtener un puntaje de 0 a 100 por encuestado. La interpretación del *Score* promedio se categoriza según los rangos de Bangor, Kortum y Miller (2008).

### 3.6.6. Técnicas para aplicar la metodología Spec Driven Development (SDD)

La metodología rectora y única del proyecto es Spec Driven Development (SDD), descrita en el apartado 2.2.2. No se aplicó el marco Scrum ni ninguna de sus ceremonias (no hubo *Daily Scrum*, *Sprint Review* ni *Retrospective* como tales); las menciones a "Sprint" en este documento (Tabla 5, apartado 4.1.3) se emplean únicamente como sinónimo llano de "iteración de desarrollo" y no implican la adopción del marco Scrum.

En cuanto a los roles, el investigador asumió una dualidad funcional secuencial propia de SDD y no de Scrum: actuó como **Analista (A)** durante la fase de especificación, elaborando el SRS y las tarjetas de funcionalidad (*Feature Cards*) a partir de la entrevista con el propietario del restaurante; como **Diseñador (D)** durante la fase de arquitectura y modelado de la base de datos; como **Implementador (I)** durante la codificación, acompañando cada función con su respectiva prueba unitaria bajo el ciclo "rojo-verde-refactor" de TDD (Beck, 2003); y como **responsable de Calidad y Despliegue** durante la verificación de cobertura y la auditoría de concordancia final.

En cuanto a los artefactos, se gestionaron los siguientes, todos propios del flujo SDD:

- **Documento de Especificación de Requisitos de Software (SRS):** `PLAN_SDD.md`, con los requisitos funcionales (RF-01 a RF-08) y no funcionales (RNF-01 a RNF-10).
- **Tarjetas de funcionalidad (*Feature Cards* FC-01 a FC-08):** una por cada módulo del sistema, redactadas antes de escribir cualquier código de dicho módulo.
- **Pruebas unitarias escritas antes del código (TDD):** por cada Feature Card, el flujo aplicado fue (1) redactar las pruebas a partir de los criterios de aceptación, (2) verificar que fallan por ausencia de implementación, (3) implementar el código mínimo necesario, (4) refactorizar manteniendo las pruebas en verde.
- **Reporte de cobertura de pruebas:** generado automáticamente al final de cada iteración mediante el flujo de integración continua (`test.yml`).

<p align="center"><i>Tabla 3 (continuación)</i></p>
<p align="center"><i>Entregables SDD aplicados en cada iteración del proyecto</i></p>

| Tarea SDD | Artefacto generado | Verificación aplicada |
|---|---|---|
| Especificación | Feature Card del módulo correspondiente | Revisión de criterios de aceptación redactados en la propia tarjeta |
| Escritura de pruebas (rojo) | Archivo `*_test.go` con las pruebas del módulo | `go test` falla por ausencia de implementación (esperado) |
| Implementación (verde) | Código fuente del módulo (`services`, `handlers`, `repository`) | `go test` pasa en verde |
| Refactorización | Ajustes de diseño sin alterar el comportamiento probado | `go test` se mantiene en verde tras cada cambio |
| Integración continua | Ejecución automática de la suite completa en cada `push` | Flujo `test.yml` de GitHub Actions, con bloqueo si la cobertura cae bajo el umbral del RNF-02 |
| Auditoría de concordancia (iteración de cierre) | Contraste manual entre el SRS y el código realmente implementado | Detección de las brechas documentadas en la Tabla 7 (apartado 4.1.4) |

*Nota.* Elaboración propia a partir del flujo de trabajo real ejecutado durante el desarrollo de RestaurantOS, íntegramente bajo Spec Driven Development.

---

## Capítulo IV — Resultados y Discusión (continuación)

### 4.1.1. Resultados de la Fase de Análisis (Variable X1)

*(Este apartado ya está desarrollado en la Tabla 4 del documento v3 — RF-01 a RF-08 y RNF-01 a RNF-10 — y no requiere cambios; el número de requisitos funcionales y no funcionales especificados en el SRS se mantiene sin variación).*

Cabe precisar, sin embargo, un hallazgo relevante para la investigación: durante la fase de auditoría de concordancia (ver 4.1.4), se detectó que dos de los ocho requisitos funcionales especificados en el SRS —RF-06 (Reportes) y RF-08 (Página pública, en su componente de menú dinámico)— no contaban con implementación real en el código, pese a estar completamente documentados desde la fase de análisis. Este hallazgo se documenta con detalle en el apartado 4.1.4 como evidencia del valor del análisis documental como técnica de control de calidad.

### 4.1.2. Resultados de la Fase de Diseño (Variable X2)

*(La descripción arquitectónica de tres capas y la organización de directorios ya documentada en el v3 se mantiene vigente). Se añade la siguiente precisión de diseño, resultado del Sprint de cierre de brechas:*

Para el módulo de Reportes (RF-06) se incorporó un paquete adicional, `internal/export`, desacoplado de la capa `handlers`, cuya única responsabilidad es transformar las estructuras de dominio (`ReporteVentas`, `ReporteMovimientos`, `[]models.Cliente`) en archivos binarios PDF (mediante la librería `gofpdf`) y Excel (mediante `excelize`). Esta decisión de diseño —aislar la generación de archivos de la lógica HTTP— permitió alcanzar una cobertura de pruebas del 100% en dicho paquete, verificando la firma binaria de cada archivo generado (cabecera `%PDF` para PDF y `PK` para el formato ZIP subyacente de `.xlsx`) sin necesidad de abrir los archivos con un lector externo.

Asimismo, se incorporó al modelo de datos la tabla `menu_publicos` (ya prevista en el diseño original pero sin consumidores en el código), habilitando un endpoint público de solo lectura (`GET /api/public/menu`) y un conjunto de endpoints administrativos protegidos por rol (`GET/POST/PUT/DELETE /api/admin/menu`) para que el menú de la landing page sea gestionable desde el dashboard, cerrando así la brecha de diseño identificada en RF-08.

### 4.1.3. Resultados de la Fase de Implementación (Variable X3)

La Tabla 5 del documento v3 describe los Sprints 0 a 5. A partir de la auditoría de concordancia realizada, se ejecutó una iteración adicional no contemplada en la planificación original, documentada a continuación como evidencia del carácter iterativo-incremental del proyecto:

<p align="center"><i>Tabla 5 (continuación)</i></p>

| Iteración | Objetivo principal | Entregables clave | Estado |
|---|---|---|---|
| Sprint 6 — Cierre de brechas de calidad | Auditar la concordancia entre el SRS y el código implementado, y cerrar las brechas detectadas en RF-03, RF-06 y RF-08. | Módulo de Reportes (`internal/export`, `internal/handlers/reportes_handler.go`, `internal/services/reportes_service.go`) con exportación PDF/Excel; módulo de gestión del Menú público (`internal/repository/menu_repo.go`, `internal/handlers/menu_handler.go`); endpoint `PUT /pedidos/:id/entregar` para completar el ciclo de estados de RF-03; corrección del flujo de reabastecimiento del *script* de datos históricos. | Completado |

*Nota.* Elaboración propia. Esta iteración no formaba parte de la planificación inicial de Feature Cards; surgió como resultado directo del análisis documental descrito en el apartado 4.1.4, lo cual constituye en sí mismo un hallazgo metodológico: la aplicación de una auditoría de concordancia SRS-código después del "cierre" nominal del proyecto permitió detectar brechas de completitud funcional que las pruebas unitarias, por sí solas, no exponen (una prueba unitaria certifica que el código existente funciona; no certifica que *falte* código).

Adicionalmente, se generó un conjunto de datos históricos simulados (`cmd/seed/main.go`) que reproduce nueve meses de operación del restaurante (noviembre de 2025 a julio de 2026), utilizando los mismos servicios de dominio (`PedidoService`, `ClienteRepo`, `MovimientoRepo`) que emplea la API en producción, en lugar de sentencias `INSERT` directas. Esta decisión de diseño garantiza que los datos de prueba respeten las mismas reglas de negocio (descuento de stock, cálculo de totales, transiciones de estado) que gobiernan al sistema real, resultando en 905 pedidos, 224 pagos y 228 eventos de reabastecimiento distribuidos de forma realista según día de la semana y estacionalidad (incremento de pedidos en la segunda quincena de diciembre).

### 4.1.4. Resultados de la Validación de Funcionamiento (Variable X4)

La evaluación de la dimensión de Funcionamiento se realizó mediante la ejecución real del comando `go test ./... -coverprofile=coverage.out` sobre el estado final del repositorio, complementada con la Ficha de Análisis Documental de concordancia SRS-código.

<p align="center"><i>Tabla 6 (actualizada con datos reales de ejecución)</i></p>

| Paquete | Pruebas | Tipo de prueba | Cobertura |
|---|---|---|---|
| `internal/auth` | Hash/verificación de contraseñas (bcrypt), generación y validación de JWT | Unitaria | 91.3% |
| `internal/middleware` | JWT válido, token inválido, ausencia de token, control de rol (`RequireRole`) | Unitaria | 100.0% |
| `internal/services` | `CalcularTotal`, `CrearPedido`, `AnularPedido`, `MarcarEntregado`, `PredecirDemanda`, `GenerarReporteVentas`, `FiltrarClientesConDeuda` | Unitaria (con *mocks*) | 96.4% |
| `internal/handlers` | Endpoints REST de los 8 módulos (clientes, pedidos, inventario, créditos, menú, reportes, auth) | Unitaria e integración (`httptest`) | 85.6% |
| `internal/repository` | Consultas GORM contra PostgreSQL real (clientes, menú, rangos de fecha de pedidos/movimientos) | Integración (BD de prueba, transacción con *rollback*) | 40.2% |
| `internal/export` | Generación de PDF/Excel para los 3 reportes | Unitaria (verificación de firma binaria) | 100.0% |
| **Total de pruebas ejecutadas** | **198 pruebas, en 17 archivos de prueba** | — | — |
| **Cobertura ponderada (paquetes de negocio, excluyendo `cmd/*`)** | — | — | **82.2%** |

*Nota.* Reporte generado mediante `go tool cover -func=coverage.out` sobre el estado del repositorio a la fecha de cierre del Sprint 6. Elaboración propia.

El resultado de 82.2% de cobertura ponderada **no alcanza** el umbral del 90% establecido en el RNF-02. Este resultado, a diferencia de una cifra simulada, constituye un hallazgo genuino de la investigación: el paquete `internal/repository` (40.2%) concentra la brecha, dado que únicamente 3 de sus 7 conjuntos de consultas (`ClienteRepo`, `MenuRepo` y los nuevos métodos `FindEntreFechas`) cuentan con pruebas de integración contra base de datos real; los repositorios `PedidoRepo` (métodos `Create`, `FindByID`, `FindAll`, `UpdateEstado`), `ProductoRepo`, `PagoRepo` y `UserRepo` se ejercitan actualmente solo de forma indirecta, a través de *mocks* en las pruebas de `services` y `handlers`, lo cual valida la lógica de negocio que los invoca pero no certifica el SQL generado por GORM contra el motor real de PostgreSQL. Se documenta esta brecha como limitación y como recomendación de trabajo futuro (ver Conclusiones).

**Hallazgos relevantes de la auditoría de calidad.** Un aporte específico de este proyecto al curso de Pruebas y Aseguramiento de Calidad de Software fue la detección de defectos reales mediante pruebas de integración contra base de datos, que las pruebas unitarias con *mocks* no habrían podido exponer:

<p align="center"><i>Tabla 7</i></p>
<p align="center"><i>Defectos detectados durante la fase de verificación</i></p>

| Nro | Defecto detectado | Técnica que lo expuso | Acción correctiva |
|---|---|---|---|
| D-01 | El campo `Disponible` del modelo `MenuPublico` usaba la etiqueta `gorm:"default:true"`; al crear un registro con `Disponible: false` explícito, GORM omitía la columna del `INSERT` (por ser el *zero value* de un campo con `default`) y la base de datos aplicaba `true` en su lugar, invirtiendo silenciosamente la intención del usuario. | Prueba de integración `TestMenuRepo_FindPublico_SoloRetornaDisponibles` contra PostgreSQL real | Se eliminó la dependencia del *default* a nivel de columna; el valor por defecto se resolvió explícitamente en la capa de aplicación (`handler`). |
| D-02 | El requisito RF-03 especifica los estados `pendiente`, `entregado` y `anulado`, pero el código solo exponía una transición hacia `anulado`; no existía forma de marcar un pedido como `entregado` desde la API ni desde la interfaz. | Revisión manual de concordancia SRS-código (Ficha de Análisis Documental) | Se incorporó el método `PedidoService.MarcarEntregado`, el endpoint `PUT /pedidos/:id/entregar` y el control correspondiente en el dashboard. |
| D-03 | Los archivos `backend/.env` y `backend/.env.test` apuntaban a credenciales y puertos que no coincidían con los contenedores Docker reales, provocando que las migraciones y pruebas de integración se ejecutaran silenciosamente contra una instancia de PostgreSQL nativa distinta a la prevista. | Verificación de infraestructura (`docker ps`, inspección de variables de entorno del contenedor) | Se homologaron las credenciales entre `.env`, `.env.test` y `docker-compose.yml`. |
| D-04 | El *pipeline* de GitHub Actions calculaba la cobertura global incluyendo los paquetes `cmd/*` (sin pruebas por diseño), lo que distorsionaba la verificación del umbral del RNF-02 (67.0% global vs. 82.2% real de los módulos de negocio). | Análisis documental del flujo `test.yml` | Se excluyeron los paquetes `cmd/*` del cálculo de cobertura en el *pipeline*. |

*Nota.* Elaboración propia. Los defectos D-01 a D-04 se documentan como evidencia empírica de que las técnicas de análisis documental y de pruebas de integración contra base de datos real complementan —y en estos casos concretos, superan— la capacidad de detección de las pruebas unitarias basadas en *mocks*.

Pese a la brecha de cobertura en `repository`, la verificación manual de los ocho requisitos funcionales (RF-01 a RF-08) mediante Postman, tras la aplicación de las acciones correctivas de la Tabla 7, arrojó un cumplimiento del **100%** (8 de 8), frente al 75% (6 de 8) que se habría reportado de no realizarse la auditoría de concordancia. Este contraste evidencia empíricamente que la *cobertura de pruebas* y la *completitud funcional* son dimensiones distintas de calidad: un módulo puede alcanzar alta cobertura de pruebas unitarias sobre código que, sin embargo, no llega a implementar la totalidad de lo especificado en el SRS.

### 4.1.5. Resultados de la Evaluación de Usabilidad (Variable X5)

> **Nota metodológica:** los valores de esta sección son **ilustrativos**, a modo de plantilla de tabulación. El investigador debe reemplazarlos por los puntajes reales obtenidos al aplicar el instrumento SUS a la muestra censal de 5 miembros del personal de Sumaq Mikhuy.

**Instrumento aplicado (System Usability Scale, Brooke 1996):**

1. Creo que me gustaría usar este sistema con frecuencia.
2. Encontré el sistema innecesariamente complejo.
3. Pensé que el sistema era fácil de usar.
4. Creo que necesitaría el apoyo de una persona con conocimientos técnicos para usar el sistema.
5. Encontré que las distintas funciones del sistema estaban bien integradas.
6. Pensé que había demasiada inconsistencia en el sistema.
7. Imagino que la mayoría de las personas aprenderían a usar el sistema muy rápidamente.
8. Encontré el sistema muy incómodo de usar.
9. Me sentí muy seguro usando el sistema.
10. Necesité aprender muchas cosas antes de poder usar el sistema.

Cada ítem se responde en una escala de Likert de 1 (totalmente en desacuerdo) a 5 (totalmente de acuerdo).

<p align="center"><i>Tabla 8 (datos ilustrativos — reemplazar por datos reales)</i></p>
<p align="center"><i>Tabulación de ejemplo del instrumento SUS</i></p>

| Encuestado | Rol | Puntaje SUS (0–100) |
|---|---|---|
| P1 | Propietario (administrador) | 82.5 |
| P2 | Colaborador — atención | 77.5 |
| P3 | Colaborador — cocina | 65.0 |
| P4 | Colaborador — atención | 85.0 |
| P5 | Colaborador — cocina (mayor edad, menor experiencia digital) | 60.0 |
| **Promedio de la muestra** | — | **74.0** |

*Nota.* Datos ilustrativos elaborados a modo de ejemplo de tabulación. Sustituir por los resultados reales de la aplicación del instrumento.

Según las escalas de aceptabilidad de Bangor, Kortum y Miller (2008), un puntaje superior a 68 determina que un sistema es "Aceptable"; el rango 68–80 corresponde a la categoría "Buena" (*Good*) usabilidad, y valores superiores a 80 a la categoría "Excelente". El puntaje ilustrativo promedio de 74.0 se ubicaría en la categoría "Buena", consistente con la heterogeneidad esperada en la alfabetización digital del personal (ítem 4, con puntuación más moderada en los encuestados de mayor edad), tal como se anticipó en la limitación 1.5.4.5 del Capítulo I.

---

## Guía de capturas de pantalla recomendadas

Para sustentar visualmente el Capítulo IV, se recomienda incluir las siguientes capturas (organizadas por el apartado que ilustran):

| Apartado | Captura recomendada | Dónde obtenerla |
|---|---|---|
| 4.1.1 (Análisis) | Vista de las 8 Feature Cards en el explorador de archivos, o el índice de `PLAN_SDD.md` | `feature-cards/` |
| 4.1.2 (Diseño) | Diagrama de arquitectura de 3 capas (elaborar en draw.io o similar) | A elaborar |
| 4.1.2 (Diseño) | Modelo entidad-relación (8 tablas) | Exportar desde `\d` de `psql` o una herramienta ER como dbdiagram.io |
| 4.1.3 (Implementación) | Login (`/login`) y Dashboard Home (`/dashboard`) | Navegador, sesión iniciada como admin |
| 4.1.3 (Implementación) | Gestión de clientes (`/dashboard/clientes`) con el modal "Nuevo cliente" abierto | Navegador |
| 4.1.3 (Implementación) | Gestión de pedidos (`/dashboard/pedidos`) con el modal "Nuevo pedido" y los badges de estado (pendiente/entregado/anulado) visibles | Navegador |
| 4.1.3 (Implementación) | Gestión de inventario (`/dashboard/inventario`) mostrando alertas de stock bajo | Navegador |
| 4.1.3 (Implementación) | Créditos y pagos (`/dashboard/creditos`) con la lista de deudores | Navegador |
| 4.1.3 (Implementación) | Gestión del menú público (`/dashboard/menu`) | Navegador |
| 4.1.3 (Implementación) | Landing page pública (`/`) mostrando el menú dinámico ya conectado | Navegador |
| 4.1.3 (Implementación) | Reportes (`/dashboard/reportes`), las 3 pestañas (ventas, deudores, inventario) | Navegador |
| 4.1.3 (Implementación) | Un archivo PDF y un Excel descargados desde Reportes, abiertos | Explorador de archivos / Excel / lector PDF |
| 4.1.4 (Funcionamiento) | Terminal con la salida completa de `go test ./... -v` (o al menos el resumen `ok`/`PASS` por paquete) | Terminal, `backend/` |
| 4.1.4 (Funcionamiento) | Terminal con `go tool cover -func=coverage.out` mostrando el desglose por función | Terminal, `backend/` |
| 4.1.4 (Funcionamiento) | Ejecución en verde del flujo `test.yml` en la pestaña "Actions" de GitHub | github.com/&lt;usuario&gt;/&lt;repo&gt;/actions |
| 4.1.4 (Funcionamiento) | Colección de Postman con los 8 RF verificados (folders por módulo) | Postman |
| 4.1.5 (Usabilidad) | Formulario del instrumento SUS aplicado (Google Forms o físico) | A elaborar cuando se aplique la encuesta real |
| 4.1.5 (Usabilidad) | Gráfico de barras con los 5 puntajes SUS individuales y el promedio | Elaborar en Excel/Sheets a partir de la Tabla 8 real |

**Recomendación práctica:** numera cada captura como Figura N y referénciala en el texto (p. ej., "como se observa en la Figura 12"), siguiendo el formato APA7 ya usado para las tablas del documento v3.

---

## Comandos para reproducir las cifras de este capítulo

```bash
cd backend

# Cobertura completa (incluye cmd/*, útil para ver el detalle)
DATABASE_URL="postgres://postgres:santiago09@localhost:5433/restaurantos_test?sslmode=disable" \
  go test ./... -coverprofile=coverage.out -v

# Porcentaje global por paquete
go tool cover -func=coverage.out

# Cobertura excluyendo cmd/* (la que usa ahora el gate de CI)
grep -v "restaurantos/cmd/" coverage.out > coverage_negocio.out
go tool cover -func=coverage_negocio.out | grep "^total:"

# Contar pruebas ejecutadas
go test ./... -v 2>&1 | grep -c "^--- PASS"
```
