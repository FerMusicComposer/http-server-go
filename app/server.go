package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Steps to build an http server:
	// First, a listener is created and bound to a tcp port. In this case is port 4221
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close() // listener closure is deferred to ensure it's closed after program exits

	fmt.Println("Server started on port 4221")
	fmt.Println("")

	// Second, the listener is set to accept incoming connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close() // connection closure is deferred to ensure it's closed after program exits

	// Third, the server parses the request and splits it in 2 parts, this way the second half is the path to check
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error parsing the request: ", err.Error())
		return
	}
	parts := strings.Split(requestLine, " ")

	if len(parts) < 2 {
		fmt.Println("Invalid request")
		return
	}

	path := parts[1]

	// Extract headers from request
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		headerParts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(headerParts) == 2 {
			headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		}
	}

	// Fourth, based on the path we respond to the incoming connection by checking which route is being requested
	// and responding with the appropriate handler or route
	// Else 404 NOT FOUND
	var response string
	routes := map[string]Handler{
		"/echo/":      handleEchoPath,
		"/user-agent": handleUserAgent,
	}

	handler, _ := findHandler(path, routes)
	switch {
	case handler != nil:
		contentType, content := handler(path, headers)
		response = generateResponse(contentType, content)
	case handler == nil && path == "/":
		response = "HTTP/1.1 200 OK\r\n\r\n"
	default:
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
	}
}

func generateResponse(contentType, content string) string {
	contentLength := len(content)
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", contentType, contentLength, content)
}

func findHandler(path string, routes map[string]Handler) (Handler, string) {
	for prefix, handler := range routes {
		if strings.HasPrefix(path, prefix) {
			return handler, prefix
		}
	}

	return nil, ""
}
