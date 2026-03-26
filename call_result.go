package ocpp16json

import "encoding/json"

// CallResult represents a parsed CALLRESULT message
// (MessageTypeId 3) with an undecoded Payload.
type CallResult struct {
	UniqueId UniqueId
	Payload  json.RawMessage
}

// NewCallResult creates a CALLRESULT message
// (MessageTypeId 3). It marshals the payload to
// json.RawMessage.
func NewCallResult(
	uniqueId UniqueId,
	payload any,
) (CallResult, error) {
	rawPayload, marshalErr := marshalPayload(payload)
	if marshalErr != nil {
		return CallResult{}, marshalErr
	}

	return CallResult{
		UniqueId: uniqueId,
		Payload:  rawPayload,
	}, nil
}

// MessageType returns the CallResult MessageTypeId (3).
func (CallResult) MessageType() MessageType {
	return MessageTypeCallResult
}

// MessageId returns the UniqueId correlation identifier.
func (rawCallResult CallResult) MessageId() string {
	return rawCallResult.UniqueId.String()
}

// MarshalJSON serializes a CALLRESULT to its canonical
// OCPP-J array: [3, "<UniqueId>", {<Payload>}]
func (rawCallResult CallResult) MarshalJSON() (
	[]byte, error,
) {
	return marshalJSONArray(
		MessageTypeCallResult,
		rawCallResult.UniqueId.String(),
		rawCallResult.Payload,
	)
}
