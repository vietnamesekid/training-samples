package main

import (
	"fmt"
	"time"
)

func main() {
	server := NewServer("localhost", WithPort(8080), WithTimeout(60*time.Second))
	fmt.Println("Server running on", server.host, ":", server.port, "with timeout", server.timeout.String())
}
