// Package message provides message parsing and serializing
package message

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"

	"github.com/JohnMurray/nbad/log"
)

const (
	// StateOk OK state for Nagios
	StateOk = iota
	// StateWarning WARN state for Nagios
	StateWarning
	// StateCritical CRIT state for Nagios
	StateCritical
	// StateUnknown  state for Nagios
	StateUnknown

	nagiosMessageLen = 720
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

// ParseMessage parses byte arrays to Nagios Message v3 spec (or as close as I can get)
func ParseMessage(bytes []byte) (*Message, error) {
	if len(bytes) >= 2 {
		version := binary.BigEndian.Uint16(bytes[:2])
		if version != 3 {
			return nil, fmt.Errorf("Can only handle message version 3, %d received", version)
		}
	}

	// // ensure we're dealing with a proper v3 message via length
	if len(bytes) != nagiosMessageLen {
		return nil, fmt.Errorf("Expected message of %d bytes, received %d", nagiosMessageLen, len(bytes))
	}

	// discard CRC for now. not sure what to do with it just yet
	crc := binary.BigEndian.Uint32(bytes[4:8])
	for i := 4; i <= 8; i++ {
		bytes[i] = 0
	}
	// only validate if a CRC is provided
	if crc != 0 {
		if valid, calcCrc := validateCrc(bytes, crc); !valid {
			return nil, fmt.Errorf("CRC validation failed. Given %d, but calculated %d", crc, calcCrc)
		}
	}

	// read the timestamp
	timestamp := binary.BigEndian.Uint32(bytes[8:12])
	// TODO: validate timestamp as semi-current ?? (maybe?)

	// read the return-code (state)
	returnCode := binary.BigEndian.Uint16(bytes[12:14])
	if returnCode != StateOk &&
		returnCode != StateWarning &&
		returnCode != StateCritical &&
		returnCode != StateUnknown {

		log.Trace().Printf("Unknown return code received %d", returnCode)
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

func validateCrc(message []byte, crc uint32) (bool, uint32) {
	calcCrc := crc32.ChecksumIEEE(message)
	return calcCrc == crc, calcCrc
}

// StateName returns a string-representation of state
func StateName(state uint16) string {
	switch state {
	case StateOk:
		return "OK"
	case StateWarning:
		return "WARNING"
	case StateCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}
