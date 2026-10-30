# matrix-frontend

Interfaz web para `matrix-qr-api`. HTML, CSS y JavaScript sin frameworks, servido con nginx.

nginx sirve los archivos estáticos y redirige `/api/` hacia la API de Go, así el navegador no necesita CORS.

## Ejecutar

Requiere `matrix-qr-api` corriendo en el puerto 8080.

```bash
docker compose up --build
```

Abrir http://localhost:8000.

Para apuntar a otra URL de la API:

```bash
GO_API_URL=https://mi-api.example.com docker compose up --build
```

## Uso

1. Ingresar con el usuario y la contraseña configurados en `matrix-qr-api`.
2. Escribir la matriz: una fila por línea, valores separados por espacios o comas. También se acepta JSON, por ejemplo `[[1,2],[3,4]]`.
3. Calcular. Se muestran Q, R y las estadísticas devueltas por `matrix-stats-api`.
