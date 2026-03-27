package ocpp16json

// CallError represents a CALLERROR message
// (MessageTypeId 4) as defined in OCPP-J 1.6
// specification section 4.2.3. A CallError consists of
// 5 elements: MessageTypeId, UniqueId, ErrorCode,
// ErrorDescription, and ErrorDetails. It is sent when a
// CALL cannot be handled — either due to a transport
// error or because the CALL content does not meet the
// requirements for a proper message.
type CallError struct {
	UniqueId         UniqueId
	ErrorCode        ErrorCode
	ErrorDescription string
	ErrorDetails     map[string]any
}

// NewCallError creates a CALLERROR message
// (MessageTypeId 4). The UniqueId MUST match the UniqueId
// of the [Call] being responded to. ErrorDescription
// should be filled in if possible; an empty string is
// valid per the spec (section 4.2.3, Table 6). If
// ErrorDetails is nil, it defaults to an empty JSON
// object as required by the spec.
func NewCallError(
	uniqueId UniqueId,
	errorCode ErrorCode,
	errorDescription string,
	errorDetails map[string]any,
) (CallError, error) {
	details := errorDetails
	if details == nil {
		details = map[string]any{}
	}

	return CallError{
		UniqueId:         uniqueId,
		ErrorCode:        errorCode,
		ErrorDescription: errorDescription,
		ErrorDetails:     details,
	}, nil
}

// MessageType returns the CallError MessageTypeId (4).
func (CallError) MessageType() MessageType {
	return MessageTypeCallError
}

// MessageId returns the UniqueId correlation identifier.
func (rawCallError CallError) MessageId() string {
	return rawCallError.UniqueId.String()
}

// MarshalJSON serializes a CALLERROR to its canonical OCPP-J
// array:
// [4, "<UniqueId>", "<ErrorCode>", "<ErrorDescription>",
//
//	{<ErrorDetails>}]
func (rawCallError CallError) MarshalJSON() (
	[]byte, error,
) {
	return marshalJSONArray(
		MessageTypeCallError,
		rawCallError.UniqueId.String(),
		rawCallError.ErrorCode.String(),
		rawCallError.ErrorDescription,
		rawCallError.ErrorDetails,
	)
}
