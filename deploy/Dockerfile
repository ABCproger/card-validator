FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/app \
    ./cmd/app

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/app /app
COPY --from=builder /build/config/config.yml /config/config.yml

USER nonroot

EXPOSE 8080

ENTRYPOINT ["/app"]
