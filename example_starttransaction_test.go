package ocpp16json_test

import (
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
	"github.com/aasanchez/ocpp16messages/starttransaction"
)

// ExampleParse_startTransactionRequest demonstrates
// parsing a StartTransaction.req CALL and decoding it
// into a validated message with ocpp16messages.
func ExampleParse_startTransactionRequest() {
	wire := []byte(
		`[2, "tx-001", "StartTransaction",` +
			` {"connectorId": 1,` +
			` "idTag": "RFID-ABC123",` +
			` "meterStart": 0,` +
			` "timestamp":` +
			` "2024-01-15T08:00:00Z"}]`,
	)

	message, _ := ocpp16json.Parse(wire)
	rawCall, _ := ocpp16json.AsCall(message)

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(
		starttransaction.Req,
	)
	_ = registry.Register("StartTransaction", decoder)

	result, decodeErr := registry.Decode(
		rawCall.Action, rawCall.Payload,
	)
	if decodeErr != nil {
		fmt.Println(decodeErr)

		return
	}

	startReq, isValid := result.(starttransaction.ReqMessage)
	if !isValid {
		fmt.Println("unexpected type")

		return
	}

	fmt.Println(startReq.IdTag)
	fmt.Println(startReq.ConnectorId)

	// Output:
	// RFID-ABC123
	// 1
}

// ExampleParse_startTransactionWrongType demonstrates
// what happens when connectorId is a string instead of
// an integer — the JSON unmarshal fails.
func ExampleParse_startTransactionWrongType() {
	wire := []byte(
		`[2, "tx-002", "StartTransaction",` +
			` {"connectorId": "one",` +
			` "idTag": "RFID-ABC123",` +
			` "meterStart": 0,` +
			` "timestamp":` +
			` "2024-01-15T08:00:00Z"}]`,
	)

	message, _ := ocpp16json.Parse(wire)
	rawCall, _ := ocpp16json.AsCall(message)

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(
		starttransaction.Req,
	)
	_ = registry.Register("StartTransaction", decoder)

	_, decodeErr := registry.Decode(
		rawCall.Action, rawCall.Payload,
	)

	// The error wraps ErrPayloadDecode — connectorId must
	// be an integer, not a string.
	fmt.Println(decodeErr != nil)

	// Output:
	// true
}
