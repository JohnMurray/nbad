package main

/**
 * File: registry.go
 *
 * The registry contains a local cache table of recent messages that have flowed
 * through the agent. This local cache aids in decision making of what should be
 * sent up to the master Nagios host.
 */

import (
	"time"
)

// Registry is just a fancy cache with a TTL
type Registry struct {
	// cache of messages
	cache map[string]*CacheEntry

	// how long before a message should be expired from the cache
	ttlInSeconds uint32
}

// CacheEntry is something to store in the Registry
type CacheEntry struct {
	message  *Message
	expireAt time.Time
}

func newRegistry(ttlInSeconds uint32) *Registry {
	registry := &Registry{
		cache:        make(map[string]*CacheEntry),
		ttlInSeconds: ttlInSeconds,
	}
	// TODO: start expiry go-routine here
	go registry.expireOldCache()
	return registry
}

// TODO: gateway needs to be notified of expiry events from the registry
func (r *Registry) expireOldCache() {
	interval := 100 * time.Millisecond
	for {
		time.Sleep(interval)

		now := time.Now()
		for k, v := range r.cache {
			if now.After(v.expireAt) {
				Logger().Trace.Printf("Expiring cache %s\n", k)
				delete(r.cache, k)
				// TODO: send notification
			}
		}
	}
}

// Contains checks to see if the message is currently in the registry.
func (r *Registry) Contains(message *Message) bool {
	if _, ok := r.cache[message.Service]; ok {
		return true
	}
	return false
}

// Update stores message in the registry or updates it if it's already there
func (r *Registry) Update(message *Message) {
	ce := &CacheEntry{
		message:  message,
		expireAt: time.Now().Add(time.Duration(r.ttlInSeconds) * time.Second),
	}
	r.cache[message.Service] = ce
}
