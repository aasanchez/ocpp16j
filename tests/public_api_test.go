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

const (
	publicActionAuthorize      = "Authorize"
	publicActionHeartbeat      = "Heartbeat"
	publicNeedPayloadDecodeFmt = "expected ErrPayloadDecode, got %v"
)

func TestRegistryDecodeCall(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	registerAuthorizeDecoder(t, registry)

	call, err := registry.DecodeCall(
		[]byte(`[2,"uid-1","Authorize",{"idTag":"RFID-123"}]`),
	)
	if err != nil {
		t.Fatalf("DecodeCall: %v", err)
	}

	if !ocpp16json.IsCall(call) {
		t.Fatal("expected decoded call to satisfy IsCall")
	}

	request, requestOK := call.Payload.(authorize.ReqMessage)
	if !requestOK {
		t.Fatalf("unexpected request payload type: %T", call.Payload)
	}

	if request.IdTag.String() != "RFID-123" {
		t.Fatalf("unexpected idTag: %q", request.IdTag.String())
	}
}

func TestRegistryDecodeCallResult(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	registerHeartbeatDecoder(t, registry)

	result, err := registry.DecodeCallResult(
		publicActionHeartbeat,
		[]byte(`[3,"uid-1",{"currentTime":"2025-01-02T15:04:05Z"}]`),
	)
	if err != nil {
		t.Fatalf("DecodeCallResult: %v", err)
	}

	if !ocpp16json.IsCallResult(result) {
		t.Fatal("expected decoded result to satisfy IsCallResult")
	}

	confirmation, confirmationOK := result.Payload.(heartbeat.ConfMessage)
	if !confirmationOK {
		t.Fatalf("unexpected payload type: %T", result.Payload)
	}

	if confirmation.CurrentTime.String() != "2025-01-02T15:04:05Z" {
		t.Fatalf(
			"unexpected currentTime: %q",
			confirmation.CurrentTime.String(),
		)
	}
}

func TestParseAndDecodeValidationErrors(t *testing.T) {
	t.Parallel()

	frame, err := ocpp16json.Parse(
		[]byte(`[2,"uid-2","Authorize",{"idTag":"1234567890123456789012345"}]`),
	)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	rawCall, err := ocpp16json.AsRawCall(frame)
	if err != nil {
		t.Fatalf("AsRawCall: %v", err)
	}

	if rawCall.Action != publicActionAuthorize {
		t.Fatalf("unexpected action: %q", rawCall.Action)
	}

	registry := ocpp16json.NewRegistry()
	registerAuthorizeDecoder(t, registry)

	_, err = registry.DecodeCall(
		[]byte(`[2,"uid-2","Authorize",{"idTag":"1234567890123456789012345"}]`),
	)
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf(publicNeedPayloadDecodeFmt, err)
	}

	if !errors.Is(err, types.ErrInvalidValue) {
		t.Fatalf("expected wrapped types.ErrInvalidValue, got %v", err)
	}

	_, err = ocpp16json.DecodePayload[map[string]string](json.RawMessage(`1`))
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf(publicNeedPayloadDecodeFmt, err)
	}
}

func TestRegistryExposesParseFailures(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	registerHeartbeatDecoder(t, registry)

	_, err := registry.DecodeCall([]byte(`{`))
	if !errors.Is(err, ocpp16json.ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}

	_, err = registry.DecodeCallResult(publicActionHeartbeat, []byte(`{`))
	if !errors.Is(err, ocpp16json.ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}
}

func TestRegistryExposesConstructorFailures(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	registerHeartbeatDecoder(t, registry)

	_, err := registry.DecodeCallResult(
		publicActionHeartbeat,
		[]byte(`[3,"uid-3",{}]`),
	)
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf("expected ErrPayloadDecode, got %v", err)
	}

	if !errors.Is(err, types.ErrEmptyValue) {
		t.Fatalf("expected wrapped types.ErrEmptyValue, got %v", err)
	}
}

func registerAuthorizeDecoder(t *testing.T, registry *ocpp16json.Registry) {
	t.Helper()

	err := registry.RegisterRequest(
		publicActionAuthorize,
		ocpp16json.JSONDecoder(authorize.Req),
	)
	if err != nil {
		t.Fatalf("RegisterRequest: %v", err)
	}
}

func registerHeartbeatDecoder(t *testing.T, registry *ocpp16json.Registry) {
	t.Helper()

	err := registry.RegisterConfirmation(
		publicActionHeartbeat,
		ocpp16json.JSONDecoder(heartbeat.Conf),
	)
	if err != nil {
		t.Fatalf("RegisterConfirmation: %v", err)
	}
}
