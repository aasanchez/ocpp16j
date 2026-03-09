package tests_fuzz

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	ocpp16json "github.com/aasanchez/ocpp16j"
	"github.com/aasanchez/ocpp16messages/authorize"
	"github.com/aasanchez/ocpp16messages/heartbeat"
	"github.com/aasanchez/ocpp16messages/types"
)

const maxFuzzStringLength = 1024

const (
	fuzzAuthorizeAction = "Authorize"
	fuzzHeartbeatAction = "Heartbeat"
)

func FuzzRegistryDecodeAuthorizeCall(f *testing.F) {
	f.Add("RFID-123")
	f.Add("")
	f.Add("12345678901234567890")
	f.Add("123456789012345678901")
	f.Add("contains\nnewline")
	f.Add("contains-tilde~")

	f.Fuzz(func(t *testing.T, idTag string) {
		if len(idTag) > maxFuzzStringLength {
			t.Skip()
		}

		registry, err := authorizeRegistry()
		if err != nil {
			t.Fatalf("RegisterRequest: %v", err)
		}

		frame, err := buildAuthorizeCallFrame(idTag)
		if err != nil {
			t.Fatalf("buildAuthorizeCallFrame: %v", err)
		}

		assertAuthorizeDecodeOutcome(t, registry, frame, idTag)
	})
}

func FuzzRegistryDecodeHeartbeatCallResult(f *testing.F) {
	f.Add("2025-01-02T15:04:05Z")
	f.Add("")
	f.Add("2025-01-02T15:04:05+01:00")
	f.Add("2025-01-02 15:04:05")
	f.Add("not-a-timestamp")

	f.Fuzz(func(t *testing.T, currentTime string) {
		if len(currentTime) > maxFuzzStringLength {
			t.Skip()
		}

		registry, err := heartbeatRegistry()
		if err != nil {
			t.Fatalf("RegisterConfirmation: %v", err)
		}

		frame, err := buildHeartbeatCallResultFrame(currentTime)
		if err != nil {
			t.Fatalf("buildHeartbeatCallResultFrame: %v", err)
		}

		assertHeartbeatDecodeOutcome(t, registry, frame)
	})
}

func authorizeRegistry() (*ocpp16json.Registry, error) {
	registry := ocpp16json.NewRegistry()

	err := registry.RegisterRequest(
		fuzzAuthorizeAction,
		ocpp16json.JSONDecoder(authorize.Req),
	)
	if err != nil {
		return nil, fmt.Errorf("RegisterRequest: %w", err)
	}

	return registry, nil
}

func heartbeatRegistry() (*ocpp16json.Registry, error) {
	registry := ocpp16json.NewRegistry()

	err := registry.RegisterConfirmation(
		fuzzHeartbeatAction,
		ocpp16json.JSONDecoder(heartbeat.Conf),
	)
	if err != nil {
		return nil, fmt.Errorf("RegisterConfirmation: %w", err)
	}

	return registry, nil
}

func buildAuthorizeCallFrame(idTag string) ([]byte, error) {
	payload, err := json.Marshal(map[string]string{"idTag": idTag})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(authorize payload): %w", err)
	}

	return append(
		[]byte(`[2,"uid-1","Authorize",`),
		append(payload, []byte(`]`)...)...,
	), nil
}

func buildHeartbeatCallResultFrame(currentTime string) ([]byte, error) {
	payload, err := json.Marshal(map[string]string{"currentTime": currentTime})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(heartbeat payload): %w", err)
	}

	return append(
		[]byte(`[3,"uid-1",`),
		append(payload, []byte(`]`)...)...,
	), nil
}

func assertAuthorizeDecodeOutcome(
	t *testing.T,
	registry *ocpp16json.Registry,
	frame []byte,
	idTag string,
) {
	t.Helper()

	decodedFrame, err := registry.DecodeCall(frame)
	if err != nil {
		assertPayloadDecodeValidationError(t, "DecodeCall", err)

		return
	}

	payloadValue, ok := decodedFrame.Payload.(authorize.ReqMessage)
	if !ok {
		t.Fatalf("unexpected payload type: %T", decodedFrame.Payload)
	}

	if payloadValue.IdTag.String() != idTag {
		t.Fatalf(
			"unexpected idTag after decode: got %q want %q",
			payloadValue.IdTag.String(),
			idTag,
		)
	}
}

func assertHeartbeatDecodeOutcome(
	t *testing.T,
	registry *ocpp16json.Registry,
	frame []byte,
) {
	t.Helper()

	decodedFrame, err := registry.DecodeCallResult(fuzzHeartbeatAction, frame)
	if err != nil {
		assertPayloadDecodeValidationError(t, "DecodeCallResult", err)

		return
	}

	payloadValue, ok := decodedFrame.Payload.(heartbeat.ConfMessage)
	if !ok {
		t.Fatalf("unexpected payload type: %T", decodedFrame.Payload)
	}

	assertHeartbeatTimeInvariant(t, payloadValue.CurrentTime.String())
}

func assertPayloadDecodeValidationError(
	t *testing.T,
	method string,
	err error,
) {
	t.Helper()

	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf("unexpected %s error: %v", method, err)
	}

	if !errors.Is(err, types.ErrEmptyValue) &&
		!errors.Is(err, types.ErrInvalidValue) {
		t.Fatalf("unexpected upstream validation error: %v", err)
	}
}

func assertHeartbeatTimeInvariant(t *testing.T, currentTimeString string) {
	t.Helper()

	if !strings.HasSuffix(currentTimeString, "Z") {
		t.Fatalf("decoded currentTime is not UTC: %q", currentTimeString)
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, currentTimeString)
	if err != nil {
		t.Fatalf("decoded currentTime is not RFC3339: %v", err)
	}

	if parsedTime.Location() != time.UTC {
		t.Fatalf(
			"decoded currentTime does not use UTC location: %q",
			currentTimeString,
		)
	}
}
