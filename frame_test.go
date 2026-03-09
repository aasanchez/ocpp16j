package ocpp16json

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestParseRawCall(t *testing.T) {
	t.Parallel()

	frame, err := Parse([]byte(`[2,"19223201","Authorize",{"idTag":"RFID-123"}]`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	call, ok := frame.(RawCall)
	if !ok {
		t.Fatalf("unexpected frame type: %T", frame)
	}

	if call.UniqueID != "19223201" {
		t.Fatalf("unexpected id: %q", call.UniqueID)
	}

	if call.Action != "Authorize" {
		t.Fatalf("unexpected action: %q", call.Action)
	}

	if string(call.Payload) != `{"idTag":"RFID-123"}` {
		t.Fatalf("unexpected payload: %s", string(call.Payload))
	}
}

func TestParseCallResult(t *testing.T) {
	t.Parallel()

	frame, err := Parse([]byte(`[3,"19223201",{"currentTime":"2025-01-02T15:04:05Z"}]`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	result, ok := frame.(RawCallResult)
	if !ok {
		t.Fatalf("unexpected frame type: %T", frame)
	}

	if result.UniqueID != "19223201" {
		t.Fatalf("unexpected id: %q", result.UniqueID)
	}
}

func TestParseCallError(t *testing.T) {
	t.Parallel()

	frame, err := Parse([]byte(`[4,"19223201","ProtocolError","bad payload",{"field":"idTag"}]`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	callError, ok := frame.(CallError)
	if !ok {
		t.Fatalf("unexpected frame type: %T", frame)
	}

	if callError.ErrorCode != "ProtocolError" {
		t.Fatalf("unexpected error code: %q", callError.ErrorCode)
	}

	if callError.ErrorDetails["field"] != "idTag" {
		t.Fatalf("unexpected details: %#v", callError.ErrorDetails)
	}
}

func TestParseRejectsUnsupportedFrameType(t *testing.T) {
	t.Parallel()

	_, err := Parse([]byte(`[9,"19223201","Authorize",{}]`))
	if !errors.Is(err, ErrUnsupportedFrameType) {
		t.Fatalf("expected ErrUnsupportedFrameType, got %v", err)
	}
}

func TestRawCallMarshalJSON(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(RawCall{
		UniqueID: "19223201",
		Action:   "Authorize",
		Payload:  json.RawMessage(`{"idTag":"RFID-123"}`),
	})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	want := []any{
		float64(2),
		"19223201",
		"Authorize",
		map[string]any{"idTag": "RFID-123"},
	}

	assertJSONArrayEqual(t, data, want)
}

func TestCallErrorMarshalJSONDefaultsEmptyDetailsObject(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(CallError{
		UniqueID:         "19223201",
		ErrorCode:        "InternalError",
		ErrorDescription: "boom",
	})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	want := []any{
		float64(4),
		"19223201",
		"InternalError",
		"boom",
		map[string]any{},
	}

	assertJSONArrayEqual(t, data, want)
}

func assertJSONArrayEqual(t *testing.T, data []byte, want []any) {
	t.Helper()

	var got []any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal(got): %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected json:\n got: %#v\nwant: %#v", got, want)
	}
}
