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

	content := string(buffer[:])
	endpoint := getEndpoint(content)
	print(content + "\n")
	headers := strings.Split(endpoint, "/")
	print(headers[0] + "\n")
	print(len(headers))

	if headers[0] == "echo" {
		// Initialisation
		content = strings.Split(strings.TrimPrefix(content, "GET /echo/"), " ")[0]
		contentLength := fmt.Sprintf("%d", len(content))

		// Generate Response
		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + contentLength + "\r\n\r\n" + content
		print(response)
		// Send Response
		_, err = connection.Write([]byte(response))
	} else if headers[0] == "user-agent" {
		// Initialisation
		userAgent := strings.Split(strings.ToLower(content), "user-agent: ")[1]
		userAgent = strings.Split(userAgent, "\r\n")[0]
		print(userAgent + "\n")
		agentLength := fmt.Sprintf("%d", len(userAgent))

		// Generate Response
		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + agentLength + "\r\n\r\n" + userAgent
		_, err = connection.Write([]byte(response))
	} else if headers[0] == "" {
		// 200 Response, Request found / valid
		response := "HTTP/1.1 200 OK\r\n\r\n"
		_, err = connection.Write([]byte(response))
	} else {
		// 404 Response, Request not found
		response := "HTTP/1.1 404 Not Found\r\n\r\n"
		_, err = connection.Write([]byte(response))
	}

}

func getEndpoint(content string) string {
	return strings.Split(strings.TrimPrefix(content, "GET /"), " ")[0]
}
