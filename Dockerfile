# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the container
COPY . .
# Explicitly enable Go modules
ENV GO111MODULE=on

# Build the Go application
RUN go build -o main cmd/webhook/webhook.go

# Expose the port the application runs on
EXPOSE 5001
# Command to run the Go application
CMD ["/app/main"]