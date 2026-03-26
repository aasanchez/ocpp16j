package ocpp16json

// DecodedMessage represents a CALL or CALLRESULT message whose
// Payload has been decoded into a typed value by a registered
// PayloadDecoder.
type DecodedMessage[T any] struct {
	messageType MessageType
	UniqueId    UniqueId
	Action      string
	Payload     T
}

// NewDecodedCall creates a DecodedMessage with MessageTypeId 2
// (Call). It validates that UniqueId and Action are non-empty.
func NewDecodedCall[T any](
	uniqueId UniqueId,
	action string,
	payload T,
) (DecodedMessage[T], error) {
	validationErr := validateAction(action)
	if validationErr != nil {
		return DecodedMessage[T]{}, validationErr
	}

	return DecodedMessage[T]{
		messageType: Call,
		UniqueId:    uniqueId,
		Action:      action,
		Payload:     payload,
	}, nil
}

// NewDecodedCallResult creates a DecodedMessage with
// MessageTypeId 3 (CallResult). It validates that UniqueId and
// Action are non-empty.
func NewDecodedCallResult[T any](
	uniqueId UniqueId,
	action string,
	payload T,
) (DecodedMessage[T], error) {
	validationErr := validateAction(action)
	if validationErr != nil {
		return DecodedMessage[T]{}, validationErr
	}

	return DecodedMessage[T]{
		messageType: CallResult,
		UniqueId:    uniqueId,
		Action:      action,
		Payload:     payload,
	}, nil
}

// MessageType returns the MessageTypeId for this decoded message.
func (decoded DecodedMessage[T]) MessageType() MessageType {
	return decoded.messageType
}

// MessageId returns the UniqueId correlation identifier.
func (decoded DecodedMessage[T]) MessageId() string {
	return decoded.UniqueId.String()
}
