package ocpp16json

// DecodedMessage represents a CALL or CALLRESULT message
// whose Payload has been decoded into a typed value T by a
// registered [PayloadDecoder]. Unlike [Call] and
// [CallResult] which carry the Payload as raw JSON, a
// DecodedMessage holds the validated domain object
// produced by an ocpp16messages constructor.
type DecodedMessage[T any] struct {
	messageType MessageType
	UniqueId    UniqueId
	Action      string
	Payload     T
}

// NewDecodedCall creates a DecodedMessage with
// MessageTypeId 2 (CALL). It validates that Action is
// non-empty. The UniqueId must already be validated via
// [NewUniqueId].
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
		messageType: MessageTypeCall,
		UniqueId:    uniqueId,
		Action:      action,
		Payload:     payload,
	}, nil
}

// NewDecodedCallResult creates a DecodedMessage with
// MessageTypeId 3 (CALLRESULT). It validates that Action
// is non-empty. The Action must be provided explicitly
// because CALLRESULT does not carry it on the wire.
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
		messageType: MessageTypeCallResult,
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
