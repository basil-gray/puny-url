# Use official golang image as the base image
FROM golang:1.23-alpine

# Install gcc and SQLite development libraries
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
ENV CGO_ENABLED=1
RUN go build -o puny-url .

RUN apk add --no-cache sqlite sqlite-dev

# Expose the port the service runs on
EXPOSE 8080

# Command to run the executable
CMD ["./puny-url"]