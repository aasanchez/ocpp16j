package ocpp16json_test

import (
	"encoding/json"
	"errors"
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
	"github.com/aasanchez/ocpp16messages/authorize"
	"github.com/aasanchez/ocpp16messages/heartbeat"
	"github.com/aasanchez/ocpp16messages/types"
)

type authorizeConfWirePayload struct {
	IdTagInfo authorizeConfWireInfo `json:"idTagInfo"`
}

type authorizeConfWireInfo struct {
	Status string `json:"status"`
}

// ExampleRegistry_authorizeFlow demonstrates a complete transport flow:
// decode an incoming CALL, validate it through ocpp16messages, apply handler
// logic, validate the confirmation payload, and wrap it back into CALLRESULT.
func ExampleRegistry_authorizeFlow() {
	registry := ocpp16json.NewRegistry()
	err := registry.RegisterRequest(
		"Authorize",
		ocpp16json.JSONDecoder(authorize.Req),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

	decoded, err := registry.DecodeCall(
		[]byte(`[2,"uid-1","Authorize",{"idTag":"RFID-123"}]`),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

	request, ok := decoded.Payload.(authorize.ReqMessage)
	if !ok {
		fmt.Println("unexpected payload type")

		return
	}

	status := authorizeStatusFor(request.IdTag.String())
	confInput := authorize.ConfInput{
		Status:      status,
		ExpiryDate:  nil,
		ParentIdTag: nil,
	}

	_, err = authorize.Conf(confInput)
	if err != nil {
		fmt.Println(err)

		return
	}

	payload, err := json.Marshal(authorizeConfWirePayload{
		IdTagInfo: authorizeConfWireInfo{
			Status: confInput.Status,
		},
	})
	if err != nil {
		fmt.Println(err)

		return
	}

	response, err := json.Marshal(ocpp16json.RawCallResult{
		UniqueID: decoded.UniqueID,
		Payload:  payload,
	})
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println("request action:", decoded.Action)
	fmt.Println("validated idTag:", request.IdTag.String())
	fmt.Println("response:", string(response))
	// Output:
	// request action: Authorize
	// validated idTag: RFID-123
	// response: [3,"uid-1",{"idTagInfo":{"status":"Accepted"}}]
}

// ExampleRegistry_authorizeValidationError demonstrates the error path for a
// CALL payload that is structurally valid JSON but fails OCPP field limits.
func ExampleRegistry_authorizeValidationError() {
	registry := ocpp16json.NewRegistry()
	err := registry.RegisterRequest(
		"Authorize",
		ocpp16json.JSONDecoder(authorize.Req),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

	_, err = registry.DecodeCall(
		[]byte(`[2,"uid-1","Authorize",{"idTag":"1234567890123456789012345"}]`),
	)
	if err == nil {
		fmt.Println("expected validation error")

		return
	}

	fmt.Println(errors.Is(err, ocpp16json.ErrPayloadDecode))
	fmt.Println(errors.Is(err, types.ErrInvalidValue))
	// Output:
	// true
	// true
}

func authorizeStatusFor(idTag string) string {
	if idTag == "RFID-BLOCKED" {
		return "Blocked"
	}

	return "Accepted"
}

// ExampleParse_rawCall demonstrates parsing a raw OCPP-J CALL frame before any
// typed payload decoding is applied.
func ExampleParse_rawCall() {
	frame, err := ocpp16json.Parse(
		[]byte(`[2,"uid-9","Authorize",{"idTag":"RFID-123"}]`),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

	call, err := ocpp16json.AsRawCall(frame)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(call.MessageType())
	fmt.Println(call.Action)
	fmt.Println(string(call.Payload))
	// Output:
	// 2
	// Authorize
	// {"idTag":"RFID-123"}
}

// ExampleRegistry_decodeCallResult demonstrates decoding a CALLRESULT payload.
// The action must be provided explicitly because OCPP-J does not include it in
// the response frame.
func ExampleRegistry_decodeCallResult() {
	registry := ocpp16json.NewRegistry()
	err := registry.RegisterConfirmation(
		"Heartbeat",
		ocpp16json.JSONDecoder(heartbeat.Conf),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

	decoded, err := registry.DecodeCallResult(
		"Heartbeat",
		[]byte(`[3,"uid-2",{"currentTime":"2025-01-02T15:04:05Z"}]`),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

	conf, ok := decoded.Payload.(heartbeat.ConfMessage)
	if !ok {
		fmt.Println("unexpected payload type")

		return
	}

	fmt.Println("action context:", decoded.Action)
	fmt.Println("current time:", conf.CurrentTime.String())
	// Output:
	// action context: Heartbeat
	// current time: 2025-01-02T15:04:05Z
}

// ExampleParse_callError demonstrates handling a raw CALLERROR frame.
func ExampleParse_callError() {
	frame, err := ocpp16json.Parse(
		[]byte(
			`[4,"uid-3","ProtocolError","invalid field",{"field":"idTag"}]`,
		),
	)
	if err != nil {
		fmt.Println(err)

		return
	}

	if !ocpp16json.IsCallError(frame) {
		fmt.Println("not a call error")

		return
	}

	callError := frame.(ocpp16json.CallError)
	fmt.Println(callError.ErrorCode)
	fmt.Println(callError.ErrorDescription)
	fmt.Println(callError.ErrorDetails["field"])
	// Output:
	// ProtocolError
	// invalid field
	// idTag
}
