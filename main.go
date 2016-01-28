package main

import (
	"net"
	"os"
)

const (
	connType = "tcp"
	connHost = "localhost"
	connPort = "5667"

	errBinding           = 1
	errAccptIncomingConn = 2

	gatewayMessageBufferSize        = 100
	messageCacheTTLSeconds          = 20
	messageExpirationChanBufferSize = 100
)

func main() {
	Logger().Info.Println("Starting up NBAd")

	address := connHost + ":" + connPort

	listener, err := net.Listen(connType, address)
	if err != nil {
		Logger().Error.Println("could not bind to "+address, err.Error())
		os.Exit(errBinding)
	}

	// close listener on program exit
	defer listener.Close()

	// sping up message registry
	messageChannel := startGateway()

	// listen for incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			Logger().Error.Println("Error accepting connection", err.Error())
			os.Exit(errAccptIncomingConn)
		}
		go handleRequest(conn, messageChannel)
	}
}

// handles incoming requests
func handleRequest(conn net.Conn, messageChannel chan *Message) {
	defer conn.Close()

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	n, err := conn.Read(buf)
	if err != nil {
		Logger().Warning.Println("Error reading incoming request", err.Error())
		return
	}

	// attempt to parse the message
	if n < 1024 {
		buf = buf[:n]
	}
	message, err := parseMessage(buf)

	// continue down processing pipeline
	if err != nil {
		Logger().Warning.Println("Failed to parse message", err.Error())
		// TODO: determine how to send proper error response
		conn.Write([]byte("Message could not be processed."))
	} else {
		Logger().Trace.Printf("Processing message: %v\n", message)
		messageChannel <- message
	}
}

func startGateway() chan *Message {
	// channel for listening to cache expiration
	expiryChan := make(chan *Message, messageExpirationChanBufferSize)
	registry := newRegistry(messageCacheTTLSeconds, expiryChan)
	gateway := newGateway(registry)

	ch := make(chan *Message, gatewayMessageBufferSize)
	go gateway.run(ch)

	return ch
}
