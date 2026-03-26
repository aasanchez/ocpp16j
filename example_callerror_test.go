package ocpp16json_test

import (
	"encoding/json"
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

// ExampleParse_unknownActionCallError demonstrates how a
// Central System responds with a CALLERROR when it receives
// a CALL with an unrecognized Action.
func ExampleParse_unknownActionCallError() {
	wire := []byte(
		`[2, "req-99", "CustomAction", {}]`,
	)

	message, _ := ocpp16json.Parse(wire)
	rawCall, _ := ocpp16json.AsRawCall(message)

	registry := ocpp16json.NewRegistry()

	_, decodeErr := registry.Decode(
		rawCall.Action, rawCall.Payload,
	)

	callError, _ := ocpp16json.NewRawCallError(
		rawCall.UniqueId,
		ocpp16json.NotImplemented,
		decodeErr.Error(),
		map[string]any{},
	)

	wireBytes, marshalErr := json.Marshal(callError)
	if marshalErr != nil {
		fmt.Println(marshalErr)

		return
	}

	fmt.Println(string(wireBytes))

	// Output:
	// [4,"req-99","NotImplemented","unknown action",{}]
}

// ExampleParse_callErrorFromWire demonstrates parsing a
// CALLERROR received from the other party.
func ExampleParse_callErrorFromWire() {
	wire := []byte(
		`[4, "req-99", "FormationViolation",` +
			` "Missing required field", {}]`,
	)

	message, _ := ocpp16json.Parse(wire)

	fmt.Println(ocpp16json.IsCallError(message))
	fmt.Println(message.MessageId())

	// Output:
	// true
	// req-99
}
