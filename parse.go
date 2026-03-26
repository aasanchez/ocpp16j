package ocpp16json

import (
	"encoding/json"
	"fmt"
)

const (
	messageTypeIndex = 0
	uniqueIdIndex    = 1
	errorCodeIndex   = 2
	errorDescIndex   = 3
)

// Parse decodes raw bytes into a typed OCPP-J message. It
// validates the message wrapper (MessageTypeId, UniqueId,
// Action, Payload) and returns a RawCall, RawCallResult, or
// RawCallError. Payload contents are not decoded — they are
// preserved as json.RawMessage for later processing.
//
//nolint:ireturn // Concrete type depends on input.
func Parse(data []byte) (Message, error) {
	var elements []json.RawMessage

	unmarshalErr := json.Unmarshal(data, &elements)
	if unmarshalErr != nil {
		return nil, fmt.Errorf(
			errorWrapFormat, ErrInvalidMessage, unmarshalErr,
		)
	}

	if len(elements) == emptyLength {
		return nil, ErrInvalidMessage
	}

	messageType, typeErr := decodeMessageType(
		elements[messageTypeIndex],
	)
	if typeErr != nil {
		return nil, typeErr
	}

	switch messageType { //nolint:exhaustive // guarded by decodeMessageType.
	case Call:
		return parseCall(elements)
	case CallResult:
		return parseCallResult(elements)
	default:
		return parseCallError(elements)
	}
}

func parseCall(
	elements []json.RawMessage,
) (RawCall, error) {
	if len(elements) != callLength {
		return RawCall{}, ErrInvalidMessage
	}

	uniqueId, idErr := decodeUniqueId(
		elements[uniqueIdIndex],
	)
	if idErr != nil {
		return RawCall{}, idErr
	}

	action, actionErr := decodeString(
		elements[callActionIndex], ErrInvalidAction,
	)
	if actionErr != nil {
		return RawCall{}, actionErr
	}

	return RawCall{
		UniqueId: uniqueId,
		Action:   action,
		Payload:  elements[callPayloadIndex],
	}, nil
}

func parseCallResult(
	elements []json.RawMessage,
) (RawCallResult, error) {
	if len(elements) != callResultLength {
		return RawCallResult{}, ErrInvalidMessage
	}

	uniqueId, idErr := decodeUniqueId(
		elements[uniqueIdIndex],
	)
	if idErr != nil {
		return RawCallResult{}, idErr
	}

	return RawCallResult{
		UniqueId: uniqueId,
		Payload:  elements[callResultPayloadIndex],
	}, nil
}

func parseCallError(
	elements []json.RawMessage,
) (RawCallError, error) {
	if len(elements) != callErrorLength {
		return RawCallError{}, ErrInvalidMessage
	}

	uniqueId, idErr := decodeUniqueId(
		elements[uniqueIdIndex],
	)
	if idErr != nil {
		return RawCallError{}, idErr
	}

	errorCode, codeErr := decodeErrorCode(
		elements[errorCodeIndex],
	)
	if codeErr != nil {
		return RawCallError{}, codeErr
	}

	errorDescription, descErr := decodeErrorDescription(
		elements[errorDescIndex],
	)
	if descErr != nil {
		return RawCallError{}, descErr
	}

	errorDetails, detailsErr := decodeErrorDetails(
		elements[callErrorDetailsIndex],
	)
	if detailsErr != nil {
		return RawCallError{}, detailsErr
	}

	return RawCallError{
		UniqueId:         uniqueId,
		ErrorCode:        errorCode,
		ErrorDescription: errorDescription,
		ErrorDetails:     errorDetails,
	}, nil
}

func decodeUniqueId(
	raw json.RawMessage,
) (UniqueId, error) {
	value, err := decodeString(raw, ErrInvalidMessageID)
	if err != nil {
		return "", err
	}

	return NewUniqueId(value)
}

func decodeErrorCode(
	raw json.RawMessage,
) (ErrorCode, error) {
	value, err := decodeString(raw, ErrErrorCodeRequired)
	if err != nil {
		return "", err
	}

	return NewErrorCode(value)
}

func decodeErrorDescription(
	raw json.RawMessage,
) (string, error) {
	var value string

	unmarshalErr := json.Unmarshal(raw, &value)
	if unmarshalErr != nil {
		return emptyString, fmt.Errorf(
			errorWrapFormat,
			ErrErrorDescriptionAbsent,
			unmarshalErr,
		)
	}

	return value, nil
}

func decodeErrorDetails(
	raw json.RawMessage,
) (map[string]any, error) {
	var details map[string]any

	unmarshalErr := json.Unmarshal(raw, &details)
	if unmarshalErr != nil {
		return nil, fmt.Errorf(
			errorWrapFormat,
			ErrErrorDetailsInvalid,
			unmarshalErr,
		)
	}

	return details, nil
}
