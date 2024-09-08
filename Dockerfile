# Use a lightweight base image
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# airをインストール (新しいリポジトリを使用)
RUN go install github.com/air-verse/air@latest

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Expose port for the application (adjust as needed)
EXPOSE 8080

# Command to run the application
CMD ["air"]