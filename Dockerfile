# Build stage
FROM golang:latest AS builder
WORKDIR /app

# Copy go.mod and go.sum, and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build the binary
COPY . .

# Disable cgo for a static build.
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 go build -o receipt-processor .

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/receipt-processor .
EXPOSE 8080
CMD ["./receipt-processor"]