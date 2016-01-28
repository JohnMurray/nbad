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
 * Initiates the gateway that listens and processes messages.
 */
func (g *Gateway) run(incomingMessageChan chan *Message) {
	// setup subscription channel for registry timeout notifications
	go g.handleExpiry(g.registry.expireChan)

	// handle all incoming (new messages) from clients
	go g.handleIncomingMessages(incomingMessageChan)
}

/**
 * Listen specificially for incoming messages and process them.
 */
func (g *Gateway) handleIncomingMessages(ch chan *Message) {
	for {
		message := <-ch
		if g.registry.contains(message) {
			fmt.Println("duplicate message")
			Logger().Trace.Printf("Duplicate message: %v\n", message)
		}
		// TODO detect state-change, act appropriately
		g.registry.update(message)
	}
}

/**
 * Listen for messages that are expiring from the registry and take action. An expired message
 * is really just the last provided state for a a service / alert. If it's in an error or warning
 * state, then we may want to clear the event. If it's in an OK state, we can leave it be.
 */
func (g *Gateway) handleExpiry(expireNotificationChan chan *Message) {
	for {
		message := <-expireNotificationChan
		if message.State == stateOk {

		}
		switch message.State {
		case stateOk:
			break
		case stateWarning:
		case stateCritical:
		case stateUnknown:
			// TODO clear the state to the upstream nagios server
		}
	}
}

func newGateway(r *Registry) *Gateway {
	g := &Gateway{registry: r}
	return g
}
