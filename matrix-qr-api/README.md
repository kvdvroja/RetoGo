# matrix-qr-api

API en Go (Fiber) que recibe una matriz rectangular, calcula su factorización QR y envía Q y R a `matrix-stats-api` para obtener estadísticas.

## Endpoints

| Método | Ruta              | Auth | Descripción                         |
|--------|-------------------|------|-------------------------------------|
| GET    | `/`               | No   | Nombre del servicio y rutas         |
| GET    | `/health`         | No   | Estado del servicio                 |
| POST   | `/api/auth/login` | No   | Devuelve un JWT                     |
| POST   | `/api/qr`         | JWT  | Factorización QR + estadísticas     |

### POST /api/auth/login

```json
{ "username": "admin", "password": "..." }
```

Respuesta:

```json
{ "token": "eyJhbGciOi...", "expiresIn": 3600 }
```

### POST /api/qr

Header `Authorization: Bearer <token>`.

```json
{ "matrix": [[12, -51, 4], [6, 167, -68], [-4, 24, -41]] }
```

Respuesta:

```json
{
  "q": [[0.857, -0.394, 0.331], [0.429, 0.903, -0.034], [-0.286, 0.171, 0.943]],
  "r": [[14, 21, -14], [0, 175, -70], [0, 0, -35]],
  "statistics": {
    "max": 175,
    "min": -70,
    "average": 5.22,
    "sum": 93.92,
    "count": 18,
    "anyDiagonal": false,
    "matrices": [
      { "name": "Q", "rows": 3, "cols": 3, "isDiagonal": false },
      { "name": "R", "rows": 3, "cols": 3, "isDiagonal": false }
    ]
  }
}
```

Errores: `400` entrada inválida, `401` token ausente o inválido, `502` la API de estadísticas no respondió correctamente.

## Decisiones

- La factorización usa rotaciones de Givens. Funciona con matrices de cualquier forma (m x n), es numéricamente estable y conecta directamente con la "rotación" que menciona el enunciado. Se devuelve la forma completa: Q es m x m ortogonal y R es m x n triangular superior.
- Los valores con magnitud menor a 1e-10 se redondean a 0 para que R sea exactamente triangular y no arrastre ruido de punto flotante.
- El tamaño máximo es 100 x 100 y el cuerpo de la petición está limitado a 1 MB.
- El token recibido se reenvía a la API de Node, que lo valida con el mismo secreto. Así ambas APIs quedan protegidas sin un servicio de autenticación adicional.
- El cliente HTTP hacia Node está detrás de una interfaz (`stats.Calculator`), lo que permite probar los handlers sin red.

## Variables de entorno

| Variable         | Por defecto             | Descripción                              |
|------------------|-------------------------|------------------------------------------|
| `PORT`           | `8080`                  | Puerto HTTP                              |
| `JWT_SECRET`     | requerido (min. 16)     | Secreto HS256 compartido con Node        |
| `AUTH_USERNAME`  | `admin`                 | Usuario para el login                    |
| `AUTH_PASSWORD`  | requerido               | Contraseña para el login                 |
| `TOKEN_TTL`      | `1h`                    | Vigencia del token                       |
| `STATS_API_URL`  | `http://localhost:3000` | URL base de la API de Node               |
| `STATS_TIMEOUT`  | `5s`                    | Timeout de la llamada a Node             |
| `ALLOWED_ORIGIN` | `*`                     | Origen permitido para CORS               |

## Ejecutar

Local:

```bash
go mod download
JWT_SECRET=un-secreto-de-prueba-largo AUTH_PASSWORD=admin123 go run ./cmd/server
```

Docker:

```bash
cp .env.example .env
docker compose up --build
```

Requiere `matrix-stats-api` corriendo en el puerto 3000 con el mismo `JWT_SECRET`. Desde el contenedor se accede a ella por `host.docker.internal`.

## Pruebas

```bash
go test ./...
```

- `internal/matrix`: Q·R reconstruye la matriz, Q es ortogonal, R es triangular superior, validaciones.
- `internal/handler`: login, protección JWT, validación de entrada, reenvío del token y manejo de errores de Node.
- `internal/stats`: cliente HTTP contra un servidor de prueba.

## Estructura

```
cmd/server        punto de entrada
internal/config   variables de entorno
internal/matrix   factorización QR
internal/auth     emisión y validación de JWT
internal/stats    cliente HTTP hacia la API de Node
internal/handler  rutas y manejo de errores
```
