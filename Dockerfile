# Use the official golang image as base
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a minimal base image for the final container
FROM alpine:latest

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory inside the container
WORKDIR /home/appuser/

# Copy the compiled binary from the builder stage and change ownership
COPY --from=builder /app/app .

# Change ownership of the directory to the non-root user
RUN chown -R appuser:appgroup /home/appuser/

# Switch to the non-root user
USER appuser

# Expose port 8080
EXPOSE 8089

# Command to run the executable
CMD ["./app"]
