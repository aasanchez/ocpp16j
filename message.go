package ocpp16json

import (
	"encoding/json"
	"fmt"
)

// MessageType identifies the OCPP-J message type as defined in
// section 4.1.3, Table 2 of the OCPP-J 1.6 specification.
type MessageType uint8

const (
	// Call is MessageTypeId 2 (Client-to-Server request).
	Call MessageType = 2
	// CallResult is MessageTypeId 3 (Server-to-Client response).
	CallResult MessageType = 3
	// CallError is MessageTypeId 4 (Server-to-Client error).
	CallError MessageType = 4
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
	return message.MessageType() == Call
}

// IsCallResult reports whether the message is a CallResult
// (MessageTypeId 3).
func IsCallResult(message Message) bool {
	return message.MessageType() == CallResult
}

// IsCallError reports whether the message is a CallError
// (MessageTypeId 4).
func IsCallError(message Message) bool {
	return message.MessageType() == CallError
}

// AsRawCall extracts the RawCall from a Message, or returns
// errMessageNotCall if the message is not a RawCall.
func AsRawCall(message Message) (RawCall, error) {
	rawCall, isRawCall := message.(RawCall)
	if !isRawCall {
		return RawCall{}, errMessageNotCall
	}

	return rawCall, nil
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
	case Call, CallResult, CallError:
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
