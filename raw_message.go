package ocpp16json

import "encoding/json"

// RawCall represents a parsed CALL message (MessageTypeId 2)
// with an undecoded Payload.
type RawCall struct {
	UniqueId string
	Action   string
	Payload  json.RawMessage
}

// RawCallResult represents a parsed CALLRESULT message
// (MessageTypeId 3) with an undecoded Payload.
type RawCallResult struct {
	UniqueId string
	Payload  json.RawMessage
}

// RawCallError represents a parsed CALLERROR message
// (MessageTypeId 4) as defined in section 4.2.3.
type RawCallError struct {
	UniqueId         string
	ErrorCode        string
	ErrorDescription string
	ErrorDetails     map[string]any
}

// MessageType returns the Call MessageTypeId (2).
func (RawCall) MessageType() MessageType {
	return Call
}

// MessageId returns the UniqueId correlation identifier.
func (rawCall RawCall) MessageId() string {
	return rawCall.UniqueId
}

// MessageType returns the CallResult MessageTypeId (3).
func (RawCallResult) MessageType() MessageType {
	return CallResult
}

// MessageId returns the UniqueId correlation identifier.
func (rawCallResult RawCallResult) MessageId() string {
	return rawCallResult.UniqueId
}

// MessageType returns the CallError MessageTypeId (4).
func (RawCallError) MessageType() MessageType {
	return CallError
}

// MessageId returns the UniqueId correlation identifier.
func (rawCallError RawCallError) MessageId() string {
	return rawCallError.UniqueId
}
