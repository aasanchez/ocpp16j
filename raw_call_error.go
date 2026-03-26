package ocpp16json

// RawCallError represents a parsed CALLERROR message
// (MessageTypeId 4) as defined in section 4.2.3.
type RawCallError struct {
	UniqueId         UniqueId
	ErrorCode        ErrorCode
	ErrorDescription string
	ErrorDetails     map[string]any
}

// NewRawCallError creates a CALLERROR message
// (MessageTypeId 4). ErrorDescription may be empty per the
// spec. If ErrorDetails is nil, it defaults to an empty
// object.
func NewRawCallError(
	uniqueId UniqueId,
	errorCode ErrorCode,
	errorDescription string,
	errorDetails map[string]any,
) (RawCallError, error) {
	details := errorDetails
	if details == nil {
		details = map[string]any{}
	}

	return RawCallError{
		UniqueId:         uniqueId,
		ErrorCode:        errorCode,
		ErrorDescription: errorDescription,
		ErrorDetails:     details,
	}, nil
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
