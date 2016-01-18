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

type Gateway struct {
	registry *Registry
}

/**
 * Listens a channel for Message's and actions on them.
 */
func (g *Gateway) run(ch chan *Message) {
	for {
		message := <-ch
		fmt.Printf("Received message: %v\n", message)
		if g.registry.Contains(message) {
			fmt.Println("duplicate message")
		}
		g.registry.Update(message)
	}
}
