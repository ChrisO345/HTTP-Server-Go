package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
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
	print(content)
	print(os.Args)

	params := getParams(content)

	method := params[0]
	path := "/" + params[1]

	response := ""

	switch path {
	case "/":
		response = "HTTP/1.1 200 OK\r\n\r\n"
	case "/echo":
		response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(params[2])) + "\r\n\r\n" + params[2]
	case "/user-agent":
		requestFields := strings.Split(content, "\r\n")
		for _, field := range requestFields {
			if strings.Contains(field, "User-Agent") {
				fieldValue := strings.Split(field, ": ")
				response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(fieldValue[1])) + "\r\n\r\n" + fieldValue[1] + "\r\n"
				break
			}
		}
	case "/files":
		dir := os.Args[2]
		fileName := params[2]
		print(fileName)
		if method == "POST" {
			data := content
			_ = os.WriteFile(dir+fileName, []byte(data), 0644)
			response = "HTTP/1.1 201 Created\r\n\r\n"
		} else {
			data, err := os.ReadFile(dir + fileName)
			if err != nil {
				response = "HTTP/1.1 404 Not Found\r\n\r\n"
			} else {
				response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
			}
		}
	default:
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))

	if err != nil {
		fmt.Println("Error writing: ", err.Error())
		os.Exit(1)
	}
}
