package ocpp16json

import "errors"

var (
	ErrInvalidFrame            = errors.New("invalid OCPP-J frame")
	ErrUnsupportedFrameType    = errors.New("unsupported OCPP-J frame type")
	ErrInvalidMessageID        = errors.New("invalid OCPP-J message id")
	ErrInvalidAction           = errors.New("invalid OCPP-J action")
	ErrPayloadRequired         = errors.New("payload is required")
	ErrPayloadDecode           = errors.New("payload decode failed")
	ErrErrorCodeRequired       = errors.New("error code is required")
	ErrErrorDescriptionAbsent  = errors.New("error description is required")
	ErrErrorDetailsInvalid     = errors.New("error details must be a JSON object")
	ErrActionAlreadyRegistered = errors.New("action already registered")
	ErrUnknownAction           = errors.New("unknown action")
)
