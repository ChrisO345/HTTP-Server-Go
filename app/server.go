package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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

	// Reads the request into a buffer
	buffer := make([]byte, 1024)
	_, err = connection.Read(buffer)
	print(string(buffer[:]))

	if strings.HasPrefix(string(buffer[:]), "GET /echo/") {
		// Initialisation
		content := string(buffer[:])
		content = strings.Split(strings.TrimPrefix(content, "GET /echo/"), " ")[0]
		contentLength := string(rune(len(content)))

		// Generate Response
		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + contentLength + "\r\n\r\n" + content

		// Send Response
		_, err = connection.Write([]byte(response))
	} else if strings.HasPrefix(string(buffer[:]), "GET / HTTP/1.1") {
		// 200 Response, Request found / valid
		response := "HTTP/1.1 200 OK\r\n\r\n"
		_, err = connection.Write([]byte(response))
	} else {
		// 404 Response, Request not found
		response := "HTTP/1.1 404 Not Found\r\n\r\n"
		_, err = connection.Write([]byte(response))
	}

}
