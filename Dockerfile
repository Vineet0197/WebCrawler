# Use an official Golang runtime as a parent image
FROM golang:1.21-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go app
RUN go build -o web-crawler cmd/main.go

# Use a minimal base image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the config file from the builder stage
COPY --from=builder /app/config/config.yaml /app/config/config.yaml

# Copy the built Go app from the builder stage
COPY --from=builder /app/web-crawler /app/web-crawler

# Make port 8108 available to the world outside this container
EXPOSE 8108

# Define environment variables
ENV PORT=8108
ENV LOG_LEVEL=INFO
ENV KAFKA_BROKER=localhost:9092
ENV KAFKA_TOPIC=web-crawler
ENV KAFKA_GROUP_ID=shoppin-web-crawler
ENV CFG_FILE=/app/config/config.yaml

# Run the web-crawler binary
CMD ["sh", "-c", "./web-crawler -f ${CFG_FILE}"]