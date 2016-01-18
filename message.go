package main

/**
 * File: message.go
 *
 * This file contains functions related to parsing and processing messages
 * incoming from the client. Sending messages upstream is handled by another
 * library that is much more complete in it's ability to _compose_ messages.
 * However since decomposing messages is usually a server-side thing, we have
 * our own message stuffs here.
 */

import (
	"encoding/binary"
	"fmt"
)

const (
	STATE_OK = iota
	STATE_WARNING
	STATE_CRITICAL
	STATE_UNKNOWN

	NAGIOS_MESSAGE_LEN = 720
)

// Message is the contents of an NSCA message
type Message struct {
	// when was the message raised
	Timestamp uint32
	// State is one of {STATE_OK, STATE_WARNING, STATE_CRITICAL, STATE_UNKNOWN}
	State uint16
	// Host is the host name to set for the NSCA message
	Host string
	// Service is the service name to set for the NSCA message [optional]
	Service string
	// Message is the "plugin output" of the NSCA message [optional]
	Message string
}

func ParseMessage(bytes []byte) (*Message, error) {
	fmt.Println("length", len(bytes))
	fmt.Printf("%v\n", bytes)
	if len(bytes) >= 2 {
		version := binary.BigEndian.Uint16(bytes[:2])
		if version != 3 {
			return nil, fmt.Errorf("Can only handle message version 3, %d received", version)
		}
	}

	// ensure we're dealing with a proper v3 message via length
	if len(bytes) != NAGIOS_MESSAGE_LEN {
		return nil, fmt.Errorf("Expected message of %d bytes, received %d", NAGIOS_MESSAGE_LEN, len(bytes))
	}

	// discard CRC for now. not sure what to do with it just yet
	// todo: figure out the right way to validate this
	binary.BigEndian.Uint32(bytes[4:8])

	// read the timestamp
	timestamp := binary.BigEndian.Uint32(bytes[8:12])
	// todo: validate timestamp as semi-current

	// read the return-code (state)
	returnCode := binary.BigEndian.Uint16(bytes[12:14])
	if returnCode != STATE_OK && returnCode != STATE_WARNING && returnCode != STATE_CRITICAL && returnCode != STATE_UNKNOWN {
		return nil, fmt.Errorf("Unknown return code received %d", returnCode)
	}

	// read hostname (64 bytes)
	hostname := string(bytes[14:78])

	// read service description / name (128 bytes)
	service := string(bytes[78:206])

	// read the description (512 bytes)
	description := string(bytes[206:718])

	// last two bytes are padding so we don't have to worry about them too much

	return &Message{
		Timestamp: timestamp,
		State:     returnCode,
		Host:      hostname,
		Service:   service,
		Message:   description,
	}, nil
}
