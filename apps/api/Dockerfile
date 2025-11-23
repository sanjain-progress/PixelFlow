# Build Stage
FROM golang:latest AS builder

WORKDIR /app

# Copy the entire monorepo context
COPY . .

# Remove go.work to avoid version conflicts
RUN rm -f go.work go.work.sum

# Build the API Service
WORKDIR /app/apps/api
RUN go mod download
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o api-service cmd/main.go

# Run Stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/apps/api/api-service .

EXPOSE 8080

CMD ["./api-service"]
