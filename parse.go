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
// Action, Payload) and returns a Call, CallResult, or
// CallError. Payload contents are not decoded — they are
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
	case MessageTypeCall:
		return parseCall(elements)
	case MessageTypeCallResult:
		return parseCallResult(elements)
	default:
		return parseCallError(elements)
	}
}

func parseCall(
	elements []json.RawMessage,
) (Call, error) {
	if len(elements) != callLength {
		return Call{}, ErrInvalidMessage
	}

	uniqueId, idErr := decodeUniqueId(
		elements[uniqueIdIndex],
	)
	if idErr != nil {
		return Call{}, idErr
	}

	action, actionErr := decodeString(
		elements[callActionIndex], ErrInvalidAction,
	)
	if actionErr != nil {
		return Call{}, actionErr
	}

	return Call{
		UniqueId: uniqueId,
		Action:   action,
		Payload:  elements[callPayloadIndex],
	}, nil
}

func parseCallResult(
	elements []json.RawMessage,
) (CallResult, error) {
	if len(elements) != callResultLength {
		return CallResult{}, ErrInvalidMessage
	}

	uniqueId, idErr := decodeUniqueId(
		elements[uniqueIdIndex],
	)
	if idErr != nil {
		return CallResult{}, idErr
	}

	return CallResult{
		UniqueId: uniqueId,
		Payload:  elements[callResultPayloadIndex],
	}, nil
}

func parseCallError(
	elements []json.RawMessage,
) (CallError, error) {
	if len(elements) != callErrorLength {
		return CallError{}, ErrInvalidMessage
	}

	uniqueId, idErr := decodeUniqueId(
		elements[uniqueIdIndex],
	)
	if idErr != nil {
		return CallError{}, idErr
	}

	errorCode, codeErr := decodeErrorCode(
		elements[errorCodeIndex],
	)
	if codeErr != nil {
		return CallError{}, codeErr
	}

	errorDescription, descErr := decodeErrorDescription(
		elements[errorDescIndex],
	)
	if descErr != nil {
		return CallError{}, descErr
	}

	errorDetails, detailsErr := decodeErrorDetails(
		elements[callErrorDetailsIndex],
	)
	if detailsErr != nil {
		return CallError{}, detailsErr
	}

	return CallError{
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
