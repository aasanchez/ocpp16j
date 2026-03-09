package tests_fuzz

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	ocpp16json "github.com/aasanchez/ocpp16j"
	"github.com/aasanchez/ocpp16messages/authorize"
	"github.com/aasanchez/ocpp16messages/heartbeat"
	"github.com/aasanchez/ocpp16messages/types"
)

const maxFuzzStringLength = 1024

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

		registry := ocpp16json.NewRegistry()
		err := registry.RegisterRequest(
			"Authorize",
			ocpp16json.JSONDecoder(authorize.Req),
		)

		if err != nil {
			t.Fatalf("RegisterRequest: %v", err)
		}

		payload, err := json.Marshal(map[string]string{"idTag": idTag})
		if err != nil {
			t.Fatalf("json.Marshal(payload): %v", err)
		}

		frame := append(
			[]byte(`[2,"uid-1","Authorize",`),
			append(payload, []byte(`]`)...)...,
		)

		decodedFrame, err := registry.DecodeCall(frame)
		if err != nil {
			if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
				t.Fatalf("unexpected DecodeCall error: %v", err)
			}

			if !errors.Is(err, types.ErrEmptyValue) &&
				!errors.Is(err, types.ErrInvalidValue) {
				t.Fatalf("unexpected upstream validation error: %v", err)
			}

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

		registry := ocpp16json.NewRegistry()
		err := registry.RegisterConfirmation(
			"Heartbeat",
			ocpp16json.JSONDecoder(heartbeat.Conf),
		)
		if err != nil {
			t.Fatalf("RegisterConfirmation: %v", err)
		}

		payload, err := json.Marshal(
			map[string]string{"currentTime": currentTime},
		)
		if err != nil {
			t.Fatalf("json.Marshal(payload): %v", err)
		}

		frame := append(
			[]byte(`[3,"uid-1",`),
			append(payload, []byte(`]`)...)...,
		)

		decodedFrame, err := registry.DecodeCallResult("Heartbeat", frame)
		if err != nil {
			if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
				t.Fatalf("unexpected DecodeCallResult error: %v", err)
			}

			if !errors.Is(err, types.ErrEmptyValue) &&
				!errors.Is(err, types.ErrInvalidValue) {
				t.Fatalf("unexpected upstream validation error: %v", err)
			}

			return
		}

		payloadValue, ok := decodedFrame.Payload.(heartbeat.ConfMessage)
		if !ok {
			t.Fatalf("unexpected payload type: %T", decodedFrame.Payload)
		}

		currentTimeString := payloadValue.CurrentTime.String()
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
	})
}
