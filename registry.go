package ocpp16json

import (
	"encoding/json"
	"sync"
)

// Registry is a thread-safe map of Action names to
// PayloadDecoders. It allows registering a decoder for each
// OCPP Action and later decoding payloads by Action name.
type Registry struct {
	mutex    sync.RWMutex
	decoders map[string]PayloadDecoder
}

// NewRegistry creates an empty Registry ready for use.
func NewRegistry() *Registry {
	return &Registry{
		mutex:    sync.RWMutex{},
		decoders: make(map[string]PayloadDecoder),
	}
}

// Register associates a PayloadDecoder with an Action name.
// It returns ErrActionAlreadyRegistered if the Action has
// already been registered.
func (registry *Registry) Register(
	action string,
	decoder PayloadDecoder,
) error {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	if _, exists := registry.decoders[action]; exists {
		return ErrActionAlreadyRegistered
	}

	registry.decoders[action] = decoder

	return nil
}

// Decode looks up the decoder for the given Action and
// applies it to the raw JSON payload. It returns
// ErrUnknownAction if no decoder is registered for the
// Action.
func (registry *Registry) Decode(
	action string,
	payload json.RawMessage,
) (any, error) {
	registry.mutex.RLock()
	decoder, exists := registry.decoders[action]
	registry.mutex.RUnlock()

	if !exists {
		return nil, ErrUnknownAction
	}

	return decoder(payload)
}
