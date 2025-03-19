# GBN Protocol Implementation

This project implements the Go-Back-N (GBN) reliable transport protocol in Go. 

## Project Structure

- `cmd/main.go`: Entry point of the application.
- `internal/protocol/gbn.go`: Implementation of the GBN protocol.
- `Dockerfile`: Instructions to build the Docker image.
- `go.mod`: Module definition for the Go project.

## Docker Setup

To build and run the Docker image, use the following commands:

```bash
docker build -t gbn-protocol .
docker run gbn-protocol
```

## GBN Protocol Overview

The Go-Back-N protocol is a sliding window protocol for reliable data communication. It allows the sender to send multiple frames before needing an acknowledgment for the first one, improving the efficiency of the communication. 

For more details on the implementation, refer to the code in `internal/protocol/gbn.go`.