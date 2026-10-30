# matrix-stats-api

API en Node.js (Express) que recibe una lista de matrices y devuelve estadísticas sobre todos sus valores.

## Endpoints

| Método | Ruta              | Auth | Descripción              |
|--------|-------------------|------|--------------------------|
| GET    | `/`               | No   | Nombre del servicio      |
| GET    | `/health`         | No   | Estado del servicio      |
| POST   | `/api/statistics` | JWT  | Estadísticas de matrices |

### POST /api/statistics

Header `Authorization: Bearer <token>`. El token lo emite `matrix-qr-api` y se firma con el mismo `JWT_SECRET`.

```json
{
  "matrices": [
    { "name": "Q", "data": [[0.6, -0.8], [0.8, 0.6]] },
    { "name": "R", "data": [[5, 4], [0, 3]] }
  ]
}
```

Respuesta:

```json
{
  "max": 5,
  "min": -0.8,
  "average": 1.65,
  "sum": 13.2,
  "count": 8,
  "anyDiagonal": false,
  "matrices": [
    { "name": "Q", "rows": 2, "cols": 2, "isDiagonal": false },
    { "name": "R", "rows": 2, "cols": 2, "isDiagonal": false }
  ]
}
```

Errores: `400` entrada inválida, `401` token ausente, inválido o expirado.

## Decisiones

- Máximo, mínimo, suma y conteo se calculan en un solo recorrido; el promedio sale de la suma y el conteo.
- Una matriz es diagonal si todos los elementos fuera de la diagonal principal son 0. Aplica también a matrices rectangulares. Se usa una tolerancia de 1e-10 porque los valores vienen de operaciones de punto flotante.
- Además de `anyDiagonal` se devuelve el detalle por matriz, para saber cuál cumple la condición.
- El payload acepta una lista de matrices con nombre en lugar de campos fijos `q` y `r`, así el servicio no depende de la factorización que se use antes.
- Límites: hasta 10 matrices de 100 x 100 y un cuerpo de 1 MB.

## Variables de entorno

| Variable         | Por defecto         | Descripción                    |
|------------------|---------------------|--------------------------------|
| `PORT`           | `3000`              | Puerto HTTP                    |
| `JWT_SECRET`     | requerido (min. 16) | Secreto HS256 compartido       |
| `ALLOWED_ORIGIN` | `*`                 | Origen permitido para CORS     |

## Ejecutar

Local:

```bash
npm install
JWT_SECRET=un-secreto-de-prueba-largo npm start
```

Docker:

```bash
cp .env.example .env
docker compose up --build
```

`JWT_SECRET` debe ser el mismo que usa `matrix-qr-api`.

## Pruebas

```bash
npm test
```

- `tests/statistics.test.js`: cálculo de estadísticas y detección de matrices diagonales.
- `tests/app.test.js`: endpoint completo con supertest (JWT válido, falsificado y expirado, validaciones, JSON mal formado).

## Estructura

```
src/server.js                 punto de entrada
src/app.js                    configuración de Express
src/middleware/auth.js        validación de JWT
src/routes/statistics.js      ruta /api/statistics
src/services/statistics.js    cálculo de estadísticas
src/validation.js             validación del cuerpo
```
