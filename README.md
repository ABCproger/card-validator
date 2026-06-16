# Card Validator

Payment card validation API built with Go, Protocol Buffers, and the Connect protocol.

Validates card number (Luhn checksum), expiration month, and expiration year. Returns `valid: true/false` with an error code on failure.

## Quick start

```bash
docker compose -f deploy/docker-compose.yml up --build
```

Service starts on `http://localhost:8080`.

## Testing the API

Import `postman_collection.json` from the repo root into Postman — all requests are ready to run.

Or use curl:

```bash
# valid card → { "valid": true }
curl -X POST http://localhost:8080/card.v1.CardService/Validate \
  -H "Content-Type: application/json" \
  -d '{"card_number":"4111111111111111","expiration_month":12,"expiration_year":2028}'

# invalid card number (Luhn fail) → { "valid": false, "error": { "code": "001", ... } }
curl -X POST http://localhost:8080/card.v1.CardService/Validate \
  -H "Content-Type: application/json" \
  -d '{"card_number":"1111111111111","expiration_month":12,"expiration_year":2028}'

# invalid month → { "valid": false, "error": { "code": "002", ... } }
curl -X POST http://localhost:8080/card.v1.CardService/Validate \
  -H "Content-Type: application/json" \
  -d '{"card_number":"4111111111111111","expiration_month":13,"expiration_year":2028}'

# expired card → { "valid": false, "error": { "code": "003", ... } }
curl -X POST http://localhost:8080/card.v1.CardService/Validate \
  -H "Content-Type: application/json" \
  -d '{"card_number":"4111111111111111","expiration_month":1,"expiration_year":2021}'

# health check → { "status": "ok" }
curl http://localhost:8080/healthz
```

## Testing via gRPC

The server exposes **gRPC server reflection**, so no `.proto` files are needed — Postman and grpcurl discover the service automatically.

### Postman

1. New Request → select **gRPC**
2. Enter URL: `localhost:8080`
3. Click **Use Server Reflection** — `card.v1.CardService/Validate` will appear
4. Select the method, enter the message body and click **Invoke**

### grpcurl

```bash
# valid card
grpcurl -plaintext \
  -d '{"card_number":"4111111111111111","expiration_month":12,"expiration_year":2028}' \
  localhost:8080 card.v1.CardService/Validate

# invalid card number
grpcurl -plaintext \
  -d '{"card_number":"1111111111111","expiration_month":12,"expiration_year":2028}' \
  localhost:8080 card.v1.CardService/Validate

# invalid month
grpcurl -plaintext \
  -d '{"card_number":"4111111111111111","expiration_month":13,"expiration_year":2028}' \
  localhost:8080 card.v1.CardService/Validate

# expired card
grpcurl -plaintext \
  -d '{"card_number":"4111111111111111","expiration_month":1,"expiration_year":2021}' \
  localhost:8080 card.v1.CardService/Validate
```

## Error codes

| Code | Reason |
|------|--------|
| `001` | Card number invalid — wrong format, length, or Luhn checksum |
| `002` | Expiration month invalid — must be 1–12 |
| `003` | Card has expired |

Validation runs in order — the first failure is returned.

## Running tests and linter

```bash
make test   # run all tests
make lint   # run golangci-lint
```

## Configuration

| Env variable | Default | Description |
|---|---|---|
| `HTTP_PORT` | `8080` | Port to listen on |
| `LOG_LEVEL` | `info` | Log level (debug/info/warn/error) |
| `APP_ENV` | `development` | Environment name |