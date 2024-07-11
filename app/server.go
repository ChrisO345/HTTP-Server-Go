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
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println("Error closing listener: ", err.Error())
			os.Exit(1)
		}
	}(l)

	for {
		// Accepts a connection
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(connection)
	}
}

func getParams(content string) []string {
	return strings.Split(strings.TrimPrefix(content, "GET /"), "/")
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection: ", err.Error())
			os.Exit(1)
		}
	}(conn)

	for {
		// Reads the request into a buffer
		buffer := make([]byte, 1024)
		length, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			os.Exit(1)
		}

		content := string(buffer[:length])

		params := getParams(content)

		path := "/" + params[0]

		switch path {
		case "/":
			_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		case "/echo":
			_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(params[2])) + "\r\n\r\n" + params[2]))
		case "/user-agent":
			userAgent := strings.Split(strings.ToLower(content), "user-agent: ")[1]
			userAgent = strings.Split(userAgent, "\r\n")[0]
			agentLength := fmt.Sprintf("%d", len(userAgent))
			_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + agentLength + "\r\n\r\n" + userAgent))
		default:
			_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			os.Exit(1)
		}
	}
}
