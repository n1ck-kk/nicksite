# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/server .
COPY --chown=appuser:appuser static ./static

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
