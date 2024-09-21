package main

import "strings"

type Handler func(path string, headers map[string]string) (contentType string, content string)

func handleRoot(path string, headers map[string]string) (contentType string, content string) {
	_, _ = path, headers
	return "", ""
}

func handleEchoPath(path string, headers map[string]string) (contentType string, content string) {
	_ = headers
	body := strings.TrimPrefix(path, "/echo/")
	return "text/plain", body
}

func handleUserAgent(path string, headers map[string]string) (contentType string, content string) {
	_ = path
	return "text/plain", headers["User-Agent"]
}
