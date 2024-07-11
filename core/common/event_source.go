package common

import (
	"sync"

	"github.com/google/uuid"
)

type EventSource[T any] struct {
	source      chan T
	subscribers map[uuid.UUID]chan T
	mtx         sync.RWMutex
	skipFn      func(T) bool
}

func (source *EventSource[T]) Subscribe() (uuid.UUID, <-chan T, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, nil, err
	}
	subscriber := make(chan T, 100)

	source.mtx.Lock()
	defer source.mtx.Unlock()
	source.subscribers[id] = subscriber
	return id, subscriber, nil
}

func (source *EventSource[T]) Unsubscribe(id uuid.UUID) {
	source.mtx.Lock()
	defer source.mtx.Unlock()
	close(source.subscribers[id])
	delete(source.subscribers, id)
}

func (source *EventSource[T]) GetSourceChannel() chan<- T {
	return source.source
}

func (source *EventSource[T]) notifySubcribers(event T) {
	source.mtx.RLock()
	defer source.mtx.RUnlock()
	for _, subscriber := range source.subscribers {
		go func() {
			subscriber <- event
		}()
	}
}

func (source *EventSource[T]) Run() {
	go func() {
		for event := range source.source {
			if source.skipFn != nil && source.skipFn(event) {
				continue
			}

			source.notifySubcribers(event)
		}
	}()
}

func NewEventSource[T any](skipFn func(T) bool) *EventSource[T] {
	return &EventSource[T]{
		source:      make(chan T),
		subscribers: make(map[uuid.UUID]chan T),
		skipFn:      skipFn,
	}
}
