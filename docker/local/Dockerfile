FROM golang:1.23

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ホットリロードツール
RUN go install github.com/air-verse/air@latest

# デバッガツール
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Expose port for the application (adjust as needed)
EXPOSE 8080

# Command to run the application
CMD ["air"]