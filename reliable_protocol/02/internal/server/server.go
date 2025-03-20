package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    // Set up the server to listen on a specific port
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error starting server:", err)
        os.Exit(1)
    }
    defer listener.Close()
    fmt.Println("Server is listening on port 8080...")

    for {
        // Accept incoming connections
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    // Handle client requests here
    fmt.Println("Client connected:", conn.RemoteAddr())
    // Additional logic for handling requests can be added here
}