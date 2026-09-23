# --- build stage ---
FROM golang:1.27.1-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

# --- runtime stage ---
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
WORKDIR /app

COPY --from=builder /out/api ./api
COPY --from=builder /src/migrations ./migrations

USER appuser
EXPOSE 8080
ENTRYPOINT ["./api"]
