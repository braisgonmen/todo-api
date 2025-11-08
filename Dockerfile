FROM golang:1.25.3-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git

WORKDIR /app

# Copiar go.mod y go.sum primero (mejor cache de layers)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar binario optimizado
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/server ./cmd/todo-api

# Etapa 2: Runtime (imagen mínima)
FROM alpine:3.18

# Instalar certificados SSL (para llamadas HTTPS si las necesitas)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copiar solo el binario desde el builder
COPY --from=builder /app/bin/server /app/server

# Usuario no-root por seguridad
RUN adduser -D -u 1000 appuser
USER appuser

# Puerto por defecto
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando de inicio
CMD ["/app/server"]