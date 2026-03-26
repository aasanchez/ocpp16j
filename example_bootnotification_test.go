package ocpp16json_test

import (
	"encoding/json"
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
	"github.com/aasanchez/ocpp16messages/bootnotification"
)

const heartbeatIntervalSeconds = 300

// ExampleParse_bootNotificationRequest demonstrates parsing
// a BootNotification.req CALL from the wire, then decoding
// the Payload into a validated BootNotification request
// using ocpp16messages.
func ExampleParse_bootNotificationRequest() {
	wire := []byte(
		`[2, "19223201", "BootNotification",` +
			` {"chargePointVendor": "VendorX",` +
			` "chargePointModel": "SingleSocket"}]`,
	)

	message, parseErr := ocpp16json.Parse(wire)
	if parseErr != nil {
		fmt.Println(parseErr)

		return
	}

	rawCall, _ := ocpp16json.AsCall(message)

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(
		bootnotification.Req,
	)
	_ = registry.Register("BootNotification", decoder)

	result, decodeErr := registry.Decode(
		rawCall.Action, rawCall.Payload,
	)
	if decodeErr != nil {
		fmt.Println(decodeErr)

		return
	}

	bootReq, isValid := result.(bootnotification.ReqMessage)
	if !isValid {
		fmt.Println("unexpected type")

		return
	}

	fmt.Println(bootReq.ChargePointVendor)
	fmt.Println(bootReq.ChargePointModel)

	// Output:
	// VendorX
	// SingleSocket
}

// ExampleParse_bootNotificationResponse demonstrates how a
// Central System builds a BootNotification.conf CALLRESULT
// and serializes it back to the wire.
func ExampleParse_bootNotificationResponse() {
	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	payload := map[string]any{
		"status":      "Accepted",
		"currentTime": "2013-02-01T20:53:32Z",
		"interval":    heartbeatIntervalSeconds,
	}

	callResult, _ := ocpp16json.NewCallResult(
		uniqueId, payload,
	)

	wireBytes, marshalErr := json.Marshal(callResult)
	if marshalErr != nil {
		fmt.Println(marshalErr)

		return
	}

	// Re-parse to verify round-trip correctness.
	roundTrip, _ := ocpp16json.Parse(wireBytes)
	fmt.Println(ocpp16json.IsCallResult(roundTrip))
	fmt.Println(roundTrip.MessageId())

	// Output:
	// true
	// 19223201
}

// ExampleParse_bootNotificationInvalidPayload shows what
// happens when the chargePointVendor exceeds 20 chars.
func ExampleParse_bootNotificationInvalidPayload() {
	wire := []byte(
		`[2, "19223201", "BootNotification",` +
			` {"chargePointVendor":` +
			` "ThisVendorNameIsWayTooLong",` +
			` "chargePointModel": "Model"}]`,
	)

	message, _ := ocpp16json.Parse(wire)
	rawCall, _ := ocpp16json.AsCall(message)

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(
		bootnotification.Req,
	)
	_ = registry.Register("BootNotification", decoder)

	_, decodeErr := registry.Decode(
		rawCall.Action, rawCall.Payload,
	)

	// The error wraps ErrPayloadDecode — the vendor name
	// exceeds the 20-character limit from the spec.
	fmt.Println(decodeErr != nil)

	// Output:
	// true
}
