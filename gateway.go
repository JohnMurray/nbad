package main

/**
 * File: gateway.go
 *
 * The gateway is the main processing center for NBA. This is where the decisions
 * are made on which alerts are proxied, cached, discarded, set OK, etc.
 */

// Gateway is where all the messages flow through
type Gateway struct {
	registry            *Registry
	incomingMessageChan chan *Message
}

/**
 * Initiates the gateway that listens and processes messages.
 */
func (g *Gateway) run() {
	// setup subscription channel for registry timeout notifications
	go g.handleExpiry(g.registry.expireChan)

	// handle all incoming (new messages) from clients
	go g.handleIncomingMessages(g.incomingMessageChan)
}

/**
 * Listen specificially for incoming messages and process them.
 */
func (g *Gateway) handleIncomingMessages(ch chan *Message) {
	for {
		message := <-ch
		if oldMessage := g.registry.get(message.Service); oldMessage != nil {
			if message.State > oldMessage.State {
				// TODO things got worse, what now?
				// send upstream
				Logger().Info.Printf("PUSH Sending state '%s' for '%s' upstream",
					stateName(message.State), message.Service)
			} else if message.State < oldMessage.State {
				// TODO things got better, what now?
				// send upstream
				Logger().Info.Printf("PUSH Sending state '%s' for '%s' upstream",
					stateName(message.State), message.Service)
			} else {
				// it's the same... so we're not gonna do anything but buffer the
				// message to the TTLs can be updated and what not
				Logger().Trace.Printf("HOLD Duplicate message for service: %s\n", message.Service)
			}
		} else {
			Logger().Trace.Printf("PUSH-NEW Sending state '%s' for '%s' upstream",
				stateName(message.State), message.Service)
		}
		g.registry.update(message)
		Logger().Trace.Printf("registry:\n%s\n", g.registry.summaryString())
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
		Logger().Info.Printf("expired message: %v with state %s\n", message, stateName(message.State))
		if message.State == stateOk {

		}
		switch message.State {
		case stateOk:
			// do nothing
		case stateWarning:
			fallthrough
		case stateCritical:
			// TODO clear the state to the upstream nagios server
			Logger().Info.Printf("PUSH Sending state '%s' for expired service '%s' upstream",
				stateName(stateOk), message.Service)
		default:
			Logger().Trace.Println("Expired message in UNKNOWN state")
		}
	}
}

func newGateway(r *Registry, incomingMessageChan chan *Message) *Gateway {
	g := &Gateway{
		registry:            r,
		incomingMessageChan: incomingMessageChan,
	}
	return g
}
