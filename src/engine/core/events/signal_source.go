package events

import (
	"errors"
	"sync"
)

type HandlerId uint64

type EventHandler[T any] func(e T) error

type listener_t[T any] struct {
	id HandlerId
	fn EventHandler[T]
}

type EventSource[T any] struct {
	mu       sync.RWMutex
	next_id  HandlerId
	handlers []listener_t[T]
}

func NewSource[T any]() *EventSource[T] {
	return &EventSource[T]{
		mu:       sync.RWMutex{},
		next_id:  0,
		handlers: make([]listener_t[T], 0),
	}
}

func (e *EventSource[T]) handle_event(event T) error {
	err_list := []error{}
	for i := range e.handlers {
		err := e.handlers[i].fn(event)
		if err != nil {
			err_list = append(err_list, err)
		}
	}
	return errors.Join(err_list...)
}

func (e *EventSource[T]) EmitImmediate(event T) error {
	if e == nil {
		panic("Event source was a null (nil) pointer!")
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.handle_event(event)
}

func (e *EventSource[T]) EmitDeferred(event T) {
	if e == nil {
		panic("Event source was a null (nil) pointer!")
	}
	go e.EmitImmediate(event)
}

func (e *EventSource[T]) Subscribe(fn EventHandler[T]) HandlerId {
	if e == nil {
		panic("Event source was a null (nil) pointer!")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	id := e.next_id
	e.next_id++
	e.handlers = append(e.handlers, listener_t[T]{
		id: id,
		fn: fn,
	})
	return id
}

func (e *EventSource[T]) Unsubscribe(id HandlerId) {
	if e == nil {
		panic("Event source was a null (nil) pointer!")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	last := len(e.handlers) - 1
	for i := range e.handlers {
		if e.handlers[i].id == id {
			e.handlers[i] = e.handlers[last]
			e.handlers = e.handlers[:last]
			return
		}
	}
}
