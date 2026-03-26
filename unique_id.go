package ocpp16json

// maxUniqueIdLength is the maximum length of a UniqueId as
// defined in OCPP-J 1.6 specification section 4.1.4, Table 3:
// "Maximum of 36 characters, to allow for GUIDs."
const maxUniqueIdLength = 36

// UniqueId is a message identifier as defined in OCPP-J 1.6
// specification section 4.1.4, Table 3. It is a string with a
// maximum length of 36 characters. A UniqueId for a CALL
// message MUST be different from all message IDs previously
// used by the same sender on the same WebSocket connection.
// A UniqueId for a CALLRESULT or CALLERROR message MUST be
// equal to that of the CALL message it is responding to.
type UniqueId string

// NewUniqueId creates a validated UniqueId. It returns
// ErrInvalidMessageID if the value is empty or exceeds
// 36 characters (spec Table 3).
func NewUniqueId(value string) (UniqueId, error) {
	if value == emptyString {
		return "", ErrInvalidMessageID
	}

	if len(value) > maxUniqueIdLength {
		return "", ErrInvalidMessageID
	}

	return UniqueId(value), nil
}

// String returns the underlying string value.
func (uniqueId UniqueId) String() string {
	return string(uniqueId)
}
