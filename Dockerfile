FROM golang:1.26-alpine

RUN adduser -D -u 1001 appuser

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main ./cmd/core-service/main.go

RUN chown -R appuser:appuser /app

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8081/health || exit 1

USER appuser
CMD ["./main"]
