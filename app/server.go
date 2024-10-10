package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

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
