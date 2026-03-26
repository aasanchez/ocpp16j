package ocpp16json

import "encoding/json"

// RawCall represents a parsed CALL message (MessageTypeId 2)
// with an undecoded Payload.
type RawCall struct {
	UniqueId UniqueId
	Action   string
	Payload  json.RawMessage
}

// NewRawCall creates a CALL message (MessageTypeId 2). It
// validates Action and marshals the payload to
// json.RawMessage.
func NewRawCall(
	uniqueId UniqueId,
	action string,
	payload any,
) (RawCall, error) {
	validationErr := validateAction(action)
	if validationErr != nil {
		return RawCall{}, validationErr
	}

	rawPayload, marshalErr := marshalPayload(payload)
	if marshalErr != nil {
		return RawCall{}, marshalErr
	}

	return RawCall{
		UniqueId: uniqueId,
		Action:   action,
		Payload:  rawPayload,
	}, nil
}

// MessageType returns the Call MessageTypeId (2).
func (RawCall) MessageType() MessageType {
	return Call
}

// MessageId returns the UniqueId correlation identifier.
func (rawCall RawCall) MessageId() string {
	return rawCall.UniqueId.String()
}

// MarshalJSON serializes a CALL to its canonical OCPP-J array:
// [2, "<UniqueId>", "<Action>", {<Payload>}]
func (rawCall RawCall) MarshalJSON() ([]byte, error) {
	return marshalJSONArray(
		Call,
		rawCall.UniqueId.String(),
		rawCall.Action,
		rawCall.Payload,
	)
}
