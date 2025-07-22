# Use the official Go image
FROM golang:1.24

# Install Air for automatic reloading
RUN go install github.com/air-verse/air@latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install swagger
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.5

# Copy the rest of the source code
COPY . .

# Expose the port your application runs on (e.g., 8080)
EXPOSE 8080

# Run Air for hot reloading
CMD ["air"]
