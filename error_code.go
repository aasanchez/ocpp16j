package ocpp16json

// ErrorCode represents a valid CALLERROR error code as defined
// in OCPP-J 1.6 specification section 4.2.3, Table 7. Only the
// values listed in the specification are valid on the wire.
type ErrorCode string

// OCPP-J 1.6 error codes from Table 7. These are the only
// valid values for the errorCode field in a CALLERROR message.
const (
	// NotImplemented indicates the requested Action is not
	// known by the receiver.
	NotImplemented ErrorCode = "NotImplemented"
	// NotSupported indicates the requested Action is
	// recognized but not supported by the receiver.
	NotSupported ErrorCode = "NotSupported"
	// InternalError indicates an internal error occurred and
	// the receiver was not able to process the requested
	// Action successfully.
	InternalError ErrorCode = "InternalError"
	// ProtocolError indicates the payload for Action is
	// incomplete.
	ProtocolError ErrorCode = "ProtocolError"
	// SecurityError indicates a security issue occurred
	// preventing the receiver from completing the Action.
	SecurityError ErrorCode = "SecurityError"
	// FormationViolation indicates the payload for Action is
	// syntactically incorrect or does not conform the PDU
	// structure for Action.
	FormationViolation ErrorCode = "FormationViolation"
	// PropertyConstraintViolation indicates the payload is
	// syntactically correct but at least one field contains
	// an invalid value.
	PropertyConstraintViolation ErrorCode = "PropertyConstraintViolation"
	// OccurenceConstraintViolation indicates the payload is
	// syntactically correct but at least one of the fields
	// violates occurrence constraints.
	// NOTE: The spec spells this "Occurence" (missing the
	// second 'r'). We match the spec spelling exactly since
	// this string appears on the wire.
	//nolint:misspell // Intentional: matches OCPP-J 1.6 spec.
	OccurenceConstraintViolation ErrorCode = "OccurenceConstraintViolation"
	// TypeConstraintViolation indicates the payload is
	// syntactically correct but at least one of the fields
	// violates data type constraints (e.g. "somestring": 12).
	TypeConstraintViolation ErrorCode = "TypeConstraintViolation"
	// GenericError covers any other error not covered by the
	// previous error codes.
	GenericError ErrorCode = "GenericError"
)

// isValidErrorCode reports whether the given ErrorCode is one
// of the 10 valid values from spec Table 7.
func isValidErrorCode(code ErrorCode) bool {
	switch code {
	case NotImplemented,
		NotSupported,
		InternalError,
		ProtocolError,
		SecurityError,
		FormationViolation,
		PropertyConstraintViolation,
		OccurenceConstraintViolation,
		TypeConstraintViolation,
		GenericError:
		return true
	default:
		return false
	}
}

// NewErrorCode creates a validated ErrorCode. It returns
// ErrErrorCodeRequired if the value is empty or not one of
// the 10 valid codes defined in spec Table 7.
func NewErrorCode(value string) (ErrorCode, error) {
	if value == emptyString {
		return "", ErrErrorCodeRequired
	}

	code := ErrorCode(value)
	if !isValidErrorCode(code) {
		return "", ErrErrorCodeRequired
	}

	return code, nil
}

// String returns the wire-format string value.
func (errorCode ErrorCode) String() string {
	return string(errorCode)
}
