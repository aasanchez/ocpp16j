package ocpp16json_test

import (
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
	"github.com/aasanchez/ocpp16messages/bootnotification"
)

// ExampleParse_dispatch demonstrates the typical workflow:
// receive a raw JSON message, detect its type, and act
// accordingly.
func ExampleParse_dispatch() {
	// A raw JSON message arrives from the WebSocket.
	wire := []byte(
		`[2, "19223201", "BootNotification",` +
			` {"chargePointVendor": "VendorX",` +
			` "chargePointModel": "Model1"}]`,
	)

	// Step 1: Parse — detects the message type.
	message, parseErr := ocpp16json.Parse(wire)
	if parseErr != nil {
		fmt.Println(parseErr)

		return
	}

	// Step 2: Detect — is it a Call, CallResult, or
	// CallError?
	switch {
	case ocpp16json.IsCall(message):
		rawCall, _ := ocpp16json.AsCall(message)

		fmt.Println("Received CALL")
		fmt.Println("Action:", rawCall.Action)

		// Step 3: Decode the Payload based on the Action.
		registry := ocpp16json.NewRegistry()
		_ = registry.Register(
			"BootNotification",
			ocpp16json.JSONDecoder(bootnotification.Req),
		)

		result, decodeErr := registry.Decode(
			rawCall.Action, rawCall.Payload,
		)
		if decodeErr != nil {
			fmt.Println("Decode failed:", decodeErr)

			return
		}

		req, isValid := result.(bootnotification.ReqMessage)
		if !isValid {
			fmt.Println("unexpected payload type")

			return
		}

		fmt.Println("Vendor:", req.ChargePointVendor)
		fmt.Println("Model:", req.ChargePointModel)

	case ocpp16json.IsCallResult(message):
		fmt.Println("Received CALLRESULT for",
			message.MessageId())

	case ocpp16json.IsCallError(message):
		fmt.Println("Received CALLERROR for",
			message.MessageId())

	default:
		fmt.Println("Unknown message type")
	}

	// Output:
	// Received CALL
	// Action: BootNotification
	// Vendor: VendorX
	// Model: Model1
}
