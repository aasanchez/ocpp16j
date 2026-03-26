package ocpp16json_test

import (
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	errFmtExpCodeReq = "expected ErrErrorCodeRequired, got %v"
	errFmtExpNilCode = "expected nil error, got %v"
	errFmtCodeExpGot = "expected %q, got %q"
	totalErrorCodes  = 10
)

func allErrorCodes() []ocpp16json.ErrorCode {
	return []ocpp16json.ErrorCode{
		ocpp16json.NotImplemented,
		ocpp16json.NotSupported,
		ocpp16json.InternalError,
		ocpp16json.ProtocolError,
		ocpp16json.SecurityError,
		ocpp16json.FormationViolation,
		ocpp16json.PropertyConstraintViolation,
		ocpp16json.OccurenceConstraintViolation,
		ocpp16json.TypeConstraintViolation,
		ocpp16json.GenericError,
	}
}

func Test_ErrorCode_AllTenConstantsExist(t *testing.T) {
	t.Parallel()

	codes := allErrorCodes()
	if len(codes) != totalErrorCodes {
		t.Fatalf(
			"expected %d error codes, got %d",
			totalErrorCodes, len(codes),
		)
	}
}

func Test_ErrorCode_AllConstantsAreDistinct(t *testing.T) {
	t.Parallel()

	seen := make(map[ocpp16json.ErrorCode]bool)

	for _, code := range allErrorCodes() {
		if seen[code] {
			t.Fatalf("duplicate error code: %q", code)
		}

		seen[code] = true
	}
}

func Test_NewErrorCode_AllValidCodesAccepted(t *testing.T) {
	t.Parallel()

	for _, expected := range allErrorCodes() {
		code, err := ocpp16json.NewErrorCode(
			expected.String(),
		)
		if err != nil {
			t.Fatalf(
				errFmtExpNilCode+" for %q", err, expected,
			)
		}

		if code != expected {
			t.Fatalf(
				errFmtCodeExpGot, expected, code,
			)
		}
	}
}

func Test_NewErrorCode_Empty(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.NewErrorCode("")
	if !errors.Is(err, ocpp16json.ErrErrorCodeRequired) {
		t.Fatalf(errFmtExpCodeReq, err)
	}
}

func Test_NewErrorCode_InvalidValue(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.NewErrorCode("MadeUpError")
	if !errors.Is(err, ocpp16json.ErrErrorCodeRequired) {
		t.Fatalf(errFmtExpCodeReq, err)
	}
}

func Test_ErrorCode_OccurenceSpelling(t *testing.T) {
	t.Parallel()

	code := ocpp16json.OccurenceConstraintViolation
	expected := "OccurenceConstraintViolation"

	if code.String() != expected {
		t.Fatalf(
			errFmtCodeExpGot, expected, code.String(),
		)
	}
}
