package ocpp16json

import (
	"encoding/json"
	"fmt"
)

func marshalPayload(payload any) (json.RawMessage, error) {
	data, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return nil, fmt.Errorf(
			errorWrapFormat,
			ErrPayloadDecode,
			marshalErr,
		)
	}

	return data, nil
}
