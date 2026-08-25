# Build stage
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o nerv-platform .

# Runtime stage
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/nerv-platform .

EXPOSE 8080

ENTRYPOINT ["./nerv-platform"]
