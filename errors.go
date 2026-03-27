package ocpp16json

import "errors"

// Sentinel errors for every failure mode in the OCPP-J
// message wrapper. Check with [errors.Is].
var (
	// ErrInvalidMessage indicates the input is not a valid
	// OCPP-J message envelope — either not valid JSON, not
	// a JSON array, or has the wrong number of elements for
	// the detected MessageTypeId.
	ErrInvalidMessage = errors.New("invalid OCPP-J message")
	// ErrUnsupportedMessageType indicates the
	// MessageTypeNumber is not 2, 3, or 4 (spec section
	// 4.1.3, Table 2). A server SHALL ignore such messages.
	ErrUnsupportedMessageType = errors.New(
		"unsupported OCPP-J message type",
	)
	// ErrInvalidMessageID indicates the UniqueId is missing,
	// empty, or exceeds 36 characters (spec section 4.1.4,
	// Table 3).
	ErrInvalidMessageID = errors.New(
		"invalid OCPP-J message id",
	)
	// ErrInvalidAction indicates the Action field in a CALL
	// is missing or empty (spec section 4.2.1, Table 4).
	ErrInvalidAction = errors.New("invalid OCPP-J action")
	// ErrPayloadRequired indicates the Payload element is
	// missing from the message array.
	ErrPayloadRequired = errors.New("payload is required")
	// ErrPayloadDecode indicates the Payload could not be
	// decoded or failed validation by an ocpp16messages
	// constructor registered in the [Registry].
	ErrPayloadDecode = errors.New("payload decode failed")
	// ErrErrorCodeRequired indicates the ErrorCode in a
	// CALLERROR is missing, empty, or not one of the 10
	// valid values from spec Table 7.
	ErrErrorCodeRequired = errors.New(
		"error code is required",
	)
	// ErrErrorDescriptionAbsent indicates the
	// ErrorDescription in a CALLERROR could not be decoded
	// as a JSON string.
	ErrErrorDescriptionAbsent = errors.New(
		"error description is required",
	)
	// ErrErrorDetailsInvalid indicates the ErrorDetails in
	// a CALLERROR is not a valid JSON object. The spec
	// requires ErrorDetails to be a JSON object; if there
	// are no details, it MUST be an empty object {}.
	ErrErrorDetailsInvalid = errors.New(
		"error details must be a JSON object",
	)
	// ErrActionAlreadyRegistered indicates a [PayloadDecoder]
	// has already been registered for the given Action name
	// in a [Registry].
	ErrActionAlreadyRegistered = errors.New(
		"action already registered",
	)
	// ErrUnknownAction indicates no [PayloadDecoder] is
	// registered for the given Action name in the
	// [Registry].
	ErrUnknownAction        = errors.New("unknown action")
	errMessageNotCall       = errors.New("message is not a Call")
	errMessageNotCallResult = errors.New(
		"message is not a CallResult",
	)
	errMessageNotCallError = errors.New(
		"message is not a CallError",
	)
)
