package tests_fuzz

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
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
	f.Add("ProtocolError", "bad payload", []byte(`{}`))
	f.Add("ProtocolError", "bad payload", []byte(`{"field":"idTag"}`))
	f.Add("ProtocolError", "bad payload", []byte(`[]`))
	f.Add("ProtocolError", "bad payload", []byte(`null`))
	f.Add("ProtocolError", "bad payload", []byte(`"oops"`))
	f.Add("", "bad payload", []byte(`{}`))
	f.Add("ProtocolError", "", []byte(`{}`))

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

		frame := append(
			[]byte(`[4,"uid-1",`),
			mustMarshalString(t, errorCode)...,
		)
		frame = append(frame, ',')
		frame = append(frame, mustMarshalString(t, errorDescription)...)
		frame = append(frame, ',')
		frame = append(frame, errorDetails...)
		frame = append(frame, ']')

		parsedFrame, err := ocpp16json.Parse(frame)
		if err != nil {
			assertParseErrorContract(t, err)

			return
		}

		callError, ok := parsedFrame.(ocpp16json.CallError)
		if !ok {
			t.Fatalf("unexpected parsed frame type: %T", parsedFrame)
		}

		if callError.ErrorCode == "" {
			t.Fatal("CallError has empty code after Parse success")
		}

		if callError.ErrorDescription == "" {
			t.Fatal("CallError has empty description after Parse success")
		}

		if callError.ErrorDetails == nil {
			t.Fatal("CallError has nil details after Parse success")
		}
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

	if frame.MessageID() == "" {
		t.Fatal("parsed frame has empty message id")
	}

	switch typedFrame := frame.(type) {
	case ocpp16json.RawCall:
		if !ocpp16json.IsCall(typedFrame) {
			t.Fatal("RawCall does not satisfy IsCall")
		}

		if typedFrame.Action == "" {
			t.Fatal("RawCall has empty action")
		}

		if len(bytes.TrimSpace(typedFrame.Payload)) == 0 {
			t.Fatal("RawCall has empty payload")
		}
	case ocpp16json.RawCallResult:
		if !ocpp16json.IsCallResult(typedFrame) {
			t.Fatal("RawCallResult does not satisfy IsCallResult")
		}

		if len(bytes.TrimSpace(typedFrame.Payload)) == 0 {
			t.Fatal("RawCallResult has empty payload")
		}
	case ocpp16json.CallError:
		if !ocpp16json.IsCallError(typedFrame) {
			t.Fatal("CallError does not satisfy IsCallError")
		}

		if typedFrame.ErrorCode == "" {
			t.Fatal("CallError has empty code")
		}

		if typedFrame.ErrorDescription == "" {
			t.Fatal("CallError has empty description")
		}

		if typedFrame.ErrorDetails == nil {
			t.Fatal("CallError has nil details")
		}
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
