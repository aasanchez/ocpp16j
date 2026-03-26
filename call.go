package ocpp16json

import "encoding/json"

// Call represents a parsed CALL message (MessageTypeId 2)
// with an undecoded Payload.
type Call struct {
	UniqueId UniqueId
	Action   string
	Payload  json.RawMessage
}

// NewCall creates a CALL message (MessageTypeId 2). It
// validates Action and marshals the payload to
// json.RawMessage.
func NewCall(
	uniqueId UniqueId,
	action string,
	payload any,
) (Call, error) {
	validationErr := validateAction(action)
	if validationErr != nil {
		return Call{}, validationErr
	}

	rawPayload, marshalErr := marshalPayload(payload)
	if marshalErr != nil {
		return Call{}, marshalErr
	}

	return Call{
		UniqueId: uniqueId,
		Action:   action,
		Payload:  rawPayload,
	}, nil
}

// MessageType returns the Call MessageTypeId (2).
func (Call) MessageType() MessageType {
	return MessageTypeCall
}

// MessageId returns the UniqueId correlation identifier.
func (rawCall Call) MessageId() string {
	return rawCall.UniqueId.String()
}

// MarshalJSON serializes a CALL to its canonical OCPP-J array:
// [2, "<UniqueId>", "<Action>", {<Payload>}]
func (rawCall Call) MarshalJSON() ([]byte, error) {
	return marshalJSONArray(
		MessageTypeCall,
		rawCall.UniqueId.String(),
		rawCall.Action,
		rawCall.Payload,
	)
}
