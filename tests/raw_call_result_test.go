package ocpp16json_test

import (
	"encoding/json"
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const expectedCallResultElements = 3

// --- RawCallResult struct ---

func Test_RawCallResult_MessageType_ReturnsCallResult(
	t *testing.T,
) {
	t.Parallel()

	rawCallResult := ocpp16json.RawCallResult{
		UniqueId: testUniqueId,
		Payload:  json.RawMessage(emptyPayload),
	}

	if rawCallResult.MessageType() != ocpp16json.CallResult {
		t.Fatalf(
			errFmtIntExpGot,
			ocpp16json.CallResult,
			rawCallResult.MessageType(),
		)
	}
}

func Test_RawCallResult_MessageId_ReturnsValue(
	t *testing.T,
) {
	t.Parallel()

	rawCallResult := ocpp16json.RawCallResult{
		UniqueId: testUniqueId,
		Payload:  json.RawMessage(emptyPayload),
	}

	if rawCallResult.MessageId() != testUniqueIdStr {
		t.Fatalf(
			errFmtStrExpGot,
			testUniqueId,
			rawCallResult.MessageId(),
		)
	}
}

// --- NewRawCallResult ---

func Test_NewRawCallResult_Success(t *testing.T) {
	t.Parallel()

	payload := map[string]string{"status": "Accepted"}

	rawCallResult, err := ocpp16json.NewRawCallResult(
		testUniqueId, payload,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if rawCallResult.UniqueId != testUniqueId {
		t.Fatalf(
			errFmtStrExpGot,
			testUniqueId, rawCallResult.UniqueId,
		)
	}

	if rawCallResult.Payload == nil {
		t.Fatal("expected non-nil payload")
	}
}

func Test_NewRawCallResult_UnmarshalablePayload(
	t *testing.T,
) {
	t.Parallel()

	_, err := ocpp16json.NewRawCallResult(
		testUniqueId, make(chan int),
	)
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf(
			errFmtExpErrGot,
			ocpp16json.ErrPayloadDecode, err,
		)
	}
}

// --- RawCallResult MarshalJSON ---

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
		t.Fatalf(errFmtNilGot, err)
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

	data, marshalErr := json.Marshal(rawCallResult)
	if marshalErr != nil {
		t.Fatalf(errFmtNilGot, marshalErr)
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
