package main

import (
	"fmt"
	"net"
	"os"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	// Binds to port 4221
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// Accepts a connection
	connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// Response to send: HTTP/1.1 200 OK\r\n\r\n
	response := "HTTP/1.1 200 OK\r\n\r\n"
	_, err = connection.Write([]byte(response))
}
