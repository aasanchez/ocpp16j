package ocpp16json

import (
	"encoding/json"
	"fmt"
	"sync"
)

// PayloadDecoder converts a raw payload into a validated typed value.
//
// Decoders are typically created with JSONDecoder and backed by constructors
// from github.com/aasanchez/ocpp16messages.
type PayloadDecoder func(json.RawMessage) (any, error)

// Registry stores request and confirmation payload decoders by action name.
//
// A Registry is safe for concurrent use.
type Registry struct {
	mu            sync.RWMutex
	requests      map[string]PayloadDecoder
	confirmations map[string]PayloadDecoder
}

// NewRegistry creates an empty action registry.
func NewRegistry() *Registry {
	return &Registry{
		mu:            sync.RWMutex{},
		requests:      make(map[string]PayloadDecoder),
		confirmations: make(map[string]PayloadDecoder),
	}
}

// RegisterRequest registers a decoder for a CALL action.
//
// It returns ErrInvalidAction for an empty action name,
// ErrActionAlreadyRegistered for duplicate registrations, and
// ErrPayloadDecode if decoder is nil.
func (r *Registry) RegisterRequest(
	action string,
	decoder PayloadDecoder,
) error {
	return r.register(action, decoder, r.requests)
}

// RegisterConfirmation registers a decoder for a CALLRESULT action.
//
// It returns ErrInvalidAction for an empty action name,
// ErrActionAlreadyRegistered for duplicate registrations, and
// ErrPayloadDecode if decoder is nil.
func (r *Registry) RegisterConfirmation(
	action string,
	decoder PayloadDecoder,
) error {
	return r.register(action, decoder, r.confirmations)
}

// DecodeCall parses and decodes a CALL frame.
//
// It validates the envelope with Parse, selects the decoder registered for the
// frame action, and returns a DecodedCall with a typed payload value.
func (r *Registry) DecodeCall(data []byte) (DecodedCall, error) {
	frame, err := Parse(data)
	if err != nil {
		return DecodedCall{}, err
	}

	call, ok := frame.(RawCall)
	if !ok {
		return DecodedCall{}, ErrInvalidFrame
	}

	payload, err := r.decodeRegistered(call.Action, call.Payload, r.requests)
	if err != nil {
		return DecodedCall{}, err
	}

	return DecodedCall{
		UniqueID: call.UniqueID,
		Action:   call.Action,
		Payload:  payload,
	}, nil
}

// DecodeCallResult parses and decodes a CALLRESULT frame.
//
// The action must be provided by the caller because OCPP-J CALLRESULT frames
// do not include it on the wire.
func (r *Registry) DecodeCallResult(
	action string,
	data []byte,
) (DecodedCallResult, error) {
	err := validateAction(action)
	if err != nil {
		return DecodedCallResult{}, err
	}

	frame, err := Parse(data)
	if err != nil {
		return DecodedCallResult{}, err
	}

	result, ok := frame.(RawCallResult)
	if !ok {
		return DecodedCallResult{}, ErrInvalidFrame
	}

	payload, err := r.decodeRegistered(action, result.Payload, r.confirmations)
	if err != nil {
		return DecodedCallResult{}, err
	}

	return DecodedCallResult{
		UniqueID: result.UniqueID,
		Action:   action,
		Payload:  payload,
	}, nil
}

func (r *Registry) decodeRegistered(
	action string,
	raw json.RawMessage,
	decoders map[string]PayloadDecoder,
) (any, error) {
	r.mu.RLock()

	decoder, ok := decoders[action]

	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownAction, action)
	}

	payload, err := decoder(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", ErrPayloadDecode, action, err)
	}

	return payload, nil
}

func (r *Registry) register(
	action string,
	decoder PayloadDecoder,
	target map[string]PayloadDecoder,
) error {
	err := validateAction(action)
	if err != nil {
		return err
	}

	if decoder == nil {
		return fmt.Errorf("%w: nil decoder for %s", ErrPayloadDecode, action)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := target[action]; exists {
		return fmt.Errorf("%w: %s", ErrActionAlreadyRegistered, action)
	}

	target[action] = decoder

	return nil
}
