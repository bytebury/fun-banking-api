# Use an official Golang runtime as the base image
FROM golang:1.21.3

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Download and install any required dependencies
RUN go get -d -v ./...

# Build the Go application
RUN go build -o main ./cmd/fun-banking

EXPOSE 8080

# Set the entry point for the application
CMD ["./app"]
