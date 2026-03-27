package ocpp16json

import "encoding/json"

// Call represents a CALL message (MessageTypeId 2) as
// defined in OCPP-J 1.6 specification section 4.2.1. A
// Call always consists of 4 elements: MessageTypeId,
// UniqueId, Action, and Payload. The Payload is preserved
// as [json.RawMessage] for later decoding via [Registry].
type Call struct {
	UniqueId UniqueId
	Action   string
	Payload  json.RawMessage
}

// NewCall creates a CALL message (MessageTypeId 2). It
// validates that Action is non-empty and marshals the
// payload to [json.RawMessage]. The UniqueId must be
// created beforehand with [NewUniqueId].
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
