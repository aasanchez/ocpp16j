package ocpp16json_test

import (
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	validCall = `[2, "19223201", "BootNotification",` +
		` {"chargePointVendor": "VendorX"}]`
	validCallResult = `[3, "19223201",` +
		` {"status": "Accepted"}]`
	validCallError = `[4, "19223201",` +
		` "NotImplemented", "Unknown action", {}]`
	errFmtParseNilGot   = "expected nil error, got %v"
	errFmtParseExpGot   = "expected %v, got %v"
	errFmtParseExpected = "expected %s, got %v"
)

// --- Parse: valid messages ---

func Test_Parse_ValidCall(t *testing.T) {
	t.Parallel()

	message, err := ocpp16json.Parse([]byte(validCall))
	if err != nil {
		t.Fatalf(errFmtParseNilGot, err)
	}

	if !ocpp16json.IsCall(message) {
		t.Fatalf(
			errFmtParseExpected,
			"Call", message.MessageType(),
		)
	}
}

func Test_Parse_ValidCall_UniqueId(t *testing.T) {
	t.Parallel()

	message, err := ocpp16json.Parse([]byte(validCall))
	if err != nil {
		t.Fatalf(errFmtParseNilGot, err)
	}

	if message.MessageId() != testUniqueIdStr {
		t.Fatalf(
			errFmtStrExpGot,
			testUniqueIdStr, message.MessageId(),
		)
	}
}

func Test_Parse_ValidCall_Action(t *testing.T) {
	t.Parallel()

	message, err := ocpp16json.Parse([]byte(validCall))
	if err != nil {
		t.Fatalf(errFmtParseNilGot, err)
	}

	rawCall, castErr := ocpp16json.AsCall(message)
	if castErr != nil {
		t.Fatalf(errFmtParseNilGot, castErr)
	}

	expectedAction := "BootNotification"
	if rawCall.Action != expectedAction {
		t.Fatalf(
			errFmtStrExpGot,
			expectedAction, rawCall.Action,
		)
	}
}

func Test_Parse_ValidCallResult(t *testing.T) {
	t.Parallel()

	message, err := ocpp16json.Parse(
		[]byte(validCallResult),
	)
	if err != nil {
		t.Fatalf(errFmtParseNilGot, err)
	}

	if !ocpp16json.IsCallResult(message) {
		t.Fatalf(
			errFmtParseExpected,
			"CallResult", message.MessageType(),
		)
	}
}

func Test_Parse_ValidCallError(t *testing.T) {
	t.Parallel()

	message, err := ocpp16json.Parse(
		[]byte(validCallError),
	)
	if err != nil {
		t.Fatalf(errFmtParseNilGot, err)
	}

	if !ocpp16json.IsCallError(message) {
		t.Fatalf(
			errFmtParseExpected,
			"CallError", message.MessageType(),
		)
	}
}

// --- Parse: invalid JSON ---

func Test_Parse_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse([]byte(`not json`))
	if !errors.Is(err, ocpp16json.ErrInvalidMessage) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessage, err,
		)
	}
}

func Test_Parse_EmptyArray(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse([]byte(`[]`))
	if !errors.Is(err, ocpp16json.ErrInvalidMessage) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessage, err,
		)
	}
}

func Test_Parse_NotAnArray(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse([]byte(`{"type": 2}`))
	if !errors.Is(err, ocpp16json.ErrInvalidMessage) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessage, err,
		)
	}
}

// --- Parse: bad MessageTypeId ---

func Test_Parse_UnknownMessageType(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[9, "123", "Foo", {}]`),
	)
	if !errors.Is(
		err, ocpp16json.ErrUnsupportedMessageType,
	) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrUnsupportedMessageType, err,
		)
	}
}

func Test_Parse_StringMessageType(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`["two", "123", "Foo", {}]`),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidMessage) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessage, err,
		)
	}
}

// --- Parse Call: wrong element count ---

func Test_Parse_Call_TooFewElements(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[2, "123", "Foo"]`),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidMessage) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessage, err,
		)
	}
}

func Test_Parse_Call_TooManyElements(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[2, "123", "Foo", {}, "extra"]`),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidMessage) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessage, err,
		)
	}
}

// --- Parse Call: bad UniqueId ---

func Test_Parse_Call_EmptyUniqueId(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[2, "", "Foo", {}]`),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidMessageID) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessageID, err,
		)
	}
}

func Test_Parse_Call_UniqueIdTooLong(t *testing.T) {
	t.Parallel()

	longId := `"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"`
	input := `[2, ` + longId + `, "Foo", {}]`

	_, err := ocpp16json.Parse([]byte(input))
	if !errors.Is(err, ocpp16json.ErrInvalidMessageID) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessageID, err,
		)
	}
}

func Test_Parse_Call_UniqueIdNotString(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[2, 123, "Foo", {}]`),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidMessageID) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessageID, err,
		)
	}
}

// --- Parse Call: bad Action ---

func Test_Parse_Call_EmptyAction(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[2, "123", "", {}]`),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidAction) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidAction, err,
		)
	}
}

func Test_Parse_Call_ActionNotString(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[2, "123", 42, {}]`),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidAction) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidAction, err,
		)
	}
}

// --- Parse Call: payload variants ---

func Test_Parse_Call_NullPayload(t *testing.T) {
	t.Parallel()

	message, err := ocpp16json.Parse(
		[]byte(`[2, "123", "Foo", null]`),
	)
	if err != nil {
		t.Fatalf(errFmtParseNilGot, err)
	}

	if !ocpp16json.IsCall(message) {
		t.Fatal("expected Call message")
	}
}

func Test_Parse_Call_EmptyObjectPayload(t *testing.T) {
	t.Parallel()

	message, err := ocpp16json.Parse(
		[]byte(`[2, "123", "Foo", {}]`),
	)
	if err != nil {
		t.Fatalf(errFmtParseNilGot, err)
	}

	if !ocpp16json.IsCall(message) {
		t.Fatal("expected Call message")
	}
}

// --- Parse CallResult: wrong element count ---

func Test_Parse_CallResult_TooFewElements(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse([]byte(`[3, "123"]`))
	if !errors.Is(err, ocpp16json.ErrInvalidMessage) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessage, err,
		)
	}
}

// --- Parse CallResult: bad UniqueId ---

func Test_Parse_CallResult_EmptyUniqueId(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[3, "", {}]`),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidMessageID) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessageID, err,
		)
	}
}

// --- Parse CallError: wrong element count ---

func Test_Parse_CallError_TooFewElements(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[4, "123", "NotImplemented", "desc"]`),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidMessage) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessage, err,
		)
	}
}

// --- Parse CallError: bad ErrorCode ---

func Test_Parse_CallError_InvalidErrorCode(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(
			`[4, "123", "MadeUp", "desc", {}]`,
		),
	)
	if !errors.Is(err, ocpp16json.ErrErrorCodeRequired) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrErrorCodeRequired, err,
		)
	}
}

func Test_Parse_CallError_EmptyErrorCode(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(`[4, "123", "", "desc", {}]`),
	)
	if !errors.Is(err, ocpp16json.ErrErrorCodeRequired) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrErrorCodeRequired, err,
		)
	}
}

// --- Parse CallError: ErrorDescription ---

func Test_Parse_CallError_EmptyDescription_Valid(
	t *testing.T,
) {
	t.Parallel()

	message, err := ocpp16json.Parse(
		[]byte(
			`[4, "123", "GenericError", "", {}]`,
		),
	)
	if err != nil {
		t.Fatalf(errFmtParseNilGot, err)
	}

	if !ocpp16json.IsCallError(message) {
		t.Fatal("expected CallError message")
	}
}

func Test_Parse_CallError_DescriptionNotString(
	t *testing.T,
) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(
			`[4, "123", "GenericError", 42, {}]`,
		),
	)
	if !errors.Is(
		err, ocpp16json.ErrErrorDescriptionAbsent,
	) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrErrorDescriptionAbsent, err,
		)
	}
}

// --- Parse CallError: bad ErrorDetails ---

func Test_Parse_CallError_DetailsNotObject(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(
			`[4, "123", "GenericError", "desc", "bad"]`,
		),
	)
	if !errors.Is(
		err, ocpp16json.ErrErrorDetailsInvalid,
	) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrErrorDetailsInvalid, err,
		)
	}
}

// --- Parse CallError: bad UniqueId ---

func Test_Parse_CallError_EmptyUniqueId(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.Parse(
		[]byte(
			`[4, "", "GenericError", "desc", {}]`,
		),
	)
	if !errors.Is(err, ocpp16json.ErrInvalidMessageID) {
		t.Fatalf(
			errFmtParseExpGot,
			ocpp16json.ErrInvalidMessageID, err,
		)
	}
}
