package ocpp16json

import (
	"testing"
)

func Test_errMessageNotCall_NonNil(t *testing.T) {
	t.Parallel()

	if errMessageNotCall == nil {
		t.Fatal("errMessageNotCall is nil")
	}
}

func Test_errMessageNotCall_Message(t *testing.T) {
	t.Parallel()

	expected := "message is not a Call"
	if errMessageNotCall.Error() != expected {
		t.Fatalf(
			"expected %q, got %q",
			expected,
			errMessageNotCall.Error(),
		)
	}
}
