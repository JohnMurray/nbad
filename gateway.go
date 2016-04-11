package main

/**
 * File: gateway.go
 *
 * The gateway is the main processing center for NBA. This is where the decisions
 * are made on which alerts are proxied, cached, discarded, set OK, etc. The Gateway
 * works by listening to a series of events. Current events that are listened to include:
 *
 *  - NSCA message received from client
 *  - StateExpiry event received from registry
 *  - InitBufferExpiry event received from registry
 */

import (
	"sync"
	"time"

	"github.com/JohnMurray/nbad/flapper"
)

// Gateway is where all the messages flow through
type Gateway struct {
	registry          *Registry
	incomingEventChan chan *GatewayEvent
	startOnce         sync.Once
}

// GatewayEvent represents union of events a Gateway expects to receive
type GatewayEvent struct {
	message          *Message
	stateExpiry      *StateExpiry
	initBufferExpiry *InitBufferExpiry
}

// StateExpiry is an event raised when the current service state expires (no recent message)
type StateExpiry struct{ service string }

// InitBufferExpiry is an event raised when buffering of initial server state has been reached
type InitBufferExpiry struct{ service string }

/**
 * Initiates the gateway that listens and processes various types of events.
 */
func (g *Gateway) run() {
	g.startOnce.Do(func() {
		go g.handleIncomingEvents()
	})
}

// handleIncomingEvents - listen for incoming events and dispatch them to the appropriate handling code
func (g *Gateway) handleIncomingEvents() {
	_flapper := flapper.NewFlapper(1, 2)
	_flapper.NoteStateChange("test")
	tick := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-tick:
			g.expireOldMessages()
		case event := <-g.incomingEventChan:
			if event == nil {
				continue
			}
			g.handleMessageStateChange(event)
		}
	}
}

// expireOldMessages - scan registry and emit events for message expirations
func (g *Gateway) expireOldMessages() {
	now := time.Now()
	for _, v := range g.registry.cache {
		if now.After(v.expireAt) {
			// send notification of message expiration
			g.incomingEventChan <- &GatewayEvent{stateExpiry: &StateExpiry{service: v.message.Service}}
		} else if now.After(v.initBufferExpireAt) {
			// send notification of init-buffer expiry
			g.incomingEventChan <- &GatewayEvent{initBufferExpiry: &InitBufferExpiry{service: v.message.Service}}
		}
	}
}

func (g *Gateway) handleMessageStateChange(event *GatewayEvent) {
	if event.message != nil {
		/*
		 * The event is an incoming message. All incoming messages should be buffered for a small
		 * period of time to make sure we're not thrashing (flip-flopping). However we have to be careful
		 * so that a flooding scenario doesn't cause us to stall indefinitely. Thus the following rules
		 * can be applied:
		 *   - if no previous service alert, store
		 *   - if previous service alert with same state (OK, WARN, etc), discard current message, don't update TTLs
		 *   - if previous service alert is different, update message and all TTLs
		 */
		// TODO ^^ these rules don't actually solve the problem... so... that's lame.  -___-
		// TODO - Update logic to reflect above description
		g.registry.update(event.message)
		Logger().Trace.Printf("registry:\n%s\n", g.registry.summaryString())
	} else if event.initBufferExpiry != nil {
		/*
		 * All messages are given an initial buffering time. This event is raised when that time is up.
		 * At this point we need to make a decision based on the sate of the message. In general we do:
		 *   - if previous state is different, proxy
		 *   - if previous state is the same, do nothing
		 *   - if previous state does not exist (expired or new), proxy
		 */
		if message := g.registry.get(event.initBufferExpiry.service); message != nil {
			if previous := g.registry.getPrev(event.initBufferExpiry.service); previous != nil {
				if message.State != previous.State {
					Logger().Info.Printf("detected state change from %s to %s for service %s",
						stateName(previous.State), stateName(message.State), message.Service)
				}
			} else {
				Logger().Info.Printf("new state of %s for service %s, sending upstream",
					stateName(message.State), message.Service)
			}
		}
	} else if event.stateExpiry != nil {
		/*
		 * An alert is cached based on it's last recorded state. When a service has not had any activity
		 * in a while, it will eventually expire with it's last known state. That is when this event is raised
		 * and we can take action on it. If an event expires in an error-state, we can set it back to a
		 * non-error state.
		 *
		 * A non-error state is defined as OK here.
		 */
		if message := g.registry.get(event.stateExpiry.service); message != nil {
			Logger().Info.Printf("expired message: %v with state %s\n", message, stateName(message.State))
			switch message.State {
			case stateOk: // do nothing
			case stateWarning:
				fallthrough
			case stateCritical:
				Logger().Info.Printf("PUSH Sending state '%s' for expired service '%s' upstream",
					stateName(stateOk), message.Service)
			default:
				Logger().Trace.Println("Expired message in UNKNOWN state")
			}
		}
	}
}

func newGateway(r *Registry, incomingEventChan chan *GatewayEvent) *Gateway {
	g := &Gateway{
		registry:          r,
		incomingEventChan: incomingEventChan,
	}
	return g
}
