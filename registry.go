package ocpp16json

import (
	"encoding/json"
	"fmt"
	"sync"
)

type PayloadDecoder func(json.RawMessage) (any, error)

type Registry struct {
	mu            sync.RWMutex
	requests      map[string]PayloadDecoder
	confirmations map[string]PayloadDecoder
}

func NewRegistry() *Registry {
	return &Registry{
		requests:      make(map[string]PayloadDecoder),
		confirmations: make(map[string]PayloadDecoder),
	}
}

func (r *Registry) RegisterRequest(action string, decoder PayloadDecoder) error {
	return r.register(action, decoder, r.requests)
}

func (r *Registry) RegisterConfirmation(
	action string,
	decoder PayloadDecoder,
) error {
	return r.register(action, decoder, r.confirmations)
}

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

func (r *Registry) DecodeCallResult(action string, data []byte) (DecodedCallResult, error) {
	if err := validateAction(action); err != nil {
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
	if err := validateAction(action); err != nil {
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
