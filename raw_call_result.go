package ocpp16json

import "encoding/json"

// RawCallResult represents a parsed CALLRESULT message
// (MessageTypeId 3) with an undecoded Payload.
type RawCallResult struct {
	UniqueId UniqueId
	Payload  json.RawMessage
}

// NewRawCallResult creates a CALLRESULT message
// (MessageTypeId 3). It marshals the payload to
// json.RawMessage.
func NewRawCallResult(
	uniqueId UniqueId,
	payload any,
) (RawCallResult, error) {
	rawPayload, marshalErr := marshalPayload(payload)
	if marshalErr != nil {
		return RawCallResult{}, marshalErr
	}

	return RawCallResult{
		UniqueId: uniqueId,
		Payload:  rawPayload,
	}, nil
}

// MessageType returns the CallResult MessageTypeId (3).
func (RawCallResult) MessageType() MessageType {
	return CallResult
}

// MessageId returns the UniqueId correlation identifier.
func (rawCallResult RawCallResult) MessageId() string {
	return rawCallResult.UniqueId.String()
}

// MarshalJSON serializes a CALLRESULT to its canonical
// OCPP-J array: [3, "<UniqueId>", {<Payload>}]
func (rawCallResult RawCallResult) MarshalJSON() (
	[]byte, error,
) {
	return marshalJSONArray(
		CallResult,
		rawCallResult.UniqueId.String(),
		rawCallResult.Payload,
	)
}
