package ocpp16json

import (
	"testing"
)

func Test_errFrameNotRawCall_NonNil(t *testing.T) {
	t.Parallel()

	if errFrameNotRawCall == nil {
		t.Fatal("errFrameNotRawCall is nil")
	}
}

func Test_errFrameNotRawCall_Message(t *testing.T) {
	t.Parallel()

	expected := "frame is not a raw call"
	if errFrameNotRawCall.Error() != expected {
		t.Fatalf(
			"expected %q, got %q",
			expected,
			errFrameNotRawCall.Error(),
		)
	}
}
