package events

import (
	"maps"
	"reflect"
	"sync"
)

type EventBus struct {
	mu              sync.RWMutex
	handlers        map[reflect.Type][]func(any)
	deferred_events []any
}

// CreateEventBus creates a new event bus.
func CreateEventBus() *EventBus {
	return &EventBus{
		mu:       sync.RWMutex{},
		handlers: make(map[reflect.Type][]func(any)),
	}
}

// AddListener registers an event handler based on the parameter type.
func (e *EventBus) AddListener(handler func(any)) {
	e.mu.Lock()
	eventType := reflect.TypeOf(handler).In(0)
	e.handlers[eventType] = append(e.handlers[eventType], handler)
	e.mu.Unlock()
}

// EmitImmediate emits an event that is handled by immediately calling the handlers associated with it. I dont recommend using it.
func (e *EventBus) EmitImmediate(event any) {
	e.mu.RLock()
	h_arr, exists := e.handlers[reflect.TypeOf(event)]
	e.mu.RUnlock()
	if exists {
		for _, h := range h_arr {
			h(event)
		}
	}
}

// EmitDeferred emits an event that is added to a queue, and gets handled when somebody (should be the main thread) calls HandleDeferred.
func (e *EventBus) EmitDeferred(event any) {
	e.mu.Lock()
	e.deferred_events = append(e.deferred_events, event)
	e.mu.Unlock()
}

// HandleDeferred processes n number of events from the queue, in FIFO order.
func (e *EventBus) HandleDeferred(n int) {
	if n <= 0 {
		return
	}

	e.mu.Lock()
	if len(e.deferred_events) == 0 {
		e.mu.Unlock()
		return
	}

	var batch []any
	if n > 0 && n < len(e.deferred_events) {
		batch = e.deferred_events[:n]
		e.deferred_events = e.deferred_events[n:]
	} else {
		batch = e.deferred_events
		e.deferred_events = make([]any, 0)
	}
	e.mu.Unlock()

	e.mu.RLock()
	batch_handlers := make(map[reflect.Type][]func(any), len(e.handlers))
	maps.Copy(batch_handlers, e.handlers)
	e.mu.RUnlock()

	for _, event := range batch {
		h_arr, exists := batch_handlers[reflect.TypeOf(event)]
		if exists {
			for _, h := range h_arr {
				h(event)
			}
		}
	}
}
