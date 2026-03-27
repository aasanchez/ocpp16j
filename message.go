package ocpp16json

import (
	"encoding/json"
	"fmt"
)

// MessageType identifies the OCPP-J message type as defined in
// section 4.1.3, Table 2 of the OCPP-J 1.6 specification.
type MessageType uint8

const (
	// MessageTypeCall is MessageTypeId 2 (CALL).
	MessageTypeCall MessageType = 2
	// MessageTypeCallResult is MessageTypeId 3 (CALLRESULT).
	MessageTypeCallResult MessageType = 3
	// MessageTypeCallError is MessageTypeId 4 (CALLERROR).
	MessageTypeCallError MessageType = 4
)

const (
	callLength             = 4
	callResultLength       = 3
	callErrorLength        = 5
	callActionIndex        = 2
	callPayloadIndex       = 3
	callResultPayloadIndex = 2
	callErrorDetailsIndex  = 4
	emptyLength            = 0
	emptyString            = ""
	invalidTypeCode        = 0
	errorWrapFormat        = "%w: %w"
)

// Message is the common interface for all OCPP-J message types.
// Every Call, CallResult, and CallError satisfies this interface.
type Message interface {
	MessageType() MessageType
	MessageId() string
}

// IsCall reports whether the message is a Call (MessageTypeId 2).
func IsCall(message Message) bool {
	return message.MessageType() == MessageTypeCall
}

// IsCallResult reports whether the message is a CallResult
// (MessageTypeId 3).
func IsCallResult(message Message) bool {
	return message.MessageType() == MessageTypeCallResult
}

// IsCallError reports whether the message is a CallError
// (MessageTypeId 4).
func IsCallError(message Message) bool {
	return message.MessageType() == MessageTypeCallError
}

// AsCall extracts the Call from a Message, or returns
// errMessageNotCall if the message is not a Call.
func AsCall(message Message) (Call, error) {
	rawCall, isCall := message.(Call)
	if !isCall {
		return Call{}, errMessageNotCall
	}

	return rawCall, nil
}

// AsCallResult extracts the CallResult from a Message, or
// returns an error if the message is not a CallResult.
func AsCallResult(
	message Message,
) (CallResult, error) {
	callResult, isCallResult := message.(CallResult)
	if !isCallResult {
		return CallResult{}, errMessageNotCallResult
	}

	return callResult, nil
}

// AsCallError extracts the CallError from a Message, or
// returns an error if the message is not a CallError.
func AsCallError(
	message Message,
) (CallError, error) {
	callError, isCallError := message.(CallError)
	if !isCallError {
		return CallError{}, errMessageNotCallError
	}

	return callError, nil
}

func validateAction(action string) error {
	if action == emptyString {
		return ErrInvalidAction
	}

	return nil
}

func decodeString(
	raw json.RawMessage,
	sentinel error,
) (string, error) {
	var value string

	unmarshalErr := json.Unmarshal(raw, &value)
	if unmarshalErr != nil {
		return emptyString, fmt.Errorf(
			errorWrapFormat, sentinel, unmarshalErr,
		)
	}

	if value == emptyString {
		return emptyString, sentinel
	}

	return value, nil
}

func decodeMessageType(
	raw json.RawMessage,
) (MessageType, error) {
	var code uint8

	unmarshalErr := json.Unmarshal(raw, &code)
	if unmarshalErr != nil {
		return invalidTypeCode, fmt.Errorf(
			errorWrapFormat, ErrInvalidMessage, unmarshalErr,
		)
	}

	messageType := MessageType(code)

	switch messageType {
	case MessageTypeCall, MessageTypeCallResult, MessageTypeCallError:
		return messageType, nil
	default:
		return invalidTypeCode, ErrUnsupportedMessageType
	}
}

func marshalJSONArray(values ...any) ([]byte, error) {
	data, marshalErr := json.Marshal(values)
	if marshalErr != nil {
		return nil, fmt.Errorf(
			errorWrapFormat, ErrInvalidMessage, marshalErr,
		)
	}

	return data, nil
}
