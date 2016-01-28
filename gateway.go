package main

/**
 * File: gateway.go
 *
 * The gateway is the main processing center for NBA. This is where the decisions
 * are made on which alerts are proxied, cached, discarded, set OK, etc.
 */

import (
	"fmt"
)

// Gateway is where all the messages flow through
type Gateway struct {
	registry *Registry
}

/**
 * Listens a channel for Message's and actions on them.
 */
func (g *Gateway) run(ch chan *Message) {
	for {
		message := <-ch
		if g.registry.Contains(message) {
			fmt.Println("duplicate message")
			Logger().Trace.Printf("Duplicate message: %v\n", message)
		}
		g.registry.Update(message)
	}
}

func newGateway(r *Registry) *Gateway {
	g := &Gateway{registry: r}
	return g
}
