# Use the official Go image as the base image
FROM golang:1.22.1

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies files
COPY backend/go.mod backend/go.sum ./

# Download all Go module dependencies
RUN go mod tidy

# Copy the backend code into the container
COPY backend/ .

# Copy the frontend directory into the container
COPY frontend/ ./frontend

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
