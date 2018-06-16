package sse

import (
	"github.com/dustin/go-broadcast"
)

// Keep track of active channels
var routeChannels = make(map[string]broadcast.Broadcaster)

// Keep track of when to remove broadcasters
var activeListeners = make(map[string]int)

// OpenListener registers a listener for a given
// route. The channel to pass data on is returned
func OpenListener(route string) chan interface{} {
	listener := make(chan interface{})
	URL(route).Register(listener)
	activeListeners[route]++
	return listener
}

// CloseListener deregisters a listener for a given URL
// from its associated channel, and closes its channel.
//
// If there are no remaining listeners on the channel,
// it also closes the channel.
func CloseListener(route string, listener chan interface{}) {
	URL(route).Unregister(listener)
	activeListeners[route]--
	close(listener)

	// Remove channel if no listeners left
	if activeListeners[route] == 0 {
		deleteBroadcast(route)
	}
}

func deleteBroadcast(route string) {
	// Close broadcast
	b, ok := routeChannels[route]
	if ok {
		b.Close()
		delete(routeChannels, route)
	}

	// Remove counter for broadcast
	_, ok = activeListeners[route]
	if ok {
		delete(activeListeners, route)
	}
}

// URL gets a channel for the given route.
// If a channel does not already exist for
// the route, one is created.
func URL(route string) broadcast.Broadcaster {
	b, ok := routeChannels[route]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		routeChannels[route] = b
		activeListeners[route] = 0
	}
	return b
}
