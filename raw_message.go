package ocpp16json

import "encoding/json"

// RawCall represents a parsed CALL message (MessageTypeId 2)
// with an undecoded Payload.
type RawCall struct {
	UniqueId UniqueId
	Action   string
	Payload  json.RawMessage
}

// RawCallResult represents a parsed CALLRESULT message
// (MessageTypeId 3) with an undecoded Payload.
type RawCallResult struct {
	UniqueId UniqueId
	Payload  json.RawMessage
}

// RawCallError represents a parsed CALLERROR message
// (MessageTypeId 4) as defined in section 4.2.3.
type RawCallError struct {
	UniqueId         UniqueId
	ErrorCode        ErrorCode
	ErrorDescription string
	ErrorDetails     map[string]any
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

// MessageType returns the CallError MessageTypeId (4).
func (RawCallError) MessageType() MessageType {
	return CallError
}

// MessageId returns the UniqueId correlation identifier.
func (rawCallError RawCallError) MessageId() string {
	return rawCallError.UniqueId.String()
}

// MarshalJSON serializes a CALLERROR to its canonical OCPP-J
// array:
// [4, "<UniqueId>", "<ErrorCode>", "<ErrorDescription>",
//
//	{<ErrorDetails>}]
func (rawCallError RawCallError) MarshalJSON() (
	[]byte, error,
) {
	return marshalJSONArray(
		CallError,
		rawCallError.UniqueId.String(),
		rawCallError.ErrorCode.String(),
		rawCallError.ErrorDescription,
		rawCallError.ErrorDetails,
	)
}
