package ocpp16json

import (
	"encoding/json"
	"errors"
	"testing"
)

const (
	testUniqueId      = "19223201"
	testAction        = "Authorize"
	expectedElements  = 3
	errFmtNilGot      = "expected nil error, got %v"
	errFmtActionGot   = "expected ErrInvalidAction, got %v"
	errFmtExpectedGot = "expected %q, got %q"
)

// --- validateUniqueId ---

func Test_validateUniqueId_Empty(t *testing.T) {
	t.Parallel()

	err := validateUniqueId("")
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected ErrInvalidMessageID, got %v", err)
	}
}

func Test_validateUniqueId_Valid(t *testing.T) {
	t.Parallel()

	err := validateUniqueId(testUniqueId)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}
}

// --- validateAction ---

func Test_validateAction_Empty(t *testing.T) {
	t.Parallel()

	err := validateAction("")
	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf(errFmtActionGot, err)
	}
}

func Test_validateAction_Valid(t *testing.T) {
	t.Parallel()

	err := validateAction(testAction)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}
}

// --- decodeString ---

func Test_decodeString_Valid(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`"Authorize"`)

	value, err := decodeString(raw, ErrInvalidAction)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if value != testAction {
		t.Fatalf(errFmtExpectedGot, testAction, value)
	}
}

func Test_decodeString_WrongJSONType(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`42`)
	_, err := decodeString(raw, ErrInvalidAction)

	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf(errFmtActionGot, err)
	}
}

func Test_decodeString_EmptyString(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`""`)
	_, err := decodeString(raw, ErrInvalidAction)

	if !errors.Is(err, ErrInvalidAction) {
		t.Fatalf(errFmtActionGot, err)
	}
}

// --- decodeMessageType ---

func Test_decodeMessageType_Call(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`2`)

	messageType, err := decodeMessageType(raw)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if messageType != Call {
		t.Fatalf("expected Call (2), got %d", messageType)
	}
}

func Test_decodeMessageType_CallResult(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`3`)

	messageType, err := decodeMessageType(raw)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if messageType != CallResult {
		t.Fatalf(
			"expected CallResult (3), got %d", messageType,
		)
	}
}

func Test_decodeMessageType_CallError(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`4`)

	messageType, err := decodeMessageType(raw)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	if messageType != CallError {
		t.Fatalf(
			"expected CallError (4), got %d", messageType,
		)
	}
}

func Test_decodeMessageType_Unknown(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`9`)
	_, err := decodeMessageType(raw)

	if !errors.Is(err, ErrUnsupportedMessageType) {
		t.Fatalf(
			"expected ErrUnsupportedMessageType, got %v",
			err,
		)
	}
}

func Test_decodeMessageType_WrongJSONType(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`"two"`)
	_, err := decodeMessageType(raw)

	if !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("expected ErrInvalidMessage, got %v", err)
	}
}

// --- marshalJSONArray ---

func Test_marshalJSONArray_Valid(t *testing.T) {
	t.Parallel()

	data, err := marshalJSONArray(
		Call, testUniqueId, testAction,
	)
	if err != nil {
		t.Fatalf(errFmtNilGot, err)
	}

	var elements []json.RawMessage

	unmarshalErr := json.Unmarshal(data, &elements)
	if unmarshalErr != nil {
		t.Fatalf("invalid JSON array: %v", unmarshalErr)
	}

	if len(elements) != expectedElements {
		t.Fatalf(
			"expected %d elements, got %d",
			expectedElements,
			len(elements),
		)
	}
}

func Test_marshalJSONArray_UnsupportedValue(t *testing.T) {
	t.Parallel()

	unsupported := make(chan int)
	_, err := marshalJSONArray(unsupported)

	if !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("expected ErrInvalidMessage, got %v", err)
	}
}
