package ocpp16json_test

import (
	"encoding/json"
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
	"github.com/aasanchez/ocpp16messages/authorize"
	"github.com/aasanchez/ocpp16messages/heartbeat"
	"github.com/aasanchez/ocpp16messages/types"
)

func TestRegistryDecodeCallAndConfirmation(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	if err := registry.RegisterRequest("Authorize", ocpp16json.JSONDecoder(authorize.Req)); err != nil {
		t.Fatalf("RegisterRequest: %v", err)
	}

	if err := registry.RegisterConfirmation("Heartbeat", ocpp16json.JSONDecoder(heartbeat.Conf)); err != nil {
		t.Fatalf("RegisterConfirmation: %v", err)
	}

	call, err := registry.DecodeCall([]byte(`[2,"uid-1","Authorize",{"idTag":"RFID-123"}]`))
	if err != nil {
		t.Fatalf("DecodeCall: %v", err)
	}

	if !ocpp16json.IsCall(call) {
		t.Fatal("expected decoded call to satisfy IsCall")
	}

	req, ok := call.Payload.(authorize.ReqMessage)
	if !ok {
		t.Fatalf("unexpected request payload type: %T", call.Payload)
	}

	if req.IdTag.String() != "RFID-123" {
		t.Fatalf("unexpected idTag: %q", req.IdTag.String())
	}

	result, err := registry.DecodeCallResult("Heartbeat", []byte(`[3,"uid-1",{"currentTime":"2025-01-02T15:04:05Z"}]`))
	if err != nil {
		t.Fatalf("DecodeCallResult: %v", err)
	}

	if !ocpp16json.IsCallResult(result) {
		t.Fatal("expected decoded result to satisfy IsCallResult")
	}

	conf, ok := result.Payload.(heartbeat.ConfMessage)
	if !ok {
		t.Fatalf("unexpected confirmation payload type: %T", result.Payload)
	}

	if conf.CurrentTime.String() != "2025-01-02T15:04:05Z" {
		t.Fatalf("unexpected currentTime: %q", conf.CurrentTime.String())
	}
}

func TestParseAndDecodeValidationErrors(t *testing.T) {
	t.Parallel()

	frame, err := ocpp16json.Parse([]byte(`[2,"uid-2","Authorize",{"idTag":"1234567890123456789012345"}]`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	rawCall, err := ocpp16json.AsRawCall(frame)
	if err != nil {
		t.Fatalf("AsRawCall: %v", err)
	}

	if rawCall.Action != "Authorize" {
		t.Fatalf("unexpected action: %q", rawCall.Action)
	}

	registry := ocpp16json.NewRegistry()
	if err := registry.RegisterRequest("Authorize", ocpp16json.JSONDecoder(authorize.Req)); err != nil {
		t.Fatalf("RegisterRequest: %v", err)
	}

	_, err = registry.DecodeCall([]byte(`[2,"uid-2","Authorize",{"idTag":"1234567890123456789012345"}]`))
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf("expected ErrPayloadDecode, got %v", err)
	}

	if !errors.Is(err, types.ErrInvalidValue) {
		t.Fatalf("expected wrapped types.ErrInvalidValue, got %v", err)
	}

	_, err = ocpp16json.DecodePayload[map[string]string](json.RawMessage(`1`))
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf("expected ErrPayloadDecode, got %v", err)
	}
}

func TestRegistryExposesParseAndConstructorFailures(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	if err := registry.RegisterConfirmation("Heartbeat", ocpp16json.JSONDecoder(heartbeat.Conf)); err != nil {
		t.Fatalf("RegisterConfirmation: %v", err)
	}

	_, err := registry.DecodeCall([]byte(`{`))
	if !errors.Is(err, ocpp16json.ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame from DecodeCall parse failure, got %v", err)
	}

	_, err = registry.DecodeCallResult("Heartbeat", []byte(`{`))
	if !errors.Is(err, ocpp16json.ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame from DecodeCallResult parse failure, got %v", err)
	}

	_, err = registry.DecodeCallResult("Heartbeat", []byte(`[3,"uid-3",{}]`))
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf("expected ErrPayloadDecode from confirmation decode, got %v", err)
	}

	if !errors.Is(err, types.ErrEmptyValue) {
		t.Fatalf("expected wrapped types.ErrEmptyValue, got %v", err)
	}
}
