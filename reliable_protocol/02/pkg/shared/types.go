package shared

// Message represents a message structure for communication between client and server.
type Message struct {
    ID      string `json:"id"`
    Content string `json:"content"`
}

// Status represents the status of a response from the server.
type Status struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}