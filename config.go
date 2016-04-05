package main

/**
 * File: config.go
 *
 * Defines a simply way to access config values from anywhere in the program by simply
 * calling `Config()`. Takes care of lazily loading configuration if it has not already
 * been loaded.
 */

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// NbadConfig is just the struct that holds all of the config values
type NbadConfig struct {
	// GatewayMessageBufferSize - The number of messages to buffer in memory for the gateway
	GatewayMessageBufferSize uint32 `json:"gateway_message_buffer_size"`

	// MessageCacheTTLInSeconds - The time before a message expires (possibly causing reset of upstream state)
	MessageCacheTTLInSeconds uint32 `json:"message_cache_ttl_in_seconds"`

	// MessageInitBufferTimeSeconds - The amount of time a message is buffered before actioned upon
	MessageInitBufferTimeSeconds uint32 `json:"message_init_buffer_ttl_in_seconds"`

	// TraceLogging - Enable trace logging (for debugging purposes) (not in JSON config file)
	TraceLogging bool
}

var configLoadOnce sync.Once
var nbadConfig *NbadConfig

// Config returns the current config file
func Config() *NbadConfig {
	return nbadConfig
}

// InitConfig - loads the config file
func InitConfig(configFile string, logger *log.Logger) {
	configLoadOnce.Do(func() {
		loadConfigFile(configFile, logger)
		validateConfig(logger)
	})

	logger.Printf("Loaded config: %v\n", nbadConfig)
}

func loadConfigFile(confFile string, logger *log.Logger) {
	logger.Printf("Loading config from file '%s'\n", confFile)

	file, err := os.Open(confFile)
	if err != nil {
		logger.Fatalf("could not load config file '%s': %v\n", confFile, err)
	}

	decoder := json.NewDecoder(file)
	configuration := &NbadConfig{}
	err = decoder.Decode(configuration)
	if err != nil {
		logger.Fatalf("could not read config file '%s': %v", confFile, err)
	}

	nbadConfig = configuration
}

func validateConfig(logger *log.Logger) {
	c := nbadConfig

	if c.MessageInitBufferTimeSeconds > c.MessageCacheTTLInSeconds {
		logger.Fatalln("init buffer ttl cannot be greater than message cache ttl")
	}
}
