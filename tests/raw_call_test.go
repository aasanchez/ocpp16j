package ocpp16json_test

import (
	"encoding/json"
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	expectedCallElements = 4
	nullPayloadStr       = "null"
)

// --- RawCall struct ---

func Test_RawCall_MessageType_ReturnsCall(t *testing.T) {
	t.Parallel()

	rawCall := ocpp16json.RawCall{
		UniqueId: testUniqueId,
		Action:   testAction,
		Payload:  json.RawMessage(emptyPayload),
	}

	if rawCall.MessageType() != ocpp16json.Call {
		t.Fatalf(
			errFmtIntExpGot,
			ocpp16json.Call,
			rawCall.MessageType(),
		)
	}
}

func Test_RawCall_MessageId_ReturnsValue(t *testing.T) {
	t.Parallel()

	rawCall := ocpp16json.RawCall{
		UniqueId: testUniqueId,
		Action:   testAction,
		Payload:  json.RawMessage(emptyPayload),
	}

	if rawCall.MessageId() != testUniqueIdStr {
		t.Fatalf(
			errFmtStrExpGot,
			testUniqueId,
			rawCall.MessageId(),
		)
	}
}

// --- NewRawCall ---

func Test_NewRawCall_Success(t *testing.T) {
	t.Parallel()

	payload := map[string]string{"key": "value"}

	rawCall, err := ocpp16json.NewRawCall(
		testUniqueId, testAction, payload,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if rawCall.UniqueId != testUniqueId {
		t.Fatalf(
			errFmtStrExpGot,
			testUniqueId, rawCall.UniqueId,
		)
	}

	if rawCall.Action != testAction {
		t.Fatalf(
			errFmtStrExpGot,
			testAction, rawCall.Action,
		)
	}

	if rawCall.Payload == nil {
		t.Fatal("expected non-nil payload")
	}
}

func Test_NewRawCall_EmptyAction(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.NewRawCall(
		testUniqueId, "", nil,
	)
	if !errors.Is(err, ocpp16json.ErrInvalidAction) {
		t.Fatalf(
			errFmtExpErrGot,
			ocpp16json.ErrInvalidAction, err,
		)
	}
}

func Test_NewRawCall_UnmarshalablePayload(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.NewRawCall(
		testUniqueId, testAction, make(chan int),
	)
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf(
			errFmtExpErrGot,
			ocpp16json.ErrPayloadDecode, err,
		)
	}
}

func Test_NewRawCall_NilPayload(t *testing.T) {
	t.Parallel()

	rawCall, err := ocpp16json.NewRawCall(
		testUniqueId, testAction, nil,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if string(rawCall.Payload) != nullPayloadStr {
		t.Fatalf(
			errFmtStrExpGot,
			nullPayloadStr, string(rawCall.Payload),
		)
	}
}

// --- RawCall MarshalJSON ---

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
		t.Fatalf(errFmtNilGot, err)
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

	data, marshalErr := json.Marshal(rawCall)
	if marshalErr != nil {
		t.Fatalf(errFmtNilGot, marshalErr)
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

	data, marshalErr := json.Marshal(rawCall)
	if marshalErr != nil {
		t.Fatalf(errFmtNilGot, marshalErr)
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
