package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Handler func(path string, headers map[string]string) (contentType string, content string)

func handleEchoPath(path string, headers map[string]string) (contentType string, content string) {
	_ = headers
	body := strings.TrimPrefix(path, "/echo/")
	return "text/plain", body
}

func handleUserAgent(path string, headers map[string]string) (contentType string, content string) {
	_ = path
	return "text/plain", headers["User-Agent"]
}

func handleFileRequest(path string, headers map[string]string) (contentType string, content string) {
	_ = headers
	filename := strings.TrimPrefix(path, "/files/")
	filepath := filepath.Join(directory, filename)

	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return "", ""
	}

	content = string(fileContent)

	return "application/octet-stream", content
}
