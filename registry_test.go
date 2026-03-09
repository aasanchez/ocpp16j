package ocpp16json

import (
	"errors"
	"strings"
	"testing"

	"github.com/aasanchez/ocpp16messages/authorize"
	"github.com/aasanchez/ocpp16messages/heartbeat"
	"github.com/aasanchez/ocpp16messages/types"
)

const (
	testActionHeartbeat     = "Heartbeat"
	errRegisterRequestFmt   = "RegisterRequest: %v"
	errNeedPayloadDecodeFmt = "expected ErrPayloadDecode, got %v"
)

func TestRegistryDecodeCall(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	err := registry.RegisterRequest(
		testActionAuthorize,
		JSONDecoder(authorize.Req),
	)
	if err != nil {
		t.Fatalf(errRegisterRequestFmt, err)
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

func TestRegistryDecodeCallResult(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	err := registry.RegisterConfirmation(
		testActionHeartbeat,
		JSONDecoder(heartbeat.Conf),
	)
	if err != nil {
		t.Fatalf("RegisterConfirmation: %v", err)
	}

	frame, err := registry.DecodeCallResult(
		testActionHeartbeat,
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

	err := registry.RegisterRequest(
		testActionAuthorize,
		JSONDecoder(authorize.Req),
	)
	if err != nil {
		t.Fatalf(errRegisterRequestFmt, err)
	}

	err = registry.RegisterRequest(
		testActionAuthorize,
		JSONDecoder(authorize.Req),
	)
	if !errors.Is(err, ErrActionAlreadyRegistered) {
		t.Fatalf("expected ErrActionAlreadyRegistered, got %v", err)
	}
}

func TestRegistryRejectsUnknownAction(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	_, err := registry.DecodeCall(
		[]byte(`[2,"19223201","Authorize",{"idTag":"RFID-123"}]`),
	)
	if !errors.Is(err, ErrUnknownAction) {
		t.Fatalf("expected ErrUnknownAction, got %v", err)
	}
}

func TestRegistryValidationBranches(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	_, err := registry.DecodeCall([]byte(`[3,"uid",{}]`))
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf(
			"expected ErrInvalidFrame from DecodeCall on CALLRESULT, got %v",
			err,
		)
	}

	_, err = registry.DecodeCallResult("", []byte(`[3,"uid",{}]`))
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction, got %v", err)
	}

	_, err = registry.DecodeCallResult(
		testActionHeartbeat,
		[]byte(`[2,"uid","Heartbeat",{}]`),
	)
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf(
			"expected ErrInvalidFrame from DecodeCallResult on CALL, got %v",
			err,
		)
	}

	err = registry.RegisterRequest(testActionAuthorize, nil)
	if !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf(
			"expected ErrPayloadDecode for nil request decoder, got %v",
			err,
		)
	}

	err = registry.RegisterConfirmation(testActionHeartbeat, nil)
	if !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf(
			"expected ErrPayloadDecode for nil confirmation decoder, got %v",
			err,
		)
	}

	err = registry.RegisterRequest(
		"",
		JSONDecoder(func(input map[string]string) (map[string]string, error) {
			return input, nil
		}),
	)
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction for empty action, got %v", err)
	}
}

func TestRegistryDecodeCallPropagatesParseFailure(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	_, err := registry.DecodeCall([]byte(`{`))
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}
}

func TestRegistryDecodeCallResultPropagatesParseFailure(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	_, err := registry.DecodeCallResult(testActionHeartbeat, []byte(`{`))
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}
}

func TestRegistryDecodeCallResultPropagatesDecoderFailure(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	err := registry.RegisterConfirmation(
		testActionHeartbeat,
		JSONDecoder(heartbeat.Conf),
	)
	if err != nil {
		t.Fatalf("RegisterConfirmation: %v", err)
	}

	_, err = registry.DecodeCallResult(
		testActionHeartbeat,
		[]byte(`[3,"uid-3",{}]`),
	)
	if !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf(errNeedPayloadDecodeFmt, err)
	}

	if !errors.Is(err, types.ErrEmptyValue) {
		t.Fatalf("expected wrapped types.ErrEmptyValue, got %v", err)
	}
}

func TestRegistryDecodeCallRejectsLongAuthorizeIDTag(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	err := registry.RegisterRequest(
		testActionAuthorize,
		JSONDecoder(authorize.Req),
	)
	if err != nil {
		t.Fatalf(errRegisterRequestFmt, err)
	}

	_, err = registry.DecodeCall(
		[]byte(
			`[2,"19223201","Authorize",{"idTag":"1234567890123456789012345"}]`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf(errNeedPayloadDecodeFmt, err)
	}

	if !errors.Is(err, types.ErrInvalidValue) {
		t.Fatalf("expected wrapped types.ErrInvalidValue, got %v", err)
	}
}

func TestRegistryDecodeCallRejectsWrongAuthorizeIDTagType(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	err := registry.RegisterRequest(
		testActionAuthorize,
		JSONDecoder(authorize.Req),
	)
	if err != nil {
		t.Fatalf(errRegisterRequestFmt, err)
	}

	_, err = registry.DecodeCall(
		[]byte(`[2,"19223201","Authorize",{"idTag":123}]`),
	)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf(errNeedPayloadDecodeFmt, err)
	}

	if !strings.Contains(
		err.Error(),
		"cannot unmarshal number into Go struct field",
	) {
		t.Fatalf("expected JSON type error, got %v", err)
	}
}
