package ocpp16json

import (
	"encoding/json"
	"fmt"
)

// JSONDecoder adapts a validating constructor into a payload decoder.
//
// The returned decoder unmarshals the raw JSON payload into Input and then
// passes it to constructor. This is the normal way to connect transport-level
// decoding in this package with domain validation from ocpp16messages.
func JSONDecoder[Input any, Output any](
	constructor func(Input) (Output, error),
) PayloadDecoder {
	return func(raw json.RawMessage) (any, error) {
		var input Input

		err := json.Unmarshal(raw, &input)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrPayloadDecode, err)
		}

		return constructor(input)
	}
}
