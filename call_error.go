package ocpp16json

// CallError represents a parsed CALLERROR message
// (MessageTypeId 4) as defined in section 4.2.3.
type CallError struct {
	UniqueId         UniqueId
	ErrorCode        ErrorCode
	ErrorDescription string
	ErrorDetails     map[string]any
}

// NewCallError creates a CALLERROR message
// (MessageTypeId 4). ErrorDescription may be empty per the
// spec. If ErrorDetails is nil, it defaults to an empty
// object.
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
