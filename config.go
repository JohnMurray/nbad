package main

import (
	"encoding/json"
	"os"
	"sync"
)

// TODO this should really be like /etc/nbad/conf.json or something real
const defaultConfLocation string = "conf.json"

// NbadConfig is just the struct that holds all of the config values
type NbadConfig struct {
	// GatewayMessageBufferSize - The number of messages to buffer in memory for the gateway
	GatewayMessageBufferSize int `json:"gateway_message_buffer_size"`

	// MessageCacheTTLInSeconds - The time before a message expires (possibly causing reset of upstream state)
	MessageCacheTTLInSeconds int `json:"message_cache_ttl_in_seconds"`

	// MessageInitBufferTimeSeconds - The amount of time a message is buffered before actioned upon
	MessageInitBufferTimeSeconds int `json:"message_init_buffer_ttl_in_seconds"`
}

var configLoadOnce sync.Once
var nbadConfig *NbadConfig

// Config returns the current config file, loading it if necessary
func Config() *NbadConfig {
	configLoadOnce.Do(func() {
		loadConfigFile()
	})

	Logger().Trace.Printf("Loaded config: %v\n", nbadConfig)

	return nbadConfig
}

func loadConfigFile() {
	// TODO make this file configurable via command line
	file, err := os.Open(defaultConfLocation)
	if err != nil {
		Logger().Error.Fatalf("could not load config file '%s': %v\n", defaultConfLocation, err)
	}

	decoder := json.NewDecoder(file)
	configuration := &NbadConfig{}
	err = decoder.Decode(configuration)
	if err != nil {
		Logger().Error.Fatalf("could not read config file '%s': %v", defaultConfLocation, err)
	}

	nbadConfig = configuration
}
