# Go Client-Server Project

This project implements a simple client-server architecture using Go. The client connects to the server to send and receive messages.

## Project Structure

```
go-client-server
├── cmd
│   ├── client
│   │   └── main.go
│   └── server
│       └── main.go
├── internal
│   ├── client
│   │   └── client.go
│   └── server
│       └── server.go
├── pkg
│   └── shared
│       └── types.go
├── go.mod
└── README.md
```

## Setup Instructions

1. Clone the repository:
   ```
   git clone <repository-url>
   cd go-client-server
   ```

2. Initialize the Go module:
   ```
   go mod tidy
   ```

3. Build the client and server:
   ```
   go build -o client ./cmd/client
   go build -o server ./cmd/server
   ```

## Usage

1. Start the server:
   ```
   ./server
   ```

2. In a new terminal, start the client:
   ```
   ./client
   ```

## Contributing

Feel free to submit issues or pull requests for improvements or bug fixes.