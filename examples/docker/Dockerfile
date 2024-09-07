# Start from the official Golang image to build the binary.
FROM golang:1.22.3-alpine3.20 AS builder

# Set the Current Working Directory inside the container.
WORKDIR /app
COPY . .
RUN go build -o app

# Start a new stage from scratch for a smaller final image.
FROM alpine:3.20

WORKDIR /app

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/app .
CMD ["./app"]
