# --- Backend build stage ---
FROM golang:1.24-alpine AS api-build
WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /out/api ./cmd/api


# --- Final runtime image ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=api-build /out/api ./api

RUN chown -R appuser:appuser /app
USER appuser

ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

CMD ["./api"]
