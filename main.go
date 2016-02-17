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

	gatewayMessageBufferSize    = 100
	messageCacheTTLSeconds      = 60
	messageInitBufferTTLSeconds = 10
)

// TODO add support for configuring all these stupid const things
// TODO add command-line option support for stuffs (not sure what yet)

func main() {
	Logger().Info.Println("Starting up NBAd")

	address := connHost + ":" + connPort

	listener, err := net.Listen(connType, address)
	if err != nil {
		Logger().Error.Println("Could not bind to "+address, err.Error())
		os.Exit(errBinding)
	}
	Logger().Info.Printf("Listening at %s\n", address)

	// close listener on program exit
	defer listener.Close()

	// sping up message registry
	messageChannel := startGateway()

	// listen for incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			Logger().Error.Println("Error accepting connection", err.Error())
			// FIXME should we do more than just print out an error here?
		}
		go handleRequest(conn, messageChannel)
	}
}

// handles incoming requests
func handleRequest(conn net.Conn, messageChannel chan *GatewayEvent) {
	defer conn.Close()

	// TODO send an initialization message (see https://github.com/Syncbak-Git/nsca/blob/master/packet.go#L163)

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
		messageChannel <- newMessageEvent(message)
	}
}

// Starts a gateway process. Returns a channel to send new messages to the gateway.
func startGateway() chan *GatewayEvent {
	// channel for sending new messages to the Gateway
	gatewayChan := make(chan *GatewayEvent, gatewayMessageBufferSize)

	registry := newRegistry(messageInitBufferTTLSeconds, messageCacheTTLSeconds, gatewayChan)
	gateway := newGateway(registry, gatewayChan)

	go gateway.run()

	return gatewayChan
}

func newMessageEvent(m *Message) *GatewayEvent {
	return &GatewayEvent{message: m}
}
