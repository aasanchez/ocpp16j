package ocpp16json

import (
	"encoding/json"
	"fmt"
)

func JSONDecoder[Input any, Output any](
	constructor func(Input) (Output, error),
) PayloadDecoder {
	return func(raw json.RawMessage) (any, error) {
		var input Input

		if err := json.Unmarshal(raw, &input); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrPayloadDecode, err)
		}

		return constructor(input)
	}
}
