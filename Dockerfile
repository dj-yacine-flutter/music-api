# Use a Go base image
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod .
COPY go.sum .

# Download the Go modules
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the binary
RUN go build -o main .

# Expose the port on which the app will be running
EXPOSE 8080

# Set the entry point to the binary we just built
CMD ["./main"]