package ocpp16json

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

const (
	testActionAuthorize = "Authorize"
	testMessageID       = "19223201"
	testStatusAccepted  = "Accepted"
	testUID             = "uid"
	testEmptyJSON       = `{}`
	testEmptyString     = ""
	errParseFmt         = "Parse: %v"
	errFrameTypeFmt     = "unexpected frame type: %T"
	errMessageTypeFmt   = "unexpected message type: %v"
	errMessageIDFmt     = "unexpected message id: %q"
	errMarshalFmt       = "json.Marshal: %v"
	errNeedInvalidIDFmt = "expected ErrInvalidMessageID, got %v"
	errNeedInvalidFrm   = "expected ErrInvalidFrame, got %v"
	errNeedPayloadFmt   = "expected ErrPayloadRequired, got %v"
	errNeedActionFmt    = "expected ErrInvalidAction, got %v"
)

func TestParseRawCall(t *testing.T) {
	t.Parallel()

	frame, err := Parse(
		[]byte(`[2,"19223201","Authorize",{"idTag":"RFID-123"}]`),
	)
	if err != nil {
		t.Fatalf(errParseFmt, err)
	}

	call, ok := frame.(RawCall)
	if !ok {
		t.Fatalf(errFrameTypeFmt, frame)
	}

	if call.UniqueID != testMessageID {
		t.Fatalf("unexpected id: %q", call.UniqueID)
	}

	if call.Action != testActionAuthorize {
		t.Fatalf("unexpected action: %q", call.Action)
	}

	if string(call.Payload) != `{"idTag":"RFID-123"}` {
		t.Fatalf("unexpected payload: %s", string(call.Payload))
	}
}

func TestParseCallResult(t *testing.T) {
	t.Parallel()

	frame, err := Parse(
		[]byte(`[3,"19223201",{"currentTime":"2025-01-02T15:04:05Z"}]`),
	)
	if err != nil {
		t.Fatalf(errParseFmt, err)
	}

	result, ok := frame.(RawCallResult)
	if !ok {
		t.Fatalf(errFrameTypeFmt, frame)
	}

	if result.UniqueID != testMessageID {
		t.Fatalf("unexpected id: %q", result.UniqueID)
	}
}

func TestParseCallError(t *testing.T) {
	t.Parallel()

	frame, err := Parse(
		[]byte(
			`[4,"19223201","ProtocolError","bad payload",{"field":"idTag"}]`,
		),
	)
	if err != nil {
		t.Fatalf(errParseFmt, err)
	}

	callError, ok := frame.(CallError)
	if !ok {
		t.Fatalf(errFrameTypeFmt, frame)
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

func TestParseRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := Parse([]byte(`{`))
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf(errNeedInvalidFrm, err)
	}
}

func TestParseRejectsEmptyArray(t *testing.T) {
	t.Parallel()

	_, err := Parse([]byte(`[]`))
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf(errNeedInvalidFrm, err)
	}
}

func TestRawCallMetadata(t *testing.T) {
	t.Parallel()

	rawCall := RawCall{
		UniqueID: "call-id",
		Action:   testEmptyString,
		Payload:  nil,
	}

	if rawCall.MessageType() != MessageTypeCall {
		t.Fatalf(errMessageTypeFmt, rawCall.MessageType())
	}

	if rawCall.MessageID() != "call-id" {
		t.Fatalf(errMessageIDFmt, rawCall.MessageID())
	}
}

func TestRawCallResultMetadata(t *testing.T) {
	t.Parallel()

	rawCallResult := RawCallResult{
		UniqueID: "result-id",
		Payload:  nil,
	}

	if rawCallResult.MessageType() != MessageTypeCallResult {
		t.Fatalf(errMessageTypeFmt, rawCallResult.MessageType())
	}

	if rawCallResult.MessageID() != "result-id" {
		t.Fatalf(errMessageIDFmt, rawCallResult.MessageID())
	}
}

func TestCallErrorMetadata(t *testing.T) {
	t.Parallel()

	callError := CallError{
		UniqueID:         "error-id",
		ErrorCode:        testEmptyString,
		ErrorDescription: testEmptyString,
		ErrorDetails:     nil,
	}

	if callError.MessageType() != MessageTypeCallError {
		t.Fatalf(errMessageTypeFmt, callError.MessageType())
	}

	if callError.MessageID() != "error-id" {
		t.Fatalf(errMessageIDFmt, callError.MessageID())
	}
}

func TestDecodedCallMetadata(t *testing.T) {
	t.Parallel()

	decodedCall := DecodedCall{
		UniqueID: "decoded-call-id",
		Action:   testEmptyString,
		Payload:  nil,
	}

	if decodedCall.MessageType() != MessageTypeCall {
		t.Fatalf(errMessageTypeFmt, decodedCall.MessageType())
	}

	if decodedCall.MessageID() != "decoded-call-id" {
		t.Fatalf(errMessageIDFmt, decodedCall.MessageID())
	}
}

func TestDecodedCallResultMetadata(t *testing.T) {
	t.Parallel()

	decodedCallResult := DecodedCallResult{
		UniqueID: "decoded-result-id",
		Action:   testEmptyString,
		Payload:  nil,
	}

	if decodedCallResult.MessageType() != MessageTypeCallResult {
		t.Fatalf(errMessageTypeFmt, decodedCallResult.MessageType())
	}

	if decodedCallResult.MessageID() != "decoded-result-id" {
		t.Fatalf(errMessageIDFmt, decodedCallResult.MessageID())
	}
}

func TestRawCallMarshalJSON(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(RawCall{
		UniqueID: testMessageID,
		Action:   testActionAuthorize,
		Payload:  json.RawMessage(`{"idTag":"RFID-123"}`),
	})
	if err != nil {
		t.Fatalf(errMarshalFmt, err)
	}

	assertJSONArrayEqual(t, data, []any{
		float64(2),
		testMessageID,
		testActionAuthorize,
		map[string]any{"idTag": "RFID-123"},
	})
}

func TestRawCallMarshalJSONRejectsMissingMessageID(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(RawCall{
		UniqueID: "",
		Action:   testActionAuthorize,
		Payload:  json.RawMessage(testEmptyJSON),
	})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf(errNeedInvalidIDFmt, err)
	}
}

func TestRawCallMarshalJSONRejectsMissingAction(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(RawCall{
		UniqueID: "uid",
		Action:   testEmptyString,
		Payload:  json.RawMessage(testEmptyJSON),
	})
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf(errNeedActionFmt, err)
	}
}

func TestRawCallMarshalJSONRejectsMissingPayload(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(RawCall{
		UniqueID: "uid",
		Action:   testActionAuthorize,
		Payload:  nil,
	})
	if !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf(errNeedPayloadFmt, err)
	}
}

func TestRawCallResultMarshalJSON(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(RawCallResult{
		UniqueID: "uid",
		Payload:  json.RawMessage(`{"status":"Accepted"}`),
	})
	if err != nil {
		t.Fatalf(errMarshalFmt, err)
	}

	assertJSONArrayEqual(t, data, []any{
		float64(3),
		"uid",
		map[string]any{"status": testStatusAccepted},
	})
}

func TestRawCallResultMarshalJSONRejectsMissingMessageID(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(RawCallResult{
		UniqueID: "",
		Payload:  json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf(errNeedInvalidIDFmt, err)
	}
}

func TestRawCallResultMarshalJSONRejectsMissingPayload(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(RawCallResult{
		UniqueID: "uid",
		Payload:  nil,
	})
	if !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf(errNeedPayloadFmt, err)
	}
}

func TestCallErrorMarshalJSONDefaultsEmptyDetailsObject(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(CallError{
		UniqueID:         testMessageID,
		ErrorCode:        "InternalError",
		ErrorDescription: "boom",
		ErrorDetails:     nil,
	})
	if err != nil {
		t.Fatalf(errMarshalFmt, err)
	}

	assertJSONArrayEqual(t, data, []any{
		float64(4),
		testMessageID,
		"InternalError",
		"boom",
		map[string]any{},
	})
}

func TestCallErrorMarshalJSONRejectsMissingMessageID(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(CallError{
		UniqueID:         "",
		ErrorCode:        "ProtocolError",
		ErrorDescription: "bad",
		ErrorDetails:     map[string]any{},
	})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf(errNeedInvalidIDFmt, err)
	}
}

func TestCallErrorMarshalJSONRejectsMissingCode(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(CallError{
		UniqueID:         "uid",
		ErrorCode:        "",
		ErrorDescription: "bad",
		ErrorDetails:     map[string]any{},
	})
	if !errors.Is(err, ErrErrorCodeRequired) {
		t.Fatalf("expected ErrErrorCodeRequired, got %v", err)
	}
}

func TestCallErrorMarshalJSONRejectsMissingDescription(t *testing.T) {
	t.Parallel()

	_, err := json.Marshal(CallError{
		UniqueID:         "uid",
		ErrorCode:        "ProtocolError",
		ErrorDescription: "",
		ErrorDetails:     map[string]any{},
	})
	if !errors.Is(err, ErrErrorDescriptionAbsent) {
		t.Fatalf("expected ErrErrorDescriptionAbsent, got %v", err)
	}
}

func TestParseCallRejectsWrongElementCount(t *testing.T) {
	t.Parallel()

	_, err := parseCall([]json.RawMessage{
		json.RawMessage(`2`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"Authorize"`),
	})
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf(errNeedInvalidFrm, err)
	}
}

func TestParseCallRejectsInvalidMessageID(t *testing.T) {
	t.Parallel()

	_, err := parseCall([]json.RawMessage{
		json.RawMessage(`2`),
		json.RawMessage(`123`),
		json.RawMessage(`"Authorize"`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf(errNeedInvalidIDFmt, err)
	}
}

func TestParseCallRejectsInvalidAction(t *testing.T) {
	t.Parallel()

	_, err := parseCall([]json.RawMessage{
		json.RawMessage(`2`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`123`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf(errNeedActionFmt, err)
	}
}

func TestParseCallRejectsMissingPayload(t *testing.T) {
	t.Parallel()

	_, err := parseCall([]json.RawMessage{
		json.RawMessage(`2`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"Authorize"`),
		json.RawMessage(` `),
	})
	if !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf(errNeedPayloadFmt, err)
	}
}

func TestParseCallResultRejectsWrongElementCount(t *testing.T) {
	t.Parallel()

	_, err := parseCallResult([]json.RawMessage{
		json.RawMessage(`3`),
		json.RawMessage(`"uid"`),
	})
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf(errNeedInvalidFrm, err)
	}
}

func TestParseCallResultRejectsInvalidMessageID(t *testing.T) {
	t.Parallel()

	_, err := parseCallResult([]json.RawMessage{
		json.RawMessage(`3`),
		json.RawMessage(`123`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf(errNeedInvalidIDFmt, err)
	}
}

func TestParseCallResultRejectsMissingPayload(t *testing.T) {
	t.Parallel()

	_, err := parseCallResult([]json.RawMessage{
		json.RawMessage(`3`),
		json.RawMessage(`"uid"`),
		json.RawMessage(` `),
	})
	if !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf(errNeedPayloadFmt, err)
	}
}

func TestParseCallErrorRejectsWrongElementCount(t *testing.T) {
	t.Parallel()

	_, err := parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`"bad"`),
	})
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf(errNeedInvalidFrm, err)
	}
}

func TestParseCallErrorRejectsInvalidMessageID(t *testing.T) {
	t.Parallel()

	_, err := parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`123`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`"bad"`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf(errNeedInvalidIDFmt, err)
	}
}

func TestParseCallErrorRejectsInvalidCode(t *testing.T) {
	t.Parallel()

	_, err := parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`123`),
		json.RawMessage(`"bad"`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrErrorCodeRequired) {
		t.Fatalf("expected ErrErrorCodeRequired, got %v", err)
	}
}

func TestParseCallErrorRejectsInvalidDescription(t *testing.T) {
	t.Parallel()

	_, err := parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`123`),
		json.RawMessage(`{}`),
	})
	if !errors.Is(err, ErrErrorDescriptionAbsent) {
		t.Fatalf("expected ErrErrorDescriptionAbsent, got %v", err)
	}
}

func TestParseCallErrorRejectsArrayDetails(t *testing.T) {
	t.Parallel()

	_, err := parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`"bad"`),
		json.RawMessage(`[]`),
	})
	if !errors.Is(err, ErrErrorDetailsInvalid) {
		t.Fatalf("expected ErrErrorDetailsInvalid, got %v", err)
	}
}

func TestParseCallErrorRejectsNullDetails(t *testing.T) {
	t.Parallel()

	_, err := parseCallError([]json.RawMessage{
		json.RawMessage(`4`),
		json.RawMessage(`"uid"`),
		json.RawMessage(`"ProtocolError"`),
		json.RawMessage(`"bad"`),
		json.RawMessage(`null`),
	})
	if !errors.Is(err, ErrErrorDetailsInvalid) {
		t.Fatalf("expected ErrErrorDetailsInvalid, got %v", err)
	}
}

func TestDecodeMessageTypeRejectsWrongJSONType(t *testing.T) {
	t.Parallel()

	_, err := decodeMessageType(json.RawMessage(`"x"`))
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}
}

func TestDecodeStringRejectsWrongJSONType(t *testing.T) {
	t.Parallel()

	_, err := decodeString(json.RawMessage(`123`), ErrInvalidAction)
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf(errNeedActionFmt, err)
	}
}

func TestDecodeStringRejectsEmptyString(t *testing.T) {
	t.Parallel()

	_, err := decodeString(json.RawMessage(`""`), ErrInvalidAction)
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf(errNeedActionFmt, err)
	}
}

func TestDecodeStringReturnsValue(t *testing.T) {
	t.Parallel()

	value, err := decodeString(json.RawMessage(`"ok"`), ErrInvalidAction)
	if err != nil {
		t.Fatalf("decodeString: %v", err)
	}

	if value != "ok" {
		t.Fatalf("unexpected value: %q", value)
	}
}

func TestValidateMessageIDRejectsEmptyValue(t *testing.T) {
	t.Parallel()

	err := validateMessageID("")
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected ErrInvalidMessageID, got %v", err)
	}
}

func TestValidateActionRejectsEmptyValue(t *testing.T) {
	t.Parallel()

	err := validateAction("")
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf("expected ErrInvalidAction, got %v", err)
	}
}

func TestMarshalJSONArrayRejectsUnsupportedValue(t *testing.T) {
	t.Parallel()

	data, err := marshalJSONArray(make(chan int))
	if !errors.Is(err, ErrInvalidFrame) {
		t.Fatalf("expected ErrInvalidFrame, got %v", err)
	}

	if data != nil {
		t.Fatalf("expected nil data, got %v", data)
	}
}

func TestDecodePayloadReturnsValue(t *testing.T) {
	t.Parallel()

	payload, err := DecodePayload[map[string]string](
		json.RawMessage(`{"k":"v"}`),
	)
	if err != nil {
		t.Fatalf("DecodePayload: %v", err)
	}

	if payload["k"] != "v" {
		t.Fatalf("unexpected payload: %v", payload)
	}
}

func TestDecodePayloadRejectsMissingPayload(t *testing.T) {
	t.Parallel()

	_, err := DecodePayload[map[string]string](json.RawMessage(` `))
	if !errors.Is(err, ErrPayloadRequired) {
		t.Fatalf("expected ErrPayloadRequired, got %v", err)
	}
}

func TestDecodePayloadRejectsInvalidPayload(t *testing.T) {
	t.Parallel()

	_, err := DecodePayload[map[string]string](json.RawMessage(`1`))
	if !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf("expected ErrPayloadDecode, got %v", err)
	}
}

func TestIsCall(t *testing.T) {
	t.Parallel()

	rawCall := RawCall{
		UniqueID: testUID,
		Action:   testActionAuthorize,
		Payload:  json.RawMessage(testEmptyJSON),
	}

	if !IsCall(rawCall) {
		t.Fatal("expected IsCall to return true")
	}

	if IsCall(nil) {
		t.Fatal("expected IsCall to return false for nil")
	}
}

func TestIsCallResult(t *testing.T) {
	t.Parallel()

	rawCallResult := RawCallResult{
		UniqueID: testUID,
		Payload:  json.RawMessage(testEmptyJSON),
	}

	if !IsCallResult(rawCallResult) {
		t.Fatal("expected IsCallResult to return true")
	}

	if IsCallResult(nil) {
		t.Fatal("expected IsCallResult to return false for nil")
	}
}

func TestIsCallError(t *testing.T) {
	t.Parallel()

	callError := CallError{
		UniqueID:         testUID,
		ErrorCode:        "ProtocolError",
		ErrorDescription: "bad",
		ErrorDetails:     map[string]any{},
	}

	if !IsCallError(callError) {
		t.Fatal("expected IsCallError to return true")
	}

	if IsCallError(nil) {
		t.Fatal("expected IsCallError to return false for nil")
	}
}

func TestAsRawCallRejectsWrongType(t *testing.T) {
	t.Parallel()

	_, err := AsRawCall(RawCallResult{
		UniqueID: testUID,
		Payload:  json.RawMessage(testEmptyJSON),
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAsRawCallReturnsValue(t *testing.T) {
	t.Parallel()

	expected := RawCall{
		UniqueID: testUID,
		Action:   testActionAuthorize,
		Payload:  json.RawMessage(testEmptyJSON),
	}

	actual, err := AsRawCall(expected)
	if err != nil {
		t.Fatalf("AsRawCall: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unexpected raw call: %#v", actual)
	}
}

func assertJSONArrayEqual(t *testing.T, data []byte, want []any) {
	t.Helper()

	var got []any

	err := json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("json.Unmarshal(got): %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected json:\n got: %#v\nwant: %#v", got, want)
	}
}
