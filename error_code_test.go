package ocpp16json

import (
	"errors"
	"testing"
)

func Test_NewErrorCode_Valid_NotImplemented(t *testing.T) {
	t.Parallel()

	code, err := NewErrorCode("NotImplemented")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if code != NotImplemented {
		t.Fatalf(
			"expected %q, got %q",
			NotImplemented, code,
		)
	}
}

func Test_NewErrorCode_Valid_GenericError(t *testing.T) {
	t.Parallel()

	code, err := NewErrorCode("GenericError")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if code != GenericError {
		t.Fatalf(
			"expected %q, got %q", GenericError, code,
		)
	}
}

func Test_NewErrorCode_Empty(t *testing.T) {
	t.Parallel()

	_, err := NewErrorCode("")
	if !errors.Is(err, ErrErrorCodeRequired) {
		t.Fatalf(
			"expected ErrErrorCodeRequired, got %v", err,
		)
	}
}

func Test_NewErrorCode_InvalidValue(t *testing.T) {
	t.Parallel()

	_, err := NewErrorCode("MadeUpError")
	if !errors.Is(err, ErrErrorCodeRequired) {
		t.Fatalf(
			"expected ErrErrorCodeRequired, got %v", err,
		)
	}
}
