package main

import (
	"fmt"
	"net"
	"os"
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

	// Second, we tell the listener to accept incoming connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close() // connection closure is deferred to ensure it's closed after program exits

	// Third, we respond to the incoming connection, confirming the success in this case
	response := "HTTP/1.1 200 OK\r\n\r\n"
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
	}

}
