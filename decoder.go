package ocpp16json

import (
	"encoding/json"
	"fmt"
)

// PayloadDecoder decodes a raw JSON payload into a validated
// domain type. Implementations are registered with a Registry
// to handle specific OCPP Actions.
type PayloadDecoder func(json.RawMessage) (any, error)

// JSONDecoder creates a PayloadDecoder that unmarshals a
// json.RawMessage into Input, then passes it to a constructor
// function that returns a validated Output. This bridges the
// raw JSON wire format to ocpp16messages constructors such as
// authorize.Req or bootnotification.Conf.
func JSONDecoder[Input any, Output any](
	constructor func(Input) (Output, error),
) PayloadDecoder {
	return func(raw json.RawMessage) (any, error) {
		var input Input

		unmarshalErr := json.Unmarshal(raw, &input)
		if unmarshalErr != nil {
			return nil, fmt.Errorf(
				errorWrapFormat,
				ErrPayloadDecode,
				unmarshalErr,
			)
		}

		output, constructErr := constructor(input)
		if constructErr != nil {
			return nil, fmt.Errorf(
				errorWrapFormat,
				ErrPayloadDecode,
				constructErr,
			)
		}

		return output, nil
	}
}
