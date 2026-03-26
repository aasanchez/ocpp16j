package ocpp16json_test

import (
	"errors"
	"strings"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	maxUniqueIdLength    = 36
	shortUniqueIdValue   = "19223201"
	validGUID            = "550e8400-e29b-41d4-a716-446655440000"
	errFmtExpMsgID       = "expected ErrInvalidMessageID, got %v"
	errFmtExpNilUniqueId = "expected nil error, got %v"
)

func Test_NewUniqueId_ValidShort(t *testing.T) {
	t.Parallel()

	uniqueId, err := ocpp16json.NewUniqueId(shortUniqueIdValue)
	if err != nil {
		t.Fatalf(errFmtExpNilUniqueId, err)
	}

	if uniqueId.String() != shortUniqueIdValue {
		t.Fatalf(
			"expected %q, got %q",
			shortUniqueIdValue, uniqueId.String(),
		)
	}
}

func Test_NewUniqueId_ValidGUID(t *testing.T) {
	t.Parallel()

	uniqueId, err := ocpp16json.NewUniqueId(validGUID)
	if err != nil {
		t.Fatalf(errFmtExpNilUniqueId, err)
	}

	if uniqueId.String() != validGUID {
		t.Fatalf(
			"expected %q, got %q",
			validGUID, uniqueId.String(),
		)
	}
}

func Test_NewUniqueId_ExactlyMaxLength(t *testing.T) {
	t.Parallel()

	value := strings.Repeat("x", maxUniqueIdLength)

	_, err := ocpp16json.NewUniqueId(value)
	if err != nil {
		t.Fatalf(errFmtExpNilUniqueId, err)
	}
}

func Test_NewUniqueId_Empty(t *testing.T) {
	t.Parallel()

	_, err := ocpp16json.NewUniqueId("")
	if !errors.Is(err, ocpp16json.ErrInvalidMessageID) {
		t.Fatalf(errFmtExpMsgID, err)
	}
}

func Test_NewUniqueId_ExceedsMaxLength(t *testing.T) {
	t.Parallel()

	value := strings.Repeat("x", maxUniqueIdLength+1)

	_, err := ocpp16json.NewUniqueId(value)
	if !errors.Is(err, ocpp16json.ErrInvalidMessageID) {
		t.Fatalf(errFmtExpMsgID, err)
	}
}
