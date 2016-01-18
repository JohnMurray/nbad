package main

/**
 * File: registry.go
 *
 * The registry contains a local cache table of recent messages that have flowed
 * through the agent. This local cache aids in decision making of what should be
 * sent up to the master Nagios host.
 */

type Registry struct {
	cache map[string]*Message
}

func (r *Registry) Contains(message *Message) bool {
	if _, ok := r.cache[message.Service]; ok {
		return true
	} else {
		return false
	}
}

func (r *Registry) Update(message *Message) {
	r.cache[message.Service] = message
}
