# Build Stage
FROM golang:latest AS builder

WORKDIR /app

# Copy the entire monorepo context
COPY . .

# Remove go.work to avoid version conflicts
RUN rm -f go.work go.work.sum

# Build the Auth Service
WORKDIR /app/apps/auth
RUN go mod download
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o auth-service cmd/main.go

# Run Stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/apps/auth/auth-service .

EXPOSE 50051

CMD ["./auth-service"]
