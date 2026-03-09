package ocpp16json

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	callFrameLength        = 4
	callResultFrameLength  = 3
	callErrorFrameLength   = 5
	callActionIndex        = 2
	callPayloadIndex       = 3
	callResultPayloadIndex = 2
	callErrorDetailsIndex  = 4
	emptyLength            = 0
	emptyString            = ""
	invalidTypeCode        = 0
	errorWrapFormat        = "%w: %w"
)

// MessageType identifies the OCPP-J frame kind.
type MessageType uint8

const (
	// MessageTypeCall identifies a CALL request frame.
	MessageTypeCall MessageType = 2
	// MessageTypeCallResult identifies a CALLRESULT response frame.
	MessageTypeCallResult MessageType = 3
	// MessageTypeCallError identifies a CALLERROR response frame.
	MessageTypeCallError MessageType = 4
)

// Frame is the common interface implemented by parsed OCPP-J frames.
//
// Both raw and decoded transport values satisfy this interface.
type Frame interface {
	MessageType() MessageType
	MessageID() string
}

// RawCall is a parsed OCPP-J CALL frame.
type RawCall struct {
	// UniqueID is the OCPP message correlation identifier.
	UniqueID string
	// Action is the OCPP action name carried by the request frame.
	Action string
	// Payload is the raw JSON payload object.
	Payload json.RawMessage
}

// RawCallResult is a parsed OCPP-J CALLRESULT frame.
type RawCallResult struct {
	// UniqueID is the OCPP message correlation identifier.
	UniqueID string
	// Payload is the raw JSON payload object.
	Payload json.RawMessage
}

// CallError is a parsed OCPP-J CALLERROR frame.
type CallError struct {
	// UniqueID is the OCPP message correlation identifier.
	UniqueID string
	// ErrorCode is the OCPP error code string.
	ErrorCode string
	// ErrorDescription is the human-readable error text.
	ErrorDescription string
	// ErrorDetails is the structured error details object.
	ErrorDetails map[string]any
}

// DecodedCall is a CALL frame with a decoded payload.
type DecodedCall struct {
	// UniqueID is the OCPP message correlation identifier.
	UniqueID string
	// Action is the OCPP action name carried by the request frame.
	Action string
	// Payload is the typed value returned by a registered decoder.
	Payload any
}

// DecodedCallResult is a CALLRESULT frame with a decoded payload.
type DecodedCallResult struct {
	// UniqueID is the OCPP message correlation identifier.
	UniqueID string
	// Action is the caller-provided action context for the response frame.
	Action string
	// Payload is the typed value returned by a registered decoder.
	Payload any
}

// MessageType returns the OCPP-J frame type.
func (rawCall RawCall) MessageType() MessageType {
	_ = rawCall

	return MessageTypeCall
}

// MessageID returns the unique message ID.
func (rawCall RawCall) MessageID() string {
	return rawCall.UniqueID
}

// MessageType returns the OCPP-J frame type.
func (rawCallResult RawCallResult) MessageType() MessageType {
	_ = rawCallResult

	return MessageTypeCallResult
}

// MessageID returns the unique message ID.
func (rawCallResult RawCallResult) MessageID() string {
	return rawCallResult.UniqueID
}

// MessageType returns the OCPP-J frame type.
func (callError CallError) MessageType() MessageType {
	_ = callError

	return MessageTypeCallError
}

// MessageID returns the unique message ID.
func (callError CallError) MessageID() string {
	return callError.UniqueID
}

// MessageType returns the OCPP-J frame type.
func (decodedCall DecodedCall) MessageType() MessageType {
	_ = decodedCall

	return MessageTypeCall
}

// MessageID returns the unique message ID.
func (decodedCall DecodedCall) MessageID() string {
	return decodedCall.UniqueID
}

// MessageType returns the OCPP-J frame type.
func (decodedCallResult DecodedCallResult) MessageType() MessageType {
	_ = decodedCallResult

	return MessageTypeCallResult
}

// MessageID returns the unique message ID.
func (decodedCallResult DecodedCallResult) MessageID() string {
	return decodedCallResult.UniqueID
}

// Parse validates a raw OCPP-J frame and returns a typed raw frame.
//
// Depending on the incoming message type, the returned value is RawCall,
// RawCallResult, or CallError.
func Parse(data []byte) (Frame, error) {
	var elements []json.RawMessage

	err := json.Unmarshal(data, &elements)
	if err != nil {
		return nil, fmt.Errorf(errorWrapFormat, ErrInvalidFrame, err)
	}

	if len(elements) == emptyLength {
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
	case MessageTypeCallError:
		return parseCallError(elements)
	}

	return nil, fmt.Errorf("%w: %d", ErrUnsupportedFrameType, messageType)
}

// MarshalJSON renders the CALL frame as an OCPP-J array.
func (rawCall RawCall) MarshalJSON() ([]byte, error) {
	err := validateMessageID(rawCall.UniqueID)
	if err != nil {
		return nil, err
	}

	err = validateAction(rawCall.Action)
	if err != nil {
		return nil, err
	}

	if len(bytes.TrimSpace(rawCall.Payload)) == emptyLength {
		return nil, ErrPayloadRequired
	}

	return marshalJSONArray(
		MessageTypeCall,
		rawCall.UniqueID,
		rawCall.Action,
		rawCall.Payload,
	)
}

// MarshalJSON renders the CALLRESULT frame as an OCPP-J array.
func (rawCallResult RawCallResult) MarshalJSON() ([]byte, error) {
	err := validateMessageID(rawCallResult.UniqueID)
	if err != nil {
		return nil, err
	}

	if len(bytes.TrimSpace(rawCallResult.Payload)) == emptyLength {
		return nil, ErrPayloadRequired
	}

	return marshalJSONArray(
		MessageTypeCallResult,
		rawCallResult.UniqueID,
		rawCallResult.Payload,
	)
}

// MarshalJSON renders the CALLERROR frame as an OCPP-J array.
func (callError CallError) MarshalJSON() ([]byte, error) {
	err := validateMessageID(callError.UniqueID)
	if err != nil {
		return nil, err
	}

	if callError.ErrorCode == emptyString {
		return nil, ErrErrorCodeRequired
	}

	if callError.ErrorDescription == emptyString {
		return nil, ErrErrorDescriptionAbsent
	}

	errorDetails := callError.ErrorDetails
	if errorDetails == nil {
		errorDetails = map[string]any{}
	}

	return marshalJSONArray(
		MessageTypeCallError,
		callError.UniqueID,
		callError.ErrorCode,
		callError.ErrorDescription,
		errorDetails,
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

	uniqueID, err := decodeString(elements[1], ErrInvalidMessageID)
	if err != nil {
		return RawCall{}, err
	}

	action, err := decodeString(elements[callActionIndex], ErrInvalidAction)
	if err != nil {
		return RawCall{}, err
	}

	if len(bytes.TrimSpace(elements[callPayloadIndex])) == emptyLength {
		return RawCall{}, ErrPayloadRequired
	}

	return RawCall{
		UniqueID: uniqueID,
		Action:   action,
		Payload:  elements[callPayloadIndex],
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

	uniqueID, err := decodeString(elements[1], ErrInvalidMessageID)
	if err != nil {
		return RawCallResult{}, err
	}

	if len(bytes.TrimSpace(elements[callResultPayloadIndex])) == emptyLength {
		return RawCallResult{}, ErrPayloadRequired
	}

	return RawCallResult{
		UniqueID: uniqueID,
		Payload:  elements[callResultPayloadIndex],
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

	uniqueID, err := decodeString(elements[1], ErrInvalidMessageID)
	if err != nil {
		return CallError{}, err
	}

	errorCode, err := decodeString(elements[2], ErrErrorCodeRequired)
	if err != nil {
		return CallError{}, err
	}

	errorDescription, err := decodeString(
		elements[3],
		ErrErrorDescriptionAbsent,
	)
	if err != nil {
		return CallError{}, err
	}

	var details map[string]any

	err = json.Unmarshal(elements[callErrorDetailsIndex], &details)
	if err != nil {
		return CallError{}, fmt.Errorf(
			errorWrapFormat,
			ErrErrorDetailsInvalid,
			err,
		)
	}

	if details == nil {
		return CallError{}, ErrErrorDetailsInvalid
	}

	return CallError{
		UniqueID:         uniqueID,
		ErrorCode:        errorCode,
		ErrorDescription: errorDescription,
		ErrorDetails:     details,
	}, nil
}

func decodeMessageType(raw json.RawMessage) (MessageType, error) {
	var code uint8

	err := json.Unmarshal(raw, &code)
	if err != nil {
		return invalidTypeCode, fmt.Errorf(
			errorWrapFormat,
			ErrInvalidFrame,
			err,
		)
	}

	messageType := MessageType(code)
	switch messageType {
	case MessageTypeCall, MessageTypeCallResult, MessageTypeCallError:
		return messageType, nil
	default:
		return invalidTypeCode, fmt.Errorf(
			"%w: %d",
			ErrUnsupportedFrameType,
			code,
		)
	}
}

func decodeString(raw json.RawMessage, sentinel error) (string, error) {
	var value string

	err := json.Unmarshal(raw, &value)
	if err != nil {
		return emptyString, fmt.Errorf(errorWrapFormat, sentinel, err)
	}

	if value == emptyString {
		return emptyString, sentinel
	}

	return value, nil
}

func validateMessageID(id string) error {
	if id == emptyString {
		return ErrInvalidMessageID
	}

	return nil
}

func validateAction(action string) error {
	if action == emptyString {
		return ErrInvalidAction
	}

	return nil
}

func marshalJSONArray(values ...any) ([]byte, error) {
	data, err := json.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf(errorWrapFormat, ErrInvalidFrame, err)
	}

	return data, nil
}

// DecodePayload unmarshals a raw payload into the requested Go type.
//
// It performs JSON decoding only. Domain-specific validation should still be
// handled by higher-level constructors when required.
func DecodePayload[T any](raw json.RawMessage) (T, error) {
	var payload T

	if len(bytes.TrimSpace(raw)) == emptyLength {
		return payload, ErrPayloadRequired
	}

	err := json.Unmarshal(raw, &payload)
	if err != nil {
		return payload, fmt.Errorf(errorWrapFormat, ErrPayloadDecode, err)
	}

	return payload, nil
}

// IsCall reports whether frame is a CALL.
func IsCall(frame Frame) bool {
	return frame != nil && frame.MessageType() == MessageTypeCall
}

// IsCallResult reports whether frame is a CALLRESULT.
func IsCallResult(frame Frame) bool {
	return frame != nil && frame.MessageType() == MessageTypeCallResult
}

// IsCallError reports whether frame is a CALLERROR.
func IsCallError(frame Frame) bool {
	return frame != nil && frame.MessageType() == MessageTypeCallError
}

// AsRawCall extracts a RawCall from frame.
//
// It returns an error when frame is not a RawCall value.
func AsRawCall(frame Frame) (RawCall, error) {
	call, ok := frame.(RawCall)
	if !ok {
		return RawCall{}, errFrameNotRawCall
	}

	return call, nil
}
