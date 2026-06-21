# Build stage
FROM golang:1.26-bookworm AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scheduler cmd/app/main.go

# Runtime stage
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /build/scheduler /app/scheduler
COPY web /app/web
RUN chmod +x /app/scheduler
EXPOSE 7540
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
CMD ["/app/scheduler"]