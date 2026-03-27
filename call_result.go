package ocpp16json

import "encoding/json"

// CallResult represents a CALLRESULT message
// (MessageTypeId 3) as defined in OCPP-J 1.6
// specification section 4.2.2. A CallResult consists of
// 3 elements: MessageTypeId, UniqueId, and Payload.
// Note that the Action is not present on the wire — the
// caller must track which Action the UniqueId corresponds
// to in order to decode the Payload.
type CallResult struct {
	UniqueId UniqueId
	Payload  json.RawMessage
}

// NewCallResult creates a CALLRESULT message
// (MessageTypeId 3). It marshals the payload to
// [json.RawMessage]. The UniqueId MUST match the
// UniqueId of the original [Call] being responded to.
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
