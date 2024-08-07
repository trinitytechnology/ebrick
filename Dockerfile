# Start from the official Golang image to build the binary.
FROM golang:1.22.5-alpine3.20 AS builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum and download dependencies (for caching).
COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN go build -o app ./examples

# Start a new stage from scratch for a smaller final image.
FROM alpine:3.20

WORKDIR /app

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/app .

RUN ls -l

CMD ["/app/app"]
