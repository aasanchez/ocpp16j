package ocpp16json

import "errors"

var (
	// ErrInvalidMessage reports an invalid OCPP-J message envelope.
	ErrInvalidMessage = errors.New("invalid OCPP-J message")
	// ErrUnsupportedMessageType reports an unknown MessageTypeNumber.
	ErrUnsupportedMessageType = errors.New(
		"unsupported OCPP-J message type",
	)
	// ErrInvalidMessageID reports a missing or invalid unique ID.
	ErrInvalidMessageID = errors.New("invalid OCPP-J message id")
	// ErrInvalidAction reports a missing or invalid action name.
	ErrInvalidAction = errors.New("invalid OCPP-J action")
	// ErrPayloadRequired reports a missing payload.
	ErrPayloadRequired = errors.New("payload is required")
	// ErrPayloadDecode reports a payload decoding or validation failure.
	ErrPayloadDecode = errors.New("payload decode failed")
	// ErrErrorCodeRequired reports a missing CALLERROR code.
	ErrErrorCodeRequired = errors.New("error code is required")
	// ErrErrorDescriptionAbsent reports a missing CALLERROR description.
	ErrErrorDescriptionAbsent = errors.New("error description is required")
	// ErrErrorDetailsInvalid reports invalid CALLERROR details.
	ErrErrorDetailsInvalid = errors.New(
		"error details must be a JSON object",
	)
	// ErrActionAlreadyRegistered reports duplicate action registration.
	ErrActionAlreadyRegistered = errors.New("action already registered")
	// ErrUnknownAction reports an action with no registered decoder.
	ErrUnknownAction  = errors.New("unknown action")
	errMessageNotCall = errors.New("message is not a Call")
)
