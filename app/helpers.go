package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

var directory string

func compressGzip(content string) ([]byte, error) {
	var buf bytes.Buffer

	gw := gzip.NewWriter(&buf)
	_, err := gw.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	err = gw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func acceptsContentEncoding(headers map[string]string) bool {
	acceptEncoding, ok := headers["Accept-Encoding"]
	if !ok {
		return false
	}

	encodings := strings.Split(acceptEncoding, ",")
	for _, encoding := range encodings {
		if strings.TrimSpace(strings.ToLower(encoding)) == "gzip" {
			return true
		}
	}

	return false
}

func generateResponse(contentType, content string, contentEncoding bool) string {
	var response string
	if contentEncoding {
		compressedContent, err := compressGzip(content)
		if err != nil {
			fmt.Printf("Error compressing content: %v\n", err)
			return ""
		}
		contentLength := len(compressedContent)
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n%s", contentType, contentLength, compressedContent)
	} else {
		contentLength := len(content)
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", contentType, contentLength, content)
	}
	return response
}

func handleConnection(conn net.Conn) {
	defer conn.Close() // connection closure is deferred to ensure it's closed after program exits

	// Third, the server parses the request and splits it in 2 parts, this way the second half is the path to check
	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error parsing the request: ", err.Error())
		}
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
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading header: ", err.Error())
			}
			break
		}
		if line == "\r\n" {
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
					fmt.Println("Error converting Content-Length:", err.Error())
				}
			}
		}
	}

	useContentEncoding := acceptsContentEncoding(headers)

	// Read request body
	fmt.Println("contentLength: ", contentLength)
	body := make([]byte, contentLength)
	if contentLength > 0 {
		_, err = reader.Read(body)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading request body:", err.Error())
			return
		}
		fmt.Printf("Read body: %s\n", string(body))
	}

	// Fourth, based on the path we respond to the incoming connection by checking which route is being requested
	// and responding with the appropriate handler or route
	// Else 404 NOT FOUND
	var response string
	handler, _ := router.FindHandler(path)
	switch {
	case handler != nil:
		contentType, content := handler(method, path, headers, body)
		if path[:7] == "/files/" && contentType == "" && content == "" {
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		} else if method == "POST" && contentType != "" {
			response = content
		} else {
			response = generateResponse(contentType, content, useContentEncoding)
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
