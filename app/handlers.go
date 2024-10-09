package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Handler func(method, path string, headers map[string]string, body []byte) (contentType string, content string)

func handleEchoPath(method, path string, headers map[string]string, body []byte) (contentType string, content string) {
	_, _, _ = method, headers, body
	bodyContent := strings.TrimPrefix(path, "/echo/")
	return "text/plain", bodyContent
}

func handleUserAgent(method, path string, headers map[string]string, body []byte) (contentType string, content string) {
	_, _, _ = method, path, body
	return "text/plain", headers["User-Agent"]
}

func handleFileRequest(method, path string, headers map[string]string, body []byte) (contentType string, content string) {
	fmt.Printf("Received body length in handler: %d\n", len(body))
	_ = headers
	filename := strings.TrimPrefix(path, "/files/")
	filepath := filepath.Join(directory, filename)

	switch {
	case method == "GET":
		fileContent, err := os.ReadFile(filepath)
		if err != nil {
			return "", ""
		}

		content = string(fileContent)

		return "application/octet-stream", content
	case method == "POST":
		fmt.Printf("Writing to filepath: %s\n", filepath)
		err := os.WriteFile(filepath, body, 0644)
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			return "", ""
		}
		fmt.Printf("Wrote %d bytes to file\n", len(body))
		return "text/plain", "HTTP/1.1 201 Created\r\n\r\n"
	default:
		return "", ""
	}
}
