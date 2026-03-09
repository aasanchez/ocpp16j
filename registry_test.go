package ocpp16json

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/aasanchez/ocpp16messages/authorize"
	"github.com/aasanchez/ocpp16messages/heartbeat"
	"github.com/aasanchez/ocpp16messages/types"
)

func TestRegistryDecodeCall(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	if err := registry.RegisterRequest("Authorize", JSONDecoder(authorize.Req)); err != nil {
		t.Fatalf("RegisterRequest: %v", err)
	}

	frame, err := registry.DecodeCall(
		[]byte(`[2,"19223201","Authorize",{"idTag":"RFID-123"}]`),
	)
	if err != nil {
		t.Fatalf("DecodeCall: %v", err)
	}

	payload, ok := frame.Payload.(authorize.ReqMessage)
	if !ok {
		t.Fatalf("unexpected payload type: %T", frame.Payload)
	}

	if payload.IdTag.String() != "RFID-123" {
		t.Fatalf("unexpected idTag: %q", payload.IdTag.String())
	}
}

func TestJSONDecoderRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	decoder := JSONDecoder(authorize.Req)
	_, err := decoder(json.RawMessage(`{"idTag":1}`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRegistryDecodeCallResult(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	if err := registry.RegisterConfirmation("Heartbeat", JSONDecoder(heartbeat.Conf)); err != nil {
		t.Fatalf("RegisterConfirmation: %v", err)
	}

	frame, err := registry.DecodeCallResult(
		"Heartbeat",
		[]byte(`[3,"19223201",{"currentTime":"2025-01-02T15:04:05Z"}]`),
	)
	if err != nil {
		t.Fatalf("DecodeCallResult: %v", err)
	}

	payload, ok := frame.Payload.(heartbeat.ConfMessage)
	if !ok {
		t.Fatalf("unexpected payload type: %T", frame.Payload)
	}

	if payload.CurrentTime.String() != "2025-01-02T15:04:05Z" {
		t.Fatalf("unexpected currentTime: %q", payload.CurrentTime.String())
	}
}

func TestRegistryRejectsDuplicateAction(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	if err := registry.RegisterRequest("Authorize", JSONDecoder(authorize.Req)); err != nil {
		t.Fatalf("RegisterRequest: %v", err)
	}

	err := registry.RegisterRequest("Authorize", JSONDecoder(authorize.Req))
	if !errors.Is(err, ErrActionAlreadyRegistered) {
		t.Fatalf("expected ErrActionAlreadyRegistered, got %v", err)
	}
}

func TestRegistryRejectsUnknownAction(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	_, err := registry.DecodeCall([]byte(`[2,"19223201","Authorize",{"idTag":"RFID-123"}]`))
	if !errors.Is(err, ErrUnknownAction) {
		t.Fatalf("expected ErrUnknownAction, got %v", err)
	}
}

func TestRegistryDecodeCallRejectsAuthorizeIDTagLongerThanCiString20(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	if err := registry.RegisterRequest("Authorize", JSONDecoder(authorize.Req)); err != nil {
		t.Fatalf("RegisterRequest: %v", err)
	}

	_, err := registry.DecodeCall(
		[]byte(`[2,"19223201","Authorize",{"idTag":"1234567890123456789012345"}]`),
	)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf("expected ErrPayloadDecode, got %v", err)
	}

	if !errors.Is(err, types.ErrInvalidValue) {
		t.Fatalf("expected wrapped types.ErrInvalidValue, got %v", err)
	}
}

func TestRegistryDecodeCallRejectsAuthorizeIDTagWithWrongJSONType(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	if err := registry.RegisterRequest("Authorize", JSONDecoder(authorize.Req)); err != nil {
		t.Fatalf("RegisterRequest: %v", err)
	}

	_, err := registry.DecodeCall(
		[]byte(`[2,"19223201","Authorize",{"idTag":123}]`),
	)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf("expected ErrPayloadDecode, got %v", err)
	}

	if !strings.Contains(err.Error(), "cannot unmarshal number into Go struct field") {
		t.Fatalf("expected JSON type error, got %v", err)
	}
}
