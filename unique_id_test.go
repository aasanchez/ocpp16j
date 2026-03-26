package ocpp16json

import (
	"errors"
	"strings"
	"testing"
)

const testUniqueIdValue = "19223201"

func Test_NewUniqueId_Valid(t *testing.T) {
	t.Parallel()

	uniqueId, err := NewUniqueId(testUniqueIdValue)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if uniqueId.String() != testUniqueIdValue {
		t.Fatalf(
			"expected %q, got %q",
			testUniqueIdValue, uniqueId,
		)
	}
}

func Test_NewUniqueId_Empty(t *testing.T) {
	t.Parallel()

	_, err := NewUniqueId("")
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf(
			"expected ErrInvalidMessageID, got %v", err,
		)
	}
}

func Test_NewUniqueId_ExactlyMaxLength(t *testing.T) {
	t.Parallel()

	value := strings.Repeat("a", maxUniqueIdLength)

	uniqueId, err := NewUniqueId(value)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if uniqueId.String() != value {
		t.Fatalf(
			"expected %q, got %q", value, uniqueId,
		)
	}
}

func Test_NewUniqueId_ExceedsMaxLength(t *testing.T) {
	t.Parallel()

	value := strings.Repeat("a", maxUniqueIdLength+1)

	_, err := NewUniqueId(value)
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf(
			"expected ErrInvalidMessageID, got %v", err,
		)
	}
}
