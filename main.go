package main

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_TYPE = "tcp"
	CONN_HOST = "localhost"
	CONN_PORT = "5667"

	ERR_BINDING             = 1
	ERR_ACCPT_INCOMING_CONN = 2
)

func main() {

	address := CONN_HOST + ":" + CONN_PORT

	listener, err := net.Listen(CONN_TYPE, address)
	if err != nil {
		fmt.Println("Could bind to "+address+": ", err.Error())
		os.Exit(ERR_BINDING)
	}

	// close listener on program exit
	defer listener.Close()

	// sping up message registry
	startRegistry()

	// listen for incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(ERR_ACCPT_INCOMING_CONN)
		}
		go handleRequest(conn)
	}
}

// handles incoming requests
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// attempt to parse the message
	if n < 1024 {
		buf = buf[:n]
	}
	message, err := ParseMessage(buf)
	if err != nil {
		// todo: determine how to send proper error response
		conn.Write([]byte("Message could not be processed."))
	}

	fmt.Printf("%v\n", message)

	conn.Close()
}

func startRegistry() {
	// todo: start a go-proc that runs the registry (needs better name)
}
