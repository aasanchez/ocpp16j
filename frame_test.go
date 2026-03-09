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

func TestFrameMetadataHelpers(t *testing.T) {
	t.Parallel()

	call := RawCall{UniqueID: "call-id"}
	result := RawCallResult{UniqueID: "result-id"}
	callError := CallError{UniqueID: "error-id"}
	decodedCall := DecodedCall{UniqueID: "decoded-call-id"}
	decodedResult := DecodedCallResult{UniqueID: "decoded-result-id"}

	if call.MessageType() != MessageTypeCall || call.MessageID() != "call-id" {
		t.Fatal("unexpected raw call metadata")
	}

	if result.MessageType() != MessageTypeCallResult || result.MessageID() != "result-id" {
		t.Fatal("unexpected raw call result metadata")
	}

	if callError.MessageType() != MessageTypeCallError || callError.MessageID() != "error-id" {
		t.Fatal("unexpected call error metadata")
	}

	if decodedCall.MessageType() != MessageTypeCall || decodedCall.MessageID() != "decoded-call-id" {
		t.Fatal("unexpected decoded call metadata")
	}

	if decodedResult.MessageType() != MessageTypeCallResult || decodedResult.MessageID() != "decoded-result-id" {
		t.Fatal("unexpected decoded call result metadata")
	}
}

func TestParseRejectsInvalidJSONAndEmptyArray(t *testing.T) {
	t.Parallel()

	_, err := Parse([]byte(`{`))
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}

	_, err = Parse([]byte(`[]`))
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame for empty array, got %v", err)
	}
}

func TestRawCallMarshalJSONValidation(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(RawCall{Action: "Authorize", Payload: json.RawMessage(`{}`)})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected ErrInvalidMessageID, got %v", err)
	}

	_, err = json.Marshal(RawCall{UniqueID: "uid", Payload: json.RawMessage(`{}`)})
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction, got %v", err)
	}

	_, err = json.Marshal(RawCall{UniqueID: "uid", Action: "Authorize"})
	if !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf("expected ErrPayloadRequired, got %v", err)
	}
}

func TestRawCallResultMarshalJSONValidation(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(RawCallResult{})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected ErrInvalidMessageID, got %v", err)
	}

	_, err = json.Marshal(RawCallResult{UniqueID: "uid"})
	if !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf("expected ErrPayloadRequired, got %v", err)
	}
}

func TestRawCallResultMarshalJSON(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(RawCallResult{
		UniqueID: "uid",
		Payload:  json.RawMessage(`{"status":"Accepted"}`),
	})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	assertJSONArrayEqual(t, data, []any{
		float64(3),
		"uid",
		map[string]any{"status": "Accepted"},
	})
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

func TestCallErrorMarshalJSONValidation(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(CallError{ErrorCode: "ProtocolError", ErrorDescription: "bad"})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected ErrInvalidMessageID, got %v", err)
	}

	_, err = json.Marshal(CallError{UniqueID: "uid", ErrorDescription: "bad"})
	if !errors.Is(err, ErrErrorCodeRequired) {
		t.Fatalf("expected ErrErrorCodeRequired, got %v", err)
	}

	_, err = json.Marshal(CallError{UniqueID: "uid", ErrorCode: "ProtocolError"})
	if !errors.Is(err, ErrErrorDescriptionAbsent) {
		t.Fatalf("expected ErrErrorDescriptionAbsent, got %v", err)
	}
}

func TestParseCallValidationBranches(t *testing.T) {
	t.Parallel()

	_, err := parseCall([]json.RawMessage{
		json.RawMessage(`2`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"Authorize"`),
	})
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}

	_, err = parseCall([]json.RawMessage{
		json.RawMessage(`2`),
		json.RawMessage(`123`),
		json.RawMessage(`"Authorize"`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected ErrInvalidMessageID, got %v", err)
	}

	_, err = parseCall([]json.RawMessage{
		json.RawMessage(`2`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`123`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction, got %v", err)
	}

	_, err = parseCall([]json.RawMessage{
		json.RawMessage(`2`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"Authorize"`),
		json.RawMessage(` `),
	})
	if !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf("expected ErrPayloadRequired, got %v", err)
	}
}

func TestParseCallResultValidationBranches(t *testing.T) {
	t.Parallel()

	_, err := parseCallResult([]json.RawMessage{
		json.RawMessage(`3`),
		json.RawMessage(`"uid"`),
	})
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}

	_, err = parseCallResult([]json.RawMessage{
		json.RawMessage(`3`),
		json.RawMessage(`123`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected ErrInvalidMessageID, got %v", err)
	}

	_, err = parseCallResult([]json.RawMessage{
		json.RawMessage(`3`),
		json.RawMessage(`"uid"`),
		json.RawMessage(` `),
	})
	if !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf("expected ErrPayloadRequired, got %v", err)
	}
}

func TestParseCallErrorValidationBranches(t *testing.T) {
	t.Parallel()

	_, err := parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`"bad"`),
	})
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}

	_, err = parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`123`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`"bad"`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected ErrInvalidMessageID, got %v", err)
	}

	_, err = parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`123`),
		json.RawMessage(`"bad"`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrErrorCodeRequired) {
		t.Fatalf("expected ErrErrorCodeRequired, got %v", err)
	}

	_, err = parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`123`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrErrorDescriptionAbsent) {
		t.Fatalf("expected ErrErrorDescriptionAbsent, got %v", err)
	}

	_, err = parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`"bad"`),
		json.RawMessage(`[]`),
	})
	if !errors.Is(err, ErrErrorDetailsInvalid) {
		t.Fatalf("expected ErrErrorDetailsInvalid, got %v", err)
	}

	_, err = parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`"bad"`),
		json.RawMessage(`null`),
	})
	if !errors.Is(err, ErrErrorDetailsInvalid) {
		t.Fatalf("expected ErrErrorDetailsInvalid for null details, got %v", err)
	}
}

func TestDecodeHelpersAndPredicates(t *testing.T) {
	t.Parallel()

	if _, err := decodeMessageType(json.RawMessage(`"x"`)); !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}

	if _, err := decodeString(json.RawMessage(`123`), ErrInvalidAction); !errors.Is(err, ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction, got %v", err)
	}

	if _, err := decodeString(json.RawMessage(`""`), ErrInvalidAction); !errors.Is(err, ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction for empty string, got %v", err)
	}

	if _, err := decodeString(json.RawMessage(`"ok"`), ErrInvalidAction); err != nil {
		t.Fatalf("expected decodeString success, got %v", err)
	}

	if err := validateMessageID(""); !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected ErrInvalidMessageID, got %v", err)
	}

	if err := validateAction(""); !errors.Is(err, ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction, got %v", err)
	}

	data, err := marshalJSONArray(make(chan int))
	if !errors.Is(err, ErrInvalidFrame) || data != nil {
		t.Fatalf("expected marshal ErrInvalidFrame, got data=%v err=%v", data, err)
	}

	payload, err := DecodePayload[map[string]string](json.RawMessage(`{"k":"v"}`))
	if err != nil || payload["k"] != "v" {
		t.Fatalf("unexpected payload decode result: payload=%v err=%v", payload, err)
	}

	if _, err := DecodePayload[map[string]string](json.RawMessage(` `)); !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf("expected ErrPayloadRequired, got %v", err)
	}

	if _, err := DecodePayload[map[string]string](json.RawMessage(`1`)); !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf("expected ErrPayloadDecode, got %v", err)
	}

	call := RawCall{UniqueID: "uid", Action: "Authorize", Payload: json.RawMessage(`{}`)}
	result := RawCallResult{UniqueID: "uid", Payload: json.RawMessage(`{}`)}
	callError := CallError{
		UniqueID:         "uid",
		ErrorCode:        "ProtocolError",
		ErrorDescription: "bad",
		ErrorDetails:     map[string]any{},
	}

	if !IsCall(call) || IsCall(nil) || IsCall(result) {
		t.Fatal("IsCall predicate mismatch")
	}

	if !IsCallResult(result) || IsCallResult(nil) || IsCallResult(call) {
		t.Fatal("IsCallResult predicate mismatch")
	}

	if !IsCallError(callError) || IsCallError(nil) || IsCallError(call) {
		t.Fatal("IsCallError predicate mismatch")
	}

	if _, err := AsRawCall(result); err == nil {
		t.Fatal("expected AsRawCall type assertion error")
	}

	raw, err := AsRawCall(call)
	if err != nil || raw.UniqueID != "uid" {
		t.Fatalf("unexpected AsRawCall result: raw=%v err=%v", raw, err)
	}
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
