package tests_fuzz

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	fuzzProtocolError = "ProtocolError"
	fuzzEmptyObject   = `{}`
	fuzzEmptyString   = ""
	fuzzBadPayload    = "bad payload"
	fuzzZeroLength    = 0
)

func FuzzParseFrameEnvelope(f *testing.F) {
	f.Add([]byte(`[2,"uid-1","Authorize",{"idTag":"RFID-123"}]`))
	f.Add([]byte(`[3,"uid-1",{"currentTime":"2025-01-02T15:04:05Z"}]`))
	f.Add([]byte(`[4,"uid-1","ProtocolError","bad payload",{}]`))
	f.Add([]byte(`[2,"","",""]`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`[9,"uid-1","Authorize",{}]`))

	f.Fuzz(func(t *testing.T, raw []byte) {
		frame, err := ocpp16json.Parse(raw)
		if err != nil {
			assertParseErrorContract(t, err)

			return
		}

		assertFrameInvariant(t, frame)
		assertFrameRoundTrip(t, frame)
	})
}

func FuzzParseCallErrorDetails(f *testing.F) {
	f.Add(fuzzProtocolError, fuzzBadPayload, []byte(fuzzEmptyObject))
	f.Add(fuzzProtocolError, fuzzBadPayload, []byte(`{"field":"idTag"}`))
	f.Add(fuzzProtocolError, fuzzBadPayload, []byte(`[]`))
	f.Add(fuzzProtocolError, fuzzBadPayload, []byte(`null`))
	f.Add(fuzzProtocolError, fuzzBadPayload, []byte(`"oops"`))
	f.Add(fuzzEmptyString, fuzzBadPayload, []byte(fuzzEmptyObject))
	f.Add(fuzzProtocolError, fuzzEmptyString, []byte(fuzzEmptyObject))

	f.Fuzz(func(
		t *testing.T,
		errorCode string,
		errorDescription string,
		errorDetails []byte,
	) {
		if len(errorCode) > maxFuzzStringLength ||
			len(errorDescription) > maxFuzzStringLength ||
			len(errorDetails) > maxFuzzStringLength {
			t.Skip()
		}

		parsedFrame, err := ocpp16json.Parse(
			buildCallErrorFrame(t, errorCode, errorDescription, errorDetails),
		)
		if err != nil {
			assertParseErrorContract(t, err)

			return
		}

		assertCallErrorInvariant(t, parsedFrame)
	})
}

func assertParseErrorContract(t *testing.T, err error) {
	t.Helper()

	switch {
	case errors.Is(err, ocpp16json.ErrInvalidFrame):
	case errors.Is(err, ocpp16json.ErrUnsupportedFrameType):
	case errors.Is(err, ocpp16json.ErrInvalidMessageID):
	case errors.Is(err, ocpp16json.ErrInvalidAction):
	case errors.Is(err, ocpp16json.ErrPayloadRequired):
	case errors.Is(err, ocpp16json.ErrErrorCodeRequired):
	case errors.Is(err, ocpp16json.ErrErrorDescriptionAbsent):
	case errors.Is(err, ocpp16json.ErrErrorDetailsInvalid):
	default:
		t.Fatalf("unexpected Parse error contract: %v", err)
	}
}

func assertFrameInvariant(t *testing.T, frame ocpp16json.Frame) {
	t.Helper()

	if frame == nil {
		t.Fatal("Parse returned nil frame without error")
	}

	if frame.MessageID() == fuzzEmptyString {
		t.Fatal("parsed frame has empty message id")
	}

	switch typedFrame := frame.(type) {
	case ocpp16json.RawCall:
		assertRawCallInvariant(t, typedFrame)
	case ocpp16json.RawCallResult:
		assertRawCallResultInvariant(t, typedFrame)
	case ocpp16json.CallError:
		assertCallErrorValueInvariant(t, typedFrame)
	default:
		t.Fatalf("unexpected frame type: %T", frame)
	}
}

func assertFrameRoundTrip(t *testing.T, frame ocpp16json.Frame) {
	t.Helper()

	originalJSON, err := json.Marshal(frame)
	if err != nil {
		t.Fatalf("json.Marshal(frame): %v", err)
	}

	reparsed, err := ocpp16json.Parse(originalJSON)
	if err != nil {
		t.Fatalf("Parse(roundtrip): %v (json=%s)", err, string(originalJSON))
	}

	reparsedJSON, err := json.Marshal(reparsed)
	if err != nil {
		t.Fatalf("json.Marshal(reparsed): %v", err)
	}

	assertJSONSemanticallyEqual(t, originalJSON, reparsedJSON)
}

func assertJSONSemanticallyEqual(t *testing.T, left []byte, right []byte) {
	t.Helper()

	var leftValue any

	err := json.Unmarshal(left, &leftValue)
	if err != nil {
		t.Fatalf("json.Unmarshal(left): %v", err)
	}

	var rightValue any

	err = json.Unmarshal(right, &rightValue)
	if err != nil {
		t.Fatalf("json.Unmarshal(right): %v", err)
	}

	if !jsonValuesEqual(leftValue, rightValue) {
		t.Fatalf(
			"semantic roundtrip mismatch:\nleft:  %s\nright: %s",
			string(left),
			string(right),
		)
	}
}

func jsonValuesEqual(left any, right any) bool {
	leftJSON, leftErr := json.Marshal(left)
	if leftErr != nil {
		return false
	}

	rightJSON, rightErr := json.Marshal(right)
	if rightErr != nil {
		return false
	}

	return bytes.Equal(leftJSON, rightJSON)
}

func mustMarshalString(t *testing.T, value string) []byte {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal(string): %v", err)
	}

	return data
}

func buildCallErrorFrame(
	t *testing.T,
	errorCode string,
	errorDescription string,
	errorDetails []byte,
) []byte {
	t.Helper()

	frame := append(
		[]byte(`[4,"uid-1",`),
		mustMarshalString(t, errorCode)...,
	)
	frame = append(frame, ',')
	frame = append(frame, mustMarshalString(t, errorDescription)...)
	frame = append(frame, ',')
	frame = append(frame, errorDetails...)

	return append(frame, ']')
}

func assertCallErrorInvariant(t *testing.T, frame ocpp16json.Frame) {
	t.Helper()

	callError, ok := frame.(ocpp16json.CallError)
	if !ok {
		t.Fatalf("unexpected parsed frame type: %T", frame)
	}

	assertCallErrorValueInvariant(t, callError)
}

func assertRawCallInvariant(t *testing.T, call ocpp16json.RawCall) {
	t.Helper()

	if !ocpp16json.IsCall(call) {
		t.Fatal("RawCall does not satisfy IsCall")
	}

	if call.Action == fuzzEmptyString {
		t.Fatal("RawCall has empty action")
	}

	assertPayloadPresent(t, call.Payload, "RawCall has empty payload")
}

func assertRawCallResultInvariant(
	t *testing.T,
	result ocpp16json.RawCallResult,
) {
	t.Helper()

	if !ocpp16json.IsCallResult(result) {
		t.Fatal("RawCallResult does not satisfy IsCallResult")
	}

	assertPayloadPresent(t, result.Payload, "RawCallResult has empty payload")
}

func assertCallErrorValueInvariant(
	t *testing.T,
	callError ocpp16json.CallError,
) {
	t.Helper()

	if !ocpp16json.IsCallError(callError) {
		t.Fatal("CallError does not satisfy IsCallError")
	}

	if callError.ErrorCode == fuzzEmptyString {
		t.Fatal("CallError has empty code")
	}

	if callError.ErrorDescription == fuzzEmptyString {
		t.Fatal("CallError has empty description")
	}

	if callError.ErrorDetails == nil {
		t.Fatal("CallError has nil details")
	}
}

func assertPayloadPresent(
	t *testing.T,
	payload json.RawMessage,
	message string,
) {
	t.Helper()

	if len(bytes.TrimSpace(payload)) == fuzzZeroLength {
		t.Fatal(message)
	}
}
