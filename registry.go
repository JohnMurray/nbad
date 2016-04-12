package main

/**
 * File: registry.go
 *
 * The registry contains a local cache table of recent messages that have flowed
 * through the agent. What exists in the registry for any given event can be thought
 * of the last known state of the message / event. This local cache aids in decision
 * making of what should be sent up to the master Nagios host.
 */

import (
	"fmt"
	"time"

	"github.com/JohnMurray/nbad/flapper"
)

// Registry is just a fancy cache with a TTL
type Registry struct {
	// cache of messages
	cache map[string]*MessageEntry

	// how long before a message should be expired from the cache
	ttlInSeconds uint

	// how long a message is initially buffered before it can be decisioned on
	initBufferTTLInSeconds uint
}

// MessageEntry is something to store in the Registry
type MessageEntry struct {
	message            *Message
	prevMessage        *Message
	initBufferExpireAt time.Time
	expireAt           time.Time
	flap               *flapper.Flapper
}

// Contains checks to see if the message is currently in the registry.
func (r *Registry) contains(message *Message) bool {
	if _, ok := r.cache[message.Service]; ok {
		return true
	}
	return false
}

// Update stores message in the registry or updates it if it's already there
func (r *Registry) update(message *Message) {
	me := &MessageEntry{
		message:            message,
		expireAt:           time.Now().Add(time.Duration(r.ttlInSeconds) * time.Second),
		initBufferExpireAt: time.Now().Add(time.Duration(Config().MessageInitBufferTimeSeconds) * time.Second),
		flap:               flapper.NewFlapper(Config().FlapCountThreshold, Config().MessageInitBufferTimeSeconds),
	}
	if prev := r.get(message.Service); prev != nil {
		me.prevMessage = prev
	}
	r.cache[message.Service] = me
}

func (r *Registry) get(key string) *Message {
	if ce, ok := r.cache[key]; ok {
		return ce.message
	}
	return nil
}

func (r *Registry) getPrev(key string) *Message {
	if ce, ok := r.cache[key]; ok {
		return ce.prevMessage
	}
	return nil
}

func (r *Registry) getFlap(key string) *flapper.Flapper {
	if ce, ok := r.cache[key]; ok {
		return ce.flap
	}
	return nil
}

func (r *Registry) summaryString() string {
	s := ""
	for k, v := range r.cache {
		entry := fmt.Sprintf("\t%s | %s | %s\n", k, stateName(v.message.State), v.message.Message)
		s = s + entry
	}
	return s
}
