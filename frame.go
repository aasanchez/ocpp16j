package ocpp16json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	callFrameLength       = 4
	callResultFrameLength = 3
	callErrorFrameLength  = 5
)

type MessageType uint8

const (
	MessageTypeCall       MessageType = 2
	MessageTypeCallResult MessageType = 3
	MessageTypeCallError  MessageType = 4
)

type Frame interface {
	MessageType() MessageType
	MessageID() string
}

type RawCall struct {
	UniqueID string
	Action   string
	Payload  json.RawMessage
}

type RawCallResult struct {
	UniqueID string
	Payload  json.RawMessage
}

type CallError struct {
	UniqueID         string
	ErrorCode        string
	ErrorDescription string
	ErrorDetails     map[string]any
}

type DecodedCall struct {
	UniqueID string
	Action   string
	Payload  any
}

type DecodedCallResult struct {
	UniqueID string
	Action   string
	Payload  any
}

func (f RawCall) MessageType() MessageType { return MessageTypeCall }
func (f RawCall) MessageID() string        { return f.UniqueID }
func (f RawCallResult) MessageType() MessageType {
	return MessageTypeCallResult
}
func (f RawCallResult) MessageID() string { return f.UniqueID }
func (f CallError) MessageType() MessageType {
	return MessageTypeCallError
}
func (f CallError) MessageID() string          { return f.UniqueID }
func (f DecodedCall) MessageType() MessageType { return MessageTypeCall }
func (f DecodedCall) MessageID() string        { return f.UniqueID }
func (f DecodedCallResult) MessageType() MessageType {
	return MessageTypeCallResult
}
func (f DecodedCallResult) MessageID() string { return f.UniqueID }

func Parse(data []byte) (Frame, error) {
	var elements []json.RawMessage

	if err := json.Unmarshal(data, &elements); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFrame, err)
	}

	if len(elements) == 0 {
		return nil, ErrInvalidFrame
	}

	messageType, err := decodeMessageType(elements[0])
	if err != nil {
		return nil, err
	}

	switch messageType {
	case MessageTypeCall:
		return parseCall(elements)
	case MessageTypeCallResult:
		return parseCallResult(elements)
	}

	return parseCallError(elements)
}

func (f RawCall) MarshalJSON() ([]byte, error) {
	if err := validateMessageID(f.UniqueID); err != nil {
		return nil, err
	}

	if err := validateAction(f.Action); err != nil {
		return nil, err
	}

	if len(bytes.TrimSpace(f.Payload)) == 0 {
		return nil, ErrPayloadRequired
	}

	return marshalJSONArray(
		MessageTypeCall,
		f.UniqueID,
		f.Action,
		json.RawMessage(f.Payload),
	)
}

func (f RawCallResult) MarshalJSON() ([]byte, error) {
	if err := validateMessageID(f.UniqueID); err != nil {
		return nil, err
	}

	if len(bytes.TrimSpace(f.Payload)) == 0 {
		return nil, ErrPayloadRequired
	}

	return marshalJSONArray(
		MessageTypeCallResult,
		f.UniqueID,
		json.RawMessage(f.Payload),
	)
}

func (f CallError) MarshalJSON() ([]byte, error) {
	if err := validateMessageID(f.UniqueID); err != nil {
		return nil, err
	}

	if f.ErrorCode == "" {
		return nil, ErrErrorCodeRequired
	}

	if f.ErrorDescription == "" {
		return nil, ErrErrorDescriptionAbsent
	}

	if f.ErrorDetails == nil {
		f.ErrorDetails = map[string]any{}
	}

	return marshalJSONArray(
		MessageTypeCallError,
		f.UniqueID,
		f.ErrorCode,
		f.ErrorDescription,
		f.ErrorDetails,
	)
}

func parseCall(elements []json.RawMessage) (RawCall, error) {
	if len(elements) != callFrameLength {
		return RawCall{}, fmt.Errorf(
			"%w: CALL frame must have %d elements",
			ErrInvalidFrame,
			callFrameLength,
		)
	}

	id, err := decodeString(elements[1], ErrInvalidMessageID)
	if err != nil {
		return RawCall{}, err
	}

	action, err := decodeString(elements[2], ErrInvalidAction)
	if err != nil {
		return RawCall{}, err
	}

	if len(bytes.TrimSpace(elements[3])) == 0 {
		return RawCall{}, ErrPayloadRequired
	}

	return RawCall{
		UniqueID: id,
		Action:   action,
		Payload:  elements[3],
	}, nil
}

func parseCallResult(elements []json.RawMessage) (RawCallResult, error) {
	if len(elements) != callResultFrameLength {
		return RawCallResult{}, fmt.Errorf(
			"%w: CALLRESULT frame must have %d elements",
			ErrInvalidFrame,
			callResultFrameLength,
		)
	}

	id, err := decodeString(elements[1], ErrInvalidMessageID)
	if err != nil {
		return RawCallResult{}, err
	}

	if len(bytes.TrimSpace(elements[2])) == 0 {
		return RawCallResult{}, ErrPayloadRequired
	}

	return RawCallResult{
		UniqueID: id,
		Payload:  elements[2],
	}, nil
}

func parseCallError(elements []json.RawMessage) (CallError, error) {
	if len(elements) != callErrorFrameLength {
		return CallError{}, fmt.Errorf(
			"%w: CALLERROR frame must have %d elements",
			ErrInvalidFrame,
			callErrorFrameLength,
		)
	}

	id, err := decodeString(elements[1], ErrInvalidMessageID)
	if err != nil {
		return CallError{}, err
	}

	errorCode, err := decodeString(elements[2], ErrErrorCodeRequired)
	if err != nil {
		return CallError{}, err
	}

	errorDescription, err := decodeString(elements[3], ErrErrorDescriptionAbsent)
	if err != nil {
		return CallError{}, err
	}

	var details map[string]any
	if err := json.Unmarshal(elements[4], &details); err != nil {
		return CallError{}, fmt.Errorf("%w: %w", ErrErrorDetailsInvalid, err)
	}

	if details == nil {
		return CallError{}, ErrErrorDetailsInvalid
	}

	return CallError{
		UniqueID:         id,
		ErrorCode:        errorCode,
		ErrorDescription: errorDescription,
		ErrorDetails:     details,
	}, nil
}

func decodeMessageType(raw json.RawMessage) (MessageType, error) {
	var code uint8

	if err := json.Unmarshal(raw, &code); err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidFrame, err)
	}

	messageType := MessageType(code)
	switch messageType {
	case MessageTypeCall, MessageTypeCallResult, MessageTypeCallError:
		return messageType, nil
	default:
		return 0, fmt.Errorf("%w: %d", ErrUnsupportedFrameType, code)
	}
}

func decodeString(raw json.RawMessage, sentinel error) (string, error) {
	var value string

	if err := json.Unmarshal(raw, &value); err != nil {
		return "", fmt.Errorf("%w: %w", sentinel, err)
	}

	if value == "" {
		return "", sentinel
	}

	return value, nil
}

func validateMessageID(id string) error {
	if id == "" {
		return ErrInvalidMessageID
	}

	return nil
}

func validateAction(action string) error {
	if action == "" {
		return ErrInvalidAction
	}

	return nil
}

func marshalJSONArray(values ...any) ([]byte, error) {
	data, err := json.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFrame, err)
	}

	return data, nil
}

func DecodePayload[T any](raw json.RawMessage) (T, error) {
	var payload T

	if len(bytes.TrimSpace(raw)) == 0 {
		return payload, ErrPayloadRequired
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		return payload, fmt.Errorf("%w: %w", ErrPayloadDecode, err)
	}

	return payload, nil
}

func IsCall(frame Frame) bool {
	return frame != nil && frame.MessageType() == MessageTypeCall
}

func IsCallResult(frame Frame) bool {
	return frame != nil && frame.MessageType() == MessageTypeCallResult
}

func IsCallError(frame Frame) bool {
	return frame != nil && frame.MessageType() == MessageTypeCallError
}

func AsRawCall(frame Frame) (RawCall, error) {
	call, ok := frame.(RawCall)
	if !ok {
		return RawCall{}, errors.New("frame is not a raw call")
	}

	return call, nil
}
