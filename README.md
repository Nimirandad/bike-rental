# Bike Rental Service API

![Coverage](https://img.shields.io/badge/coverage-94.5%25-brightgreen)

API RESTful para gestión de renta de bicicletas construida con Go, SQLite y arquitectura hexagonal.

## Tabla de Contenidos

- [Instalación](#instalación)
- [Ejecución](#ejecución)
- [Testing](#testing)
- [Documentación Swagger](#documentación-swagger)
- [Comandos Make Disponibles](#comandos-make-disponibles)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Configuración](#configuración)
- [Descripción](#descripción)
- [Características](#características)
- [Arquitectura](#arquitectura)
- [Base de Datos](#base-de-datos)
- [API Endpoints](#api-endpoints)
- [Reglas de Negocio](#reglas-de-negocio)

---

## Instalación

### Prerequisitos

- **Go 1.24+**
- **Docker** (opcional, para containerización)
- **Make** (para comandos automatizados)

### Clonar el Repositorio

```bash
git clone <repository-url>
cd bike-rental
```

### Configurar Variables de Entorno

Crea un archivo `.env` en la raíz:

```bash
# .env
PORT=8080
SQLITE_PATH=data/bike_rental.db
JWT_SECRET=change-this-secret-in-production
ADMIN_CREDENTIALS=YWRtaW46YmlrZXJlbnRhbGFkbWlu  # admin:bikerentaladmin
LOG_LEVEL=info
```

**Generar credenciales admin personalizadas**:
```bash
echo -n "admin:tu_password" | base64
```

### Instalar Dependencias

```bash
go mod download
```

---


---

## Ejecución

### Opción 1: Desarrollo Local

```bash
# 1. Ejecutar migraciones (crea DB + seed con 150 bikes)
make migrate

# 2. Generar documentación Swagger (opcional)
make swagger

# 3. Ejecutar servidor
make run

# Servidor disponible en http://localhost:8080
# Swagger UI en http://localhost:8080/swagger/index.html
```

**Comandos útiles**:
```bash
make migrate-fresh    # Resetea DB completamente
make migrate-no-seed  # Migración sin datos de ejemplo
make test             # Ejecuta tests
make test-coverage    # Tests + reporte HTML cobertura
```

---

### Opción 2: Docker

```bash
# 1. Construir imagen distroless (~10MB)
make docker-build

# 2. Ejecutar contenedor (carga .env automáticamente)
make docker-run

# 3. Ver logs
make docker-logs

# 4. Detener
make docker-stop
```

**Comandos Docker**:
```bash
make docker-build     # Construye imagen distroless
make docker-run       # Ejecuta contenedor con .env
make docker-stop      # Detiene y elimina contenedor
make docker-clean     # Limpia imágenes
make docker-logs      # Muestra logs en tiempo real
```

**Docker manual**:
```bash
# Build
docker build -f Dockerfile.distroless -t bike-rental:latest .

# Run
docker run -d \
  --name bike-rental-api \
  --env-file .env \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  bike-rental:latest

# Ver logs
docker logs -f bike-rental-api
```

---


---

## Testing

### Ejecutar Tests

```bash
# Tests completos
make test

# Tests con cobertura + reporte HTML
make test-coverage
# Abre coverage.html en el navegador
```
---

## Documentación Swagger

### Acceso a Swagger UI

Una vez ejecutado el servidor:

**URL**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Características

- 18 endpoints documentados
- Schemas de request/response
- Ejemplos de uso
- Try it out (testing interactivo)
- Autenticación JWT y Basic Auth

### Cómo Usar Swagger UI

1. **Autenticarse como Usuario**:
   - Click en botón "Authorize"
   - Selecciona "BearerAuth"
   - Ingresa: `Bearer <tu-jwt-token>`
   - Obtén el token desde: `POST /api/v1/users/login`

2. **Autenticarse como Admin**:
   - Click en "Authorize"
   - Selecciona "BasicAuth"
   - Username: `admin`
   - Password: `bikerentaladmin`

3. **Probar Endpoints**:
   - Expande cualquier endpoint
   - Click en "Try it out"
   - Modifica parámetros/body
   - Click "Execute"
   - Ver response en tiempo real

### Regenerar Documentación

```bash
make swagger
```

Esto genera/actualiza:
- `docs/docs.go`
- `docs/swagger.json`
- `docs/swagger.yaml`

---

## Comandos Make Disponibles

```bash
make help              # Muestra todos los comandos

# Desarrollo Local
make migrate           # DB migrations + seed
make migrate-fresh     # Resetea DB completamente
make migrate-no-seed   # Migración sin datos
make run               # Ejecuta servidor local
make build-linux       # Compila binario Linux (Docker)
make test              # Ejecuta tests
make test-coverage     # Tests + cobertura HTML
make swagger           # Genera Swagger docs
make clean             # Limpia artifacts

# Docker
make docker-build      # Build imagen distroless (~10MB)
make docker-run        # Ejecuta contenedor
make docker-stop       # Detiene contenedor
make docker-clean      # Elimina imágenes
make docker-logs       # Muestra logs
```
---

## Estructura del Proyecto

```
bike-rental/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point
├── internal/
│   ├── app/
│   │   └── app.go                  # Inicialización aplicación
│   ├── config/
│   │   ├── config.go               # Carga configuración
│   │   └── defaults.go             # Valores por defecto
│   ├── constants/
│   │   └── constants.go            # Constantes del sistema
│   ├── database/
│   │   ├── schema.sql              # DDL tablas
│   │   ├── seed.sql                # Datos iniciales
│   │   └── sqlite.go               # Conexión DB
│   ├── handlers/                   # Capa HTTP
│   │   ├── admin_handler.go
│   │   ├── bikes_handler.go
│   │   ├── health_handler.go
│   │   ├── rentals_handler.go
│   │   └── users_handler.go
│   ├── logger/
│   │   └── logger.go               # Configuración zerolog
│   ├── models/                     # Entidades de dominio
│   │   ├── bikes.go
│   │   ├── rentals.go
│   │   └── users.go
│   ├── repositories/               # Capa de datos
│   │   ├── admin_repository.go
│   │   ├── bike_repository.go
│   │   ├── rental_repository.go
│   │   └── user_repository.go
│   ├── routes/
│   │   └── routes.go               # Definición de rutas
│   ├── server/
│   │   ├── server.go
│   │   └── middlewares/            # Auth, logging, CORS
│   │       └── logging.go
│   ├── services/                   # Lógica de negocio
│   │   ├── admin_service.go
│   │   ├── bike_service.go
│   │   ├── health_service.go
│   │   ├── rental_service.go
│   │   └── user_service.go
│   ├── types/                      # DTOs
│   │   ├── admin.go
│   │   ├── rentals.go
│   │   ├── responses.go
│   │   └── users.go
│   └── utils/                      # Utilidades
│       ├── jwt.go
│       └── password.go
├── docs/                           # Swagger autogenerado
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── scripts/
│   └── migrate.sh                  # Script de migraciones
├── data/
│   └── bike_rental.db              # SQLite database
├── bin/                            # Binarios compilados
├── build.sh                        # Script de build
├── Dockerfile.distroless           # Imagen minimalista
├── Makefile                        # Comandos automatizados
├── .env                            # Variables de entorno
├── go.mod
├── go.sum
└── README.md                       # Este archivo
```

---

## Configuración

### Variables de Entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `PORT` | `8080` | Puerto del servidor HTTP |
| `SQLITE_PATH` | `data/bike_rental.db` | Ruta de la base de datos |
| `JWT_SECRET` | - | Secret para firmar JWT |
| `ADMIN_CREDENTIALS` | - | Base64 de `user:password` para admin |
| `LOG_LEVEL` | `info` | debug, info, warn, error |



---

## Descripción

Sistema de renta de bicicletas que permite a los usuarios:
- Registrarse y autenticarse
- Ver bicicletas disponibles cercanas
- Iniciar y finalizar rentas
- Consultar historial de rentas

Los administradores pueden:
- Gestionar flota de bicicletas (crear, actualizar, ver todas)
- Administrar usuarios
- Supervisar todas las rentas


## Características

- **Autenticación JWT** para usuarios
- **Basic Auth** para administradores
- **Geolocalización** de bicicletas (latitud/longitud)
- **Cálculo automático** de costos por minuto
- **Paginación** en listados
- **Logging estructurado** con zerolog
- **Docker distroless** (~10MB)
- **Documentación Swagger/OpenAPI**
- **Testing** con 89.2% cobertura
- **SQLite** embebido (sin dependencias externas)


## Arquitectura

Arquitectura hexagonal (Ports & Adapters):

```
bike-rental/
├── cmd/api/                    # Entry point
├── internal/
│   ├── app/                    # Aplicación principal
│   ├── config/                 # Configuración
│   ├── constants/              # Constantes del sistema
│   ├── database/               # Schemas y conexión
│   ├── handlers/               # HTTP handlers (adapters)
│   ├── logger/                 # Logging
│   ├── models/                 # Entidades de dominio
│   ├── repositories/           # Persistencia (ports)
│   ├── routes/                 # Rutas HTTP
│   ├── server/                 # Servidor HTTP
│   │   └── middlewares/        # Auth, logging, CORS
│   ├── services/               # Lógica de negocio (core)
│   ├── types/                  # DTOs y requests/responses
│   └── utils/                  # Utilidades
└── docs/                       # Swagger docs autogenerados
```

### Stack Tecnológico

- **Lenguaje**: Go 1.23
- **HTTP Router**: Chi
- **Base de Datos**: SQLite 3
- **Autenticación**: JWT (golang-jwt/jwt)
- **Logging**: zerolog
- **Documentación**: Swagger/OpenAPI (swaggo)
- **Testing**: testify, go-sqlmock

---


## Base de Datos

SQLite database con 3 tablas principales.

### Tabla: `users`

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | INTEGER | Primary key (autoincremental) |
| `email` | TEXT | Email único del usuario |
| `hashed_password` | TEXT | Contraseña hasheada (bcrypt) |
| `first_name` | TEXT | Nombre |
| `last_name` | TEXT | Apellido |
| `created_at` | DATETIME | Fecha de creación |
| `updated_at` | DATETIME | Última actualización |

**Índices**: `idx_users_email` (email)

### Tabla: `bikes`

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | INTEGER | Primary key (autoincremental) |
| `is_available` | INTEGER | Disponibilidad (0/1) |
| `latitude` | REAL | Latitud GPS |
| `longitude` | REAL | Longitud GPS |
| `price_per_minute` | REAL | Precio por minuto (€) |
| `created_at` | DATETIME | Fecha de creación |
| `updated_at` | DATETIME | Última actualización |

**Índices**: `idx_bikes_available` (is_available)

**Datos seed**: 150 bicicletas en Londres, Manchester, Birmingham, Leeds y Glasgow.

### Tabla: `rentals`

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | INTEGER | Primary key (autoincremental) |
| `user_id` | INTEGER | FK a users |
| `bike_id` | INTEGER | FK a bikes |
| `status` | TEXT | Estado: "running", "ended" |
| `start_time` | DATETIME | Inicio de la renta |
| `end_time` | DATETIME | Fin de la renta (nullable) |
| `start_latitude` | REAL | Ubicación inicial |
| `start_longitude` | REAL | Ubicación inicial |
| `end_latitude` | REAL | Ubicación final (nullable) |
| `end_longitude` | REAL | Ubicación final (nullable) |
| `duration_minutes` | INTEGER | Duración en minutos |
| `cost` | REAL | Costo total (€) |
| `created_at` | DATETIME | Fecha de creación |
| `updated_at` | DATETIME | Última actualización |

**Índices**: 
- `idx_rentals_user` (user_id)
- `idx_rentals_bike` (bike_id)
- `idx_rentals_status` (status)                                     


---

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

### Usuarios (Públicos)

#### POST `/users/register`
Registra un nuevo usuario.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response** (201):
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

**Errores**:
- `400`: Email o contraseña inválidos
- `409`: Email ya registrado

---

#### POST `/users/login`
Autenticación de usuario.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response** (200):
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe"
    }
  }
}
```

**Errores**:
- `400`: Credenciales faltantes
- `401`: Credenciales inválidas

---

#### GET `/users/profile`
Obtiene el perfil del usuario autenticado.

**Headers**: `Authorization: Bearer <token>`

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2026-02-15T10:30:00Z"
  }
}
```

**Errores**:
- `401`: Token inválido o ausente
- `404`: Usuario no encontrado

---

#### PATCH `/users/profile`
Actualiza el perfil del usuario.

**Headers**: `Authorization: Bearer <token>`

**Request Body**:
```json
{
  "first_name": "John Updated",
  "last_name": "Doe Updated",
  "password": "newPassword123"
}
```

**Response** (200):
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John Updated",
    "last_name": "Doe Updated"
  }
}
```

**Errores**:
- `401`: No autenticado
- `400`: Datos inválidos

---

### Bicicletas

#### GET `/bikes/available`
Lista bicicletas disponibles (paginado).

**Headers**: `Authorization: Bearer <token>`

**Query Parameters**:
- `page` (default: 1): Número de página
- `page_size` (default: 20): Elementos por página

**Response** (200):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "is_available": true,
        "latitude": 51.5074,
        "longitude": -0.1278,
        "price_per_minute": 0.65
      }
    ],
    "page": 1,
    "page_size": 20,
    "total_items": 150,
    "total_pages": 8
  }
}
```

**Errores**:
- `401`: No autenticado
- `400`: Parámetros de paginación inválidos

---

### Rentas

#### POST `/rentals/start`
Inicia una renta de bicicleta.

**Headers**: `Authorization: Bearer <token>`

**Request Body**:
```json
{
  "bike_id": 1,
  "start_latitude": 51.5074,
  "start_longitude": -0.1278
}
```

**Response** (201):
```json
{
  "success": true,
  "message": "Rental started successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "bike_id": 1,
    "status": "running",
    "start_time": "2026-02-15T10:30:00Z",
    "start_latitude": 51.5074,
    "start_longitude": -0.1278
  }
}
```

**Errores**:
- `401`: No autenticado
- `400`: Bicicleta no disponible
- `400`: Usuario ya tiene una renta activa
- `404`: Bicicleta no encontrada

---

#### POST `/rentals/end`
Finaliza una renta activa.

**Headers**: `Authorization: Bearer <token>`

**Request Body**:
```json
{
  "end_latitude": 51.5155,
  "end_longitude": -0.0922
}
```

**Response** (200):
```json
{
  "success": true,
  "message": "Rental ended successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "bike_id": 1,
    "status": "ended",
    "start_time": "2026-02-15T10:30:00Z",
    "end_time": "2026-02-15T11:00:00Z",
    "duration_minutes": 30,
    "cost": 19.50,
    "start_latitude": 51.5074,
    "start_longitude": -0.1278,
    "end_latitude": 51.5155,
    "end_longitude": -0.0922
  }
}
```

**Cálculo de costo**: `duration_minutes * price_per_minute`

**Errores**:
- `401`: No autenticado
- `400`: No hay renta activa
- `404`: Renta no encontrada

---

#### GET `/rentals/history`
Obtiene el historial de rentas del usuario.

**Headers**: `Authorization: Bearer <token>`

**Query Parameters**:
- `page` (default: 1)
- `page_size` (default: 10)

**Response** (200):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "bike_id": 1,
        "status": "ended",
        "start_time": "2026-02-15T10:30:00Z",
        "end_time": "2026-02-15T11:00:00Z",
        "duration_minutes": 30,
        "cost": 19.50
      }
    ],
    "page": 1,
    "page_size": 10,
    "total_items": 5,
    "total_pages": 1
  }
}
```

---

### Admin (Requiere Basic Auth)

**Credenciales por defecto**:
- Usuario: `admin`
- Password: `bikerentaladmin`
- Header: `Authorization: Basic YWRtaW46YmlrZXJlbnRhbGFkbWlu`

#### POST `/admin/bikes`
Crea una nueva bicicleta.

**Headers**: `Authorization: Basic <credentials>`

**Request Body**:
```json
{
  "latitude": 51.5074,
  "longitude": -0.1278,
  "price_per_minute": 0.65
}
```

**Response** (201):
```json
{
  "success": true,
  "message": "Bike added successfully",
  "data": {
    "id": 151,
    "is_available": true,
    "latitude": 51.5074,
    "longitude": -0.1278,
    "price_per_minute": 0.65
  }
}
```

---

#### PATCH `/admin/bikes/{bike-id}`
Actualiza una bicicleta.

**Headers**: `Authorization: Basic <credentials>`

**Request Body** (todos opcionales):
```json
{
  "is_available": false,
  "latitude": 51.5080,
  "longitude": -0.1280,
  "price_per_minute": 0.70
}
```

**Response** (200):
```json
{
  "success": true,
  "message": "Bike updated successfully",
  "data": { ... }
}
```

---

#### GET `/admin/bikes`
Lista todas las bicicletas (paginado).

**Headers**: `Authorization: Basic <credentials>`

**Query Parameters**: `page`, `page_size`

---

#### GET `/admin/users`
Lista todos los usuarios (paginado).

**Headers**: `Authorization: Basic <credentials>`

**Response** (200):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "email": "user@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "created_at": "2026-02-15T10:30:00Z"
      }
    ],
    "page": 1,
    "page_size": 20,
    "total_items": 100,
    "total_pages": 5
  }
}
```

---

#### GET `/admin/users/{user-id}`
Obtiene detalles de un usuario específico.

---

#### PATCH `/admin/users/{user-id}`
Actualiza datos de un usuario.

**Request Body**:
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "password": "newPassword"
}
```

---

#### GET `/admin/rentals`
Lista todos las rentas del sistema.

**Query Parameters**: `page`, `page_size`

**Response** (200):
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "user_id": 1,
        "bike_id": 1,
        "status": "ended",
        "start_time": "2026-02-15T10:30:00Z",
        "end_time": "2026-02-15T11:00:00Z",
        "duration_minutes": 30,
        "cost": 19.50
      }
    ],
    "page": 1,
    "page_size": 20,
    "total_items": 500,
    "total_pages": 25
  }
}
```

---

#### GET `/admin/rentals/{rental-id}`
Obtiene detalles de una renta específica.

---

#### PATCH `/admin/rentals/{rental-id}`
Actualiza una renta (ej. cambiar status a "ended" o "running").

**Request Body**:
```json
{
  "status": "ended"
}
```

---

### Health Check

#### GET `/status`
Verifica el estado de la API y la base de datos.

**No requiere autenticación**

**Response** (200):
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2026-02-15T12:00:00Z"
}
```

**Response** (503) - Si hay problemas:
```json
{
  "status": "unhealthy",
  "database": "disconnected",
  "timestamp": "2026-02-15T12:00:00Z"
}
```

---


---

## Reglas de Negocio

### Usuarios

1. **Registro**:
   - Email debe ser único
   - Email debe tener formato válido (`@` presente)
   - Contraseña mínimo 6 caracteres
   - Nombre y apellido son obligatorios

2. **Autenticación**:
   - JWT válido por 24 horas
   - Password hasheado con bcrypt (cost 10)

### Bicicletas

1. **Disponibilidad**:
   - Una bicicleta solo puede tener una renta activa
   - Bicicleta se marca como no disponible al iniciar renta
   - Bicicleta vuelve a estar disponible al finalizar renta

2. **Precios**:
   - Precio por minuto configurable por bicicleta
   - Rango típico: €0.35 - €0.70 por minuto

### Rentas

1. **Inicio**:
   - Usuario debe estar autenticado
   - Usuario solo puede tener 1 renta activa a la vez
   - Bicicleta debe estar disponible
   - Se requiere ubicación inicial

2. **Finalización**:
   - Se requiere ubicación final
   - Cálculo automático de:
     - Duración: `end_time - start_time` (redondeado a minutos)
     - Costo: `duration_minutes * bike.price_per_minute`
   - Status cambia a "ended"
   - Bicicleta vuelve a estar disponible

3. **Estados posibles**:
   - `running`: Renta en curso
   - `ended`: Finalizado normalmente

### Paginación

- Default `page_size`: 20 elementos
- Máximo `page_size`: 100 elementos
- `page` inicia en 1
- Response incluye: `total_items`, `total_pages`, `page`, `page_size`

---

