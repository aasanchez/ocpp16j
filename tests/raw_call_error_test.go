package ocpp16json_test

import (
	"encoding/json"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	expectedCallErrorElements = 5
	emptyDescription          = ""
	expectedEmptyDetails      = 0
)

// --- RawCallError struct ---

func Test_RawCallError_MessageType_ReturnsCallError(
	t *testing.T,
) {
	t.Parallel()

	rawCallError := ocpp16json.RawCallError{
		UniqueId:         testUniqueId,
		ErrorCode:        testErrorCode,
		ErrorDescription: testErrorDesc,
		ErrorDetails:     map[string]any{},
	}

	if rawCallError.MessageType() != ocpp16json.CallError {
		t.Fatalf(
			errFmtIntExpGot,
			ocpp16json.CallError,
			rawCallError.MessageType(),
		)
	}
}

func Test_RawCallError_MessageId_ReturnsValue(
	t *testing.T,
) {
	t.Parallel()

	rawCallError := ocpp16json.RawCallError{
		UniqueId:         testUniqueId,
		ErrorCode:        testErrorCode,
		ErrorDescription: testErrorDesc,
		ErrorDetails:     map[string]any{},
	}

	if rawCallError.MessageId() != testUniqueIdStr {
		t.Fatalf(
			errFmtStrExpGot,
			testUniqueId,
			rawCallError.MessageId(),
		)
	}
}

// --- NewRawCallError ---

func Test_NewRawCallError_Success(t *testing.T) {
	t.Parallel()

	rawCallError, err := ocpp16json.NewRawCallError(
		testUniqueId,
		ocpp16json.NotImplemented,
		testErrorDesc,
		map[string]any{},
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if rawCallError.ErrorCode != ocpp16json.NotImplemented {
		t.Fatalf(
			errFmtStrExpGot,
			ocpp16json.NotImplemented,
			rawCallError.ErrorCode,
		)
	}
}

func Test_NewRawCallError_EmptyDescription(t *testing.T) {
	t.Parallel()

	rawCallError, err := ocpp16json.NewRawCallError(
		testUniqueId,
		ocpp16json.GenericError,
		emptyDescription,
		nil,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if rawCallError.ErrorDescription != emptyDescription {
		t.Fatalf(
			"expected empty description, got %q",
			rawCallError.ErrorDescription,
		)
	}
}

func Test_NewRawCallError_NilDetails_DefaultsToEmpty(
	t *testing.T,
) {
	t.Parallel()

	rawCallError, err := ocpp16json.NewRawCallError(
		testUniqueId,
		ocpp16json.GenericError,
		testErrorDesc,
		nil,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if rawCallError.ErrorDetails == nil {
		t.Fatal("expected non-nil ErrorDetails")
	}

	if len(rawCallError.ErrorDetails) != expectedEmptyDetails {
		t.Fatalf(
			"expected empty ErrorDetails, got %v",
			rawCallError.ErrorDetails,
		)
	}
}

// --- RawCallError MarshalJSON ---

func Test_RawCallError_MarshalJSON_ProducesValidArray(
	t *testing.T,
) {
	t.Parallel()

	rawCallError, _ := ocpp16json.NewRawCallError(
		testUniqueId,
		ocpp16json.GenericError,
		testErrorDesc,
		nil,
	)

	data, err := json.Marshal(rawCallError)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	var elements []json.RawMessage

	_ = json.Unmarshal(data, &elements)

	if len(elements) != expectedCallErrorElements {
		t.Fatalf(
			errFmtElementCount,
			expectedCallErrorElements, len(elements),
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

	data, marshalErr := json.Marshal(rawCallError)
	if marshalErr != nil {
		t.Fatalf(errFmtNilGot, marshalErr)
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

	data, marshalErr := json.Marshal(rawCallError)
	if marshalErr != nil {
		t.Fatalf(errFmtNilGot, marshalErr)
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
