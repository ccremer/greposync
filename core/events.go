package core

import (
	"fmt"
)

// EventSource represents a uniform input value for event processing.
type EventSource struct {
	// Url is a property that uniquely identifies the Git repository that was involved in the event handling.
	Url *GitURL
}

// EventResult represents a uniform return value for events.
type EventResult struct {
	// Url is a property that uniquely identifies the Git repository that was involved in the event handling.
	Url *GitURL
	// Error can optionally contain any error that occurred during event handling.
	Error error
}

// EventHandler is a service that can handle certain configured event types.
type EventHandler interface {
	// Handle processes the given EventSource synchronously.
	// Any errors are to be defined in the EventResult::Error property.
	Handle(source EventSource) EventResult
}

// EventName uniquely identifies a named event type.
// There should only be exactly one handler for each event name.
type EventName string

var handlers = map[EventName]EventHandler{}

// IsSuccessful returns true if the Error property is nil.
func (e EventResult) IsSuccessful() bool {
	return e.Error == nil
}

// FireEvent handles the given event in the background and returns a channel for retrieving the EventResult.
// If there is no handler for the given EventName, then an EventResult is returned with an error.
// There is only one EventResult and the channel is immediately closed after populating it.
func FireEvent(name EventName, source EventSource) chan EventResult {
	ch := make(chan EventResult, 1)
	go func() {
		defer close(ch)
		if handler, exists := handlers[name]; exists && handler != nil {
			ch <- handler.Handle(source)
			return
		}
		ch <- noHandlerExists(name, source)
	}()
	return ch
}

// RegisterHandler adds the given handler to the internal map.
// It overwrites an existing handler if the handler is already registered with the given name.
func RegisterHandler(name EventName, handler EventHandler) {
	handlers[name] = handler
}

func noHandlerExists(name EventName, source EventSource) EventResult {
	return EventResult{
		Url:   source.Url,
		Error: fmt.Errorf("no event handler exists for '%s'", name),
	}
}

// ToResult is a convenience method that converts the given EventSource to an EventResult with same properties.
func ToResult(source EventSource, err error) EventResult {
	return EventResult{
		Url:   source.Url,
		Error: err,
	}
}
