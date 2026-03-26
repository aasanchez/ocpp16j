package ocpp16json_test

import (
	"encoding/json"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

// Shared test constants used across all test files in this
// package. Defined here because message_test.go is the
// foundational test file.
const (
	expectedCall       ocpp16json.MessageType = 2
	expectedCallResult ocpp16json.MessageType = 3
	expectedCallError  ocpp16json.MessageType = 4
	testUniqueIdStr                           = "19223201"
	testUniqueId       ocpp16json.UniqueId    = testUniqueIdStr
	testAction                                = "Authorize"
	testErrorCode      ocpp16json.ErrorCode   = "NotImplemented"
	testErrorDesc                             = "Unknown action"
	emptyPayload                              = `{}`
	errFmtIntExpGot                           = "expected %d, got %d"
	errFmtStrExpGot                           = "expected %q, got %q"
	errFmtNilGot                              = "expected nil error, got %v"
	errFmtExpErrGot                           = "expected %v, got %v"
	errFmtElementCount                        = "expected %d elements, got %d"
)

// --- MessageType constants ---

func Test_Call_Equals_2(t *testing.T) {
	t.Parallel()

	if ocpp16json.Call != expectedCall {
		t.Fatalf(
			errFmtIntExpGot,
			expectedCall,
			ocpp16json.Call,
		)
	}
}

func Test_CallResult_Equals_3(t *testing.T) {
	t.Parallel()

	if ocpp16json.CallResult != expectedCallResult {
		t.Fatalf(
			errFmtIntExpGot,
			expectedCallResult,
			ocpp16json.CallResult,
		)
	}
}

func Test_CallError_Equals_4(t *testing.T) {
	t.Parallel()

	if ocpp16json.CallError != expectedCallError {
		t.Fatalf(
			errFmtIntExpGot,
			expectedCallError,
			ocpp16json.CallError,
		)
	}
}

// --- DecodedMessage (Call) ---

func Test_DecodedCall_MessageType_ReturnsCall(
	t *testing.T,
) {
	t.Parallel()

	decodedCall, err := ocpp16json.NewDecodedCall[any](
		testUniqueId, testAction, nil,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if decodedCall.MessageType() != ocpp16json.Call {
		t.Fatalf(
			errFmtIntExpGot,
			ocpp16json.Call,
			decodedCall.MessageType(),
		)
	}
}

func Test_DecodedCall_MessageId_ReturnsValue(
	t *testing.T,
) {
	t.Parallel()

	decodedCall, err := ocpp16json.NewDecodedCall[any](
		testUniqueId, testAction, nil,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if decodedCall.MessageId() != testUniqueIdStr {
		t.Fatalf(
			errFmtStrExpGot,
			testUniqueId,
			decodedCall.MessageId(),
		)
	}
}

func Test_DecodedCall_EmptyAction_ReturnsError(
	t *testing.T,
) {
	t.Parallel()

	_, err := ocpp16json.NewDecodedCall[any](
		testUniqueId, "", nil,
	)
	if err == nil {
		t.Fatal("expected error for empty action, got nil")
	}
}

// --- DecodedMessage (CallResult) ---

func Test_DecodedCallResult_MessageType_ReturnsCallResult(
	t *testing.T,
) {
	t.Parallel()

	decodedCallResult, err := ocpp16json.NewDecodedCallResult[any](
		testUniqueId, testAction, nil,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if decodedCallResult.MessageType() != ocpp16json.CallResult {
		t.Fatalf(
			errFmtIntExpGot,
			ocpp16json.CallResult,
			decodedCallResult.MessageType(),
		)
	}
}

func Test_DecodedCallResult_MessageId_ReturnsValue(
	t *testing.T,
) {
	t.Parallel()

	decodedCallResult, err := ocpp16json.NewDecodedCallResult[any](
		testUniqueId, testAction, nil,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if decodedCallResult.MessageId() != testUniqueIdStr {
		t.Fatalf(
			errFmtStrExpGot,
			testUniqueId,
			decodedCallResult.MessageId(),
		)
	}
}

func Test_DecodedCallResult_EmptyAction_ReturnsError(
	t *testing.T,
) {
	t.Parallel()

	_, err := ocpp16json.NewDecodedCallResult[any](
		testUniqueId, "", nil,
	)
	if err == nil {
		t.Fatal("expected error for empty action, got nil")
	}
}

// --- IsCall ---

func Test_IsCall_RawCall_ReturnsTrue(t *testing.T) {
	t.Parallel()

	rawCall := ocpp16json.RawCall{
		UniqueId: testUniqueId,
		Action:   testAction,
		Payload:  json.RawMessage(emptyPayload),
	}

	if !ocpp16json.IsCall(rawCall) {
		t.Fatal("expected IsCall to return true for RawCall")
	}
}

func Test_IsCall_DecodedCall_ReturnsTrue(t *testing.T) {
	t.Parallel()

	decodedCall, err := ocpp16json.NewDecodedCall[any](
		testUniqueId, testAction, nil,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if !ocpp16json.IsCall(decodedCall) {
		t.Fatal(
			"expected IsCall to return true for DecodedMessage",
		)
	}
}

func Test_IsCall_RawCallResult_ReturnsFalse(t *testing.T) {
	t.Parallel()

	rawCallResult := ocpp16json.RawCallResult{
		UniqueId: testUniqueId,
		Payload:  json.RawMessage(emptyPayload),
	}

	if ocpp16json.IsCall(rawCallResult) {
		t.Fatal(
			"expected IsCall false for RawCallResult",
		)
	}
}

// --- IsCallResult ---

func Test_IsCallResult_RawCallResult_ReturnsTrue(
	t *testing.T,
) {
	t.Parallel()

	rawCallResult := ocpp16json.RawCallResult{
		UniqueId: testUniqueId,
		Payload:  json.RawMessage(emptyPayload),
	}

	if !ocpp16json.IsCallResult(rawCallResult) {
		t.Fatal(
			"expected IsCallResult true for RawCallResult",
		)
	}
}

func Test_IsCallResult_DecodedCallResult_ReturnsTrue(
	t *testing.T,
) {
	t.Parallel()

	decodedCallResult, err := ocpp16json.NewDecodedCallResult[any](
		testUniqueId, testAction, nil,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if !ocpp16json.IsCallResult(decodedCallResult) {
		t.Fatal(
			"expected IsCallResult true for DecodedMessage",
		)
	}
}

func Test_IsCallResult_RawCall_ReturnsFalse(t *testing.T) {
	t.Parallel()

	rawCall := ocpp16json.RawCall{
		UniqueId: testUniqueId,
		Action:   testAction,
		Payload:  json.RawMessage(emptyPayload),
	}

	if ocpp16json.IsCallResult(rawCall) {
		t.Fatal(
			"expected IsCallResult false for RawCall",
		)
	}
}

// --- IsCallError ---

func Test_IsCallError_RawCallError_ReturnsTrue(
	t *testing.T,
) {
	t.Parallel()

	rawCallError := ocpp16json.RawCallError{
		UniqueId:         testUniqueId,
		ErrorCode:        testErrorCode,
		ErrorDescription: testErrorDesc,
		ErrorDetails:     map[string]any{},
	}

	if !ocpp16json.IsCallError(rawCallError) {
		t.Fatal(
			"expected IsCallError true for RawCallError",
		)
	}
}

func Test_IsCallError_RawCall_ReturnsFalse(t *testing.T) {
	t.Parallel()

	rawCall := ocpp16json.RawCall{
		UniqueId: testUniqueId,
		Action:   testAction,
		Payload:  json.RawMessage(emptyPayload),
	}

	if ocpp16json.IsCallError(rawCall) {
		t.Fatal(
			"expected IsCallError false for RawCall",
		)
	}
}

// --- AsRawCall ---

func Test_AsRawCall_RawCall_ReturnsValue(t *testing.T) {
	t.Parallel()

	rawCall := ocpp16json.RawCall{
		UniqueId: testUniqueId,
		Action:   testAction,
		Payload:  json.RawMessage(emptyPayload),
	}

	result, err := ocpp16json.AsRawCall(rawCall)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if result.UniqueId != testUniqueId {
		t.Fatalf(
			errFmtStrExpGot,
			testUniqueId,
			result.UniqueId,
		)
	}
}

func Test_AsRawCall_RawCallResult_ReturnsError(
	t *testing.T,
) {
	t.Parallel()

	rawCallResult := ocpp16json.RawCallResult{
		UniqueId: testUniqueId,
		Payload:  json.RawMessage(emptyPayload),
	}

	_, err := ocpp16json.AsRawCall(rawCallResult)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
