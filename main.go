package main

import (
	"net"
	"os"

	"github.com/codegangsta/cli"
)

const (
	connType = "tcp"
	connHost = "localhost"
	connPort = "5667"

	defaultConfLocation = "/etc/nbad/conf.json"

	errBinding           = 1
	errAccptIncomingConn = 2
)

func main() {
	app := cli.NewApp()
	app.Name = "nbad"
	app.Usage = "NSCA Buffering Agent (daemon) - Emulates NSCA interface as local buffer/proxy"

	var configFile string
	var trace bool

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       defaultConfLocation,
			Usage:       "Location of config file on disk",
			Destination: &configFile,
		},
		cli.BoolFlag{
			Name:        "trace, t",
			EnvVar:      "NBAD_TRACE",
			Usage:       "Turn on trace-logging",
			Destination: &trace,
		},
	}
	app.Version = "1.0"
	app.Action = func(c *cli.Context) {
		// load configuration
		InitConfig(configFile, TempLogger("STARTUP"))
		Config().TraceLogging = trace

		startServer()
	}
	app.Run(os.Args)
}

func startServer() {
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
		go handleIncomingConn(conn, messageChannel)
	}
}

// handles incoming requests
func handleIncomingConn(conn net.Conn, messageChannel chan *GatewayEvent) {
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
	conn.Close()
}

// Starts a gateway process. Returns a channel to send new messages to the gateway.
func startGateway() chan *GatewayEvent {
	// channel for sending new messages to the Gateway
	gatewayChan := make(chan *GatewayEvent, Config().GatewayMessageBufferSize)

	registry := &Registry{
		cache:                  make(map[string]*MessageEntry),
		ttlInSeconds:           Config().MessageCacheTTLInSeconds,
		initBufferTTLInSeconds: Config().MessageInitBufferTimeSeconds,
	}
	gateway := newGateway(registry, gatewayChan)

	go gateway.run()

	return gatewayChan
}

func newMessageEvent(m *Message) *GatewayEvent {
	return &GatewayEvent{message: m}
}
