package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var directory string

func main() {
	flag.StringVar(&directory, "directory", "", "the directory from which files will be served")
	flag.Parse()

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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
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

func handleConnection(conn net.Conn) {
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

	method := parts[0]
	path := parts[1]

	// Extract headers from request
	headers := make(map[string]string)
	var contentLength int
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		headerParts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(headerParts) == 2 {
			key := strings.TrimSpace(headerParts[0])
			value := strings.TrimSpace(headerParts[1])
			headers[key] = value
			fmt.Println(headers)
			if key == "Content-Length" {
				fmt.Println("Found Content-Length header")
				contentLength, err = strconv.Atoi(value)
				if err != nil {
					fmt.Println("Error setting content lenght: ", err)
				}
			}
		}
	}

	// Read request body
	fmt.Println("contentLength: ", contentLength)
	body := make([]byte, contentLength)
	_, err = reader.Read(body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		return
	}
	fmt.Printf("Read body: %s\n", string(body))
	// Fourth, based on the path we respond to the incoming connection by checking which route is being requested
	// and responding with the appropriate handler or route
	// Else 404 NOT FOUND
	var response string
	routes := map[string]Handler{
		"/echo/":      handleEchoPath,
		"/user-agent": handleUserAgent,
		"/files/":     handleFileRequest,
	}

	fmt.Printf("Body length before passing to handler: %d\n", len(body))
	fmt.Println("body content: ", body)

	handler, _ := findHandler(path, routes)
	switch {
	case handler != nil:
		contentType, content := handler(method, path, headers, body)
		if path[:7] == "/files/" && contentType == "" && content == "" {
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		} else if method == "POST" && contentType != "" {
			response = content
		} else {
			response = generateResponse(contentType, content)
		}
	case path == "/":
		response = "HTTP/1.1 200 OK\r\n\r\n"
	default:
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
	}
}
