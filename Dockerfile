# Build stage
FROM golang:1.23.3-alpine3.19 AS builder

# Set build arguments for user and group IDs
ARG UID=65534
ARG GID=65534

# Set build argument for binary name
ARG BINARY_NAME=recallKingApi

# Set the working directory
WORKDIR /app

# Copy only the necessary files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the binary with the specified name
RUN CGO_ENABLED=0 GOOS=linux go build -o ${BINARY_NAME} ./cmd/main

# Run stage
FROM alpine:3.19

# Set build arguments for user and group IDs
ARG UID=65534
ARG GID=65534
ARG BINARY_NAME=recallKingApi # Must match the build stage

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/${BINARY_NAME} .

# Set user and group
RUN adduser -u ${UID} -G ${GID} -D appuser
USER appuser

# Copy environment variables (use with caution!)
COPY .env .

# Expose the ports the app needs
EXPOSE 8080

# Run the app
CMD [ "./${BINARY_NAME}" ]