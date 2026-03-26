package ocpp16json_test

import (
	"encoding/json"
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
	"github.com/aasanchez/ocpp16messages/heartbeat"
)

// ExampleParse_heartbeatRequest demonstrates parsing a
// Heartbeat.req CALL. The Heartbeat request has an empty
// payload — the spec allows both null and {}.
func ExampleParse_heartbeatRequest() {
	wire := []byte(
		`[2, "42", "Heartbeat", {}]`,
	)

	message, parseErr := ocpp16json.Parse(wire)
	if parseErr != nil {
		fmt.Println(parseErr)

		return
	}

	rawCall, _ := ocpp16json.AsRawCall(message)

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(heartbeat.Req)
	_ = registry.Register("Heartbeat", decoder)

	_, decodeErr := registry.Decode(
		rawCall.Action, rawCall.Payload,
	)
	if decodeErr != nil {
		fmt.Println(decodeErr)

		return
	}

	fmt.Println(rawCall.Action)
	fmt.Println(rawCall.MessageId())

	// Output:
	// Heartbeat
	// 42
}

// ExampleParse_heartbeatResponse demonstrates building a
// Heartbeat.conf CALLRESULT with the Central System's
// current time.
func ExampleParse_heartbeatResponse() {
	uniqueId, _ := ocpp16json.NewUniqueId("42")

	payload := map[string]string{
		"currentTime": "2024-01-15T10:30:00Z",
	}

	callResult, _ := ocpp16json.NewRawCallResult(
		uniqueId, payload,
	)

	wireBytes, marshalErr := json.Marshal(callResult)
	if marshalErr != nil {
		fmt.Println(marshalErr)

		return
	}

	fmt.Println(string(wireBytes))

	// Output:
	// [3,"42",{"currentTime":"2024-01-15T10:30:00Z"}]
}
