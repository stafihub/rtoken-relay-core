// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package core

import (
	"fmt"
	"sync"

	log "github.com/ChainSafe/log15"
)

// Handler consumes a message and makes the requried on-chain interactions.
type Handler interface {
	HandleMessage(msg *Message)
}

// Router forwards messages from their source to their destination
type Router struct {
	registry map[RSymbol]Handler
	lock     *sync.RWMutex
	log      log.Logger
	stop     chan int
}

func NewRouter(log log.Logger) *Router {
	return &Router{
		registry: make(map[RSymbol]Handler),
		lock:     &sync.RWMutex{},
		log:      log,
		stop:     make(chan int),
	}
}

// Send passes a message to the destination Writer if it exists
func (r *Router) Send(msg *Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if msg.Reason != NewEra {
		r.log.Trace("Routing message", "source", msg.Source, "dest", msg.Destination, "Reason", msg.Reason)
	}

	h := r.registry[msg.Destination]
	if h == nil {
		return fmt.Errorf("unknown destination symbol: %s", msg.Destination)
	}

	go h.HandleMessage(msg)
	return nil
}

// Listen registers a Writer with a ChainId which Router.Send can then use to propagate messages
func (r *Router) Listen(symbol RSymbol, w Handler) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.log.Debug("Registering new chain in router", "symbol", symbol)
	r.registry[symbol] = w
}

func (r *Router) StopMsgHandler() {
	close(r.stop)
}
