package ocpp16json_test

import (
	"encoding/json"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	expectedCallElements       = 4
	expectedCallResultElements = 3
	expectedCallErrorElements  = 5
	errFmtMarshalNilGot        = "expected nil error, got %v"
	errFmtElementCount         = "expected %d elements, got %d"
)

func Test_RawCall_MarshalJSON_ProducesValidArray(
	t *testing.T,
) {
	t.Parallel()

	rawCall := ocpp16json.RawCall{
		UniqueId: testUniqueId,
		Action:   testAction,
		Payload:  json.RawMessage(emptyPayload),
	}

	data, err := json.Marshal(rawCall)
	if err != nil {
		t.Fatalf(errFmtMarshalNilGot, err)
	}

	var elements []json.RawMessage

	unmarshalErr := json.Unmarshal(data, &elements)
	if unmarshalErr != nil {
		t.Fatalf("invalid JSON array: %v", unmarshalErr)
	}

	if len(elements) != expectedCallElements {
		t.Fatalf(
			errFmtElementCount,
			expectedCallElements,
			len(elements),
		)
	}
}

func Test_RawCall_MarshalJSON_CorrectMessageTypeId(
	t *testing.T,
) {
	t.Parallel()

	rawCall := ocpp16json.RawCall{
		UniqueId: testUniqueId,
		Action:   testAction,
		Payload:  json.RawMessage(emptyPayload),
	}

	data, err := json.Marshal(rawCall)
	if err != nil {
		t.Fatalf(errFmtMarshalNilGot, err)
	}

	var elements []json.RawMessage

	_ = json.Unmarshal(data, &elements)

	var messageTypeId uint8

	_ = json.Unmarshal(elements[0], &messageTypeId)

	if messageTypeId != uint8(ocpp16json.Call) {
		t.Fatalf(
			errFmtIntExpGot,
			ocpp16json.Call,
			messageTypeId,
		)
	}
}

func Test_RawCall_MarshalJSON_CorrectUniqueId(
	t *testing.T,
) {
	t.Parallel()

	rawCall := ocpp16json.RawCall{
		UniqueId: testUniqueId,
		Action:   testAction,
		Payload:  json.RawMessage(emptyPayload),
	}

	data, err := json.Marshal(rawCall)
	if err != nil {
		t.Fatalf(errFmtMarshalNilGot, err)
	}

	var elements []json.RawMessage

	_ = json.Unmarshal(data, &elements)

	var uniqueId string

	_ = json.Unmarshal(elements[1], &uniqueId)

	if uniqueId != testUniqueIdStr {
		t.Fatalf(
			errFmtStrExpGot, testUniqueIdStr, uniqueId,
		)
	}
}

func Test_RawCallResult_MarshalJSON_ProducesValidArray(
	t *testing.T,
) {
	t.Parallel()

	rawCallResult := ocpp16json.RawCallResult{
		UniqueId: testUniqueId,
		Payload:  json.RawMessage(emptyPayload),
	}

	data, err := json.Marshal(rawCallResult)
	if err != nil {
		t.Fatalf(errFmtMarshalNilGot, err)
	}

	var elements []json.RawMessage

	_ = json.Unmarshal(data, &elements)

	if len(elements) != expectedCallResultElements {
		t.Fatalf(
			errFmtElementCount,
			expectedCallResultElements,
			len(elements),
		)
	}
}

func Test_RawCallResult_MarshalJSON_CorrectMessageTypeId(
	t *testing.T,
) {
	t.Parallel()

	rawCallResult := ocpp16json.RawCallResult{
		UniqueId: testUniqueId,
		Payload:  json.RawMessage(emptyPayload),
	}

	data, err := json.Marshal(rawCallResult)
	if err != nil {
		t.Fatalf(errFmtMarshalNilGot, err)
	}

	var elements []json.RawMessage

	_ = json.Unmarshal(data, &elements)

	var messageTypeId uint8

	_ = json.Unmarshal(elements[0], &messageTypeId)

	if messageTypeId != uint8(ocpp16json.CallResult) {
		t.Fatalf(
			errFmtIntExpGot,
			ocpp16json.CallResult,
			messageTypeId,
		)
	}
}

func Test_RawCallError_MarshalJSON_ProducesValidArray(
	t *testing.T,
) {
	t.Parallel()

	rawCallError := ocpp16json.RawCallError{
		UniqueId:         testUniqueId,
		ErrorCode:        testErrorCode,
		ErrorDescription: testErrorDesc,
		ErrorDetails:     map[string]any{},
	}

	data, err := json.Marshal(rawCallError)
	if err != nil {
		t.Fatalf(errFmtMarshalNilGot, err)
	}

	var elements []json.RawMessage

	_ = json.Unmarshal(data, &elements)

	if len(elements) != expectedCallErrorElements {
		t.Fatalf(
			errFmtElementCount,
			expectedCallErrorElements,
			len(elements),
		)
	}
}

func Test_RawCallError_MarshalJSON_CorrectErrorCode(
	t *testing.T,
) {
	t.Parallel()

	rawCallError := ocpp16json.RawCallError{
		UniqueId:         testUniqueId,
		ErrorCode:        testErrorCode,
		ErrorDescription: testErrorDesc,
		ErrorDetails:     map[string]any{},
	}

	data, err := json.Marshal(rawCallError)
	if err != nil {
		t.Fatalf(errFmtMarshalNilGot, err)
	}

	var elements []json.RawMessage

	_ = json.Unmarshal(data, &elements)

	var errorCode string

	_ = json.Unmarshal(elements[2], &errorCode)

	if errorCode != testErrorCode.String() {
		t.Fatalf(
			errFmtStrExpGot,
			testErrorCode.String(),
			errorCode,
		)
	}
}

func Test_RawCallError_MarshalJSON_CorrectMessageTypeId(
	t *testing.T,
) {
	t.Parallel()

	rawCallError := ocpp16json.RawCallError{
		UniqueId:         testUniqueId,
		ErrorCode:        testErrorCode,
		ErrorDescription: testErrorDesc,
		ErrorDetails:     map[string]any{},
	}

	data, err := json.Marshal(rawCallError)
	if err != nil {
		t.Fatalf(errFmtMarshalNilGot, err)
	}

	var elements []json.RawMessage

	_ = json.Unmarshal(data, &elements)

	var messageTypeId uint8

	_ = json.Unmarshal(elements[0], &messageTypeId)

	if messageTypeId != uint8(ocpp16json.CallError) {
		t.Fatalf(
			errFmtIntExpGot,
			ocpp16json.CallError,
			messageTypeId,
		)
	}
}
