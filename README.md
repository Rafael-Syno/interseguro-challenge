# Factorización QR + Estadísticas

Dos APIs y un frontend:

- `api-go` (Go + Fiber, :8080) — factoriza una matriz en Q y R (Householder) y llama a la API de estadísticas.
- `api-node` (Node + Express, :3000) — máximo, mínimo, promedio, suma y si alguna matriz es diagonal.
- `frontend` (HTML + JS, :8081) — formulario para probar el flujo completo.

Autenticación JWT HS256 con secreto compartido entre ambas APIs.

## Levantar

```bash
cp .env.example .env
docker compose up --build
```

- Frontend: http://localhost:8081
- Health: http://localhost:8080/health · http://localhost:3000/health

Sin Docker:

```bash
cd api-node && npm install && npm start
cd api-go   && go mod tidy && go run .
```

## Uso

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/token \
  -H 'Content-Type: application/json' \
  -d '{"client_id":"interseguro","client_secret":"challenge2001"}' | jq -r .access_token)

curl -s -X POST http://localhost:8080/api/v1/matrix/qr \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"matrix":[[12,-51,4],[6,167,-68],[-4,24,-41]]}' | jq
```

Devuelve `input`, `factorization` (Q, R) y `statistics`. Con `"skipStatistics": true` responde solo Q y R.

## Endpoints

| Método | Ruta | Auth |
|---|---|---|
| GET | `:8080/health` | no |
| POST | `:8080/api/v1/auth/token` | no |
| POST | `:8080/api/v1/matrix/qr` | JWT |
| GET | `:3000/health` | no |
| POST | `:3000/api/v1/statistics` | JWT |

## Pruebas

```bash
cd api-go   && go test ./...
cd api-node && npm test
```

## Producción

https://interseguro-challenge.duckdns.org — desplegado en AWS EC2 (t2.micro, us-east-2), nginx + Let's Encrypt.

## Despliegue

Ambas imágenes son stateless: una app por imagen en Azure Container Apps (o Cloud Run / App Runner), `stats-api-node` con ingress interno y `qr-api-go` con ingress externo, y `JWT_SECRET` en el gestor de secretos. El frontend, como sitio estático.
