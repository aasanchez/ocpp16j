package ocpp16json_test

import (
	"errors"
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

// authReqInput mirrors an authorize request input.
type authReqInput struct {
	IdTag string `json:"idTag"`
}

// authReqOutput is the validated output.
type authReqOutput struct {
	IdTag string
}

// errIdTagRequired is a test-only sentinel error.
var errIdTagRequired = errors.New("idTag is required")

func authReqConstructor(
	input authReqInput,
) (authReqOutput, error) {
	if input.IdTag == "" {
		return authReqOutput{}, errIdTagRequired
	}

	return authReqOutput(input), nil
}

func ExampleRegistry() {
	// 1. Create a registry and register a decoder.
	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(
		authReqConstructor,
	)

	registerErr := registry.Register(
		"Authorize", decoder,
	)
	if registerErr != nil {
		fmt.Println(errPrefix, registerErr)

		return
	}

	// 2. Parse a raw CALL message from the wire.
	data := []byte(
		`[2, "abc123", "Authorize",` +
			` {"idTag": "RFID-001"}]`,
	)

	message, parseErr := ocpp16json.Parse(data)
	if parseErr != nil {
		fmt.Println(errPrefix, parseErr)

		return
	}

	// 3. Extract the Call and decode its payload.
	rawCall, castErr := ocpp16json.AsCall(message)
	if castErr != nil {
		fmt.Println(errPrefix, castErr)

		return
	}

	result, decodeErr := registry.Decode(
		rawCall.Action, rawCall.Payload,
	)
	if decodeErr != nil {
		fmt.Println(errPrefix, decodeErr)

		return
	}

	authorizeReq, isValid := result.(authReqOutput)
	if !isValid {
		fmt.Println("unexpected type")

		return
	}

	fmt.Println(authorizeReq.IdTag)

	// Output:
	// RFID-001
}
