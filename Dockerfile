FROM golang:1.23.4-alpine

WORKDIR /app

# Install necessary build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]