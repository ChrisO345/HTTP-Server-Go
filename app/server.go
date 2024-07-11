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
	defer l.Close()

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
	return strings.Split(strings.Split(content, " ")[1], "/")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Reads the request into a buffer
	buffer := make([]byte, 1024)
	length, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	content := string(buffer[:length])

	params := getParams(content)

	path := "/" + params[1]

	print(content + "\n")
	print(path + "\n")

	switch path {
	case "/":
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	case "/echo":
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(params[2])) + "\r\n\r\n" + params[2]))
	case "/user-agent":
		requestFields := strings.Split(content, "\r\n")
		for _, field := range requestFields {
			if strings.Contains(field, "User-Agent") {
				fieldValue := strings.Split(field, ": ")
				_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(fieldValue[1])) + "\r\n\r\n" + fieldValue[1] + "\r\n"))
				break
			}
		}
	default:
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	if err != nil {
		fmt.Println("Error writing: ", err.Error())
		os.Exit(1)
	}
}
