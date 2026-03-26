package ocpp16json_test

import (
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

func Test_ErrInvalidMessage_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrInvalidMessage == nil {
		t.Fatal("ErrInvalidMessage is nil")
	}
}

func Test_ErrUnsupportedMessageType_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrUnsupportedMessageType == nil {
		t.Fatal("ErrUnsupportedMessageType is nil")
	}
}

func Test_ErrInvalidMessageID_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrInvalidMessageID == nil {
		t.Fatal("ErrInvalidMessageID is nil")
	}
}

func Test_ErrInvalidAction_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrInvalidAction == nil {
		t.Fatal("ErrInvalidAction is nil")
	}
}

func Test_ErrPayloadRequired_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrPayloadRequired == nil {
		t.Fatal("ErrPayloadRequired is nil")
	}
}

func Test_ErrPayloadDecode_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrPayloadDecode == nil {
		t.Fatal("ErrPayloadDecode is nil")
	}
}

func Test_ErrErrorCodeRequired_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrErrorCodeRequired == nil {
		t.Fatal("ErrErrorCodeRequired is nil")
	}
}

func Test_ErrErrorDescriptionAbsent_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrErrorDescriptionAbsent == nil {
		t.Fatal("ErrErrorDescriptionAbsent is nil")
	}
}

func Test_ErrErrorDetailsInvalid_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrErrorDetailsInvalid == nil {
		t.Fatal("ErrErrorDetailsInvalid is nil")
	}
}

func Test_ErrActionAlreadyRegistered_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrActionAlreadyRegistered == nil {
		t.Fatal("ErrActionAlreadyRegistered is nil")
	}
}

func Test_ErrUnknownAction_NonNil(t *testing.T) {
	t.Parallel()

	if ocpp16json.ErrUnknownAction == nil {
		t.Fatal("ErrUnknownAction is nil")
	}
}

const nextIndex = 1

func Test_SentinelErrors_AreDistinct(t *testing.T) {
	t.Parallel()

	sentinels := []error{
		ocpp16json.ErrInvalidMessage,
		ocpp16json.ErrUnsupportedMessageType,
		ocpp16json.ErrInvalidMessageID,
		ocpp16json.ErrInvalidAction,
		ocpp16json.ErrPayloadRequired,
		ocpp16json.ErrPayloadDecode,
		ocpp16json.ErrErrorCodeRequired,
		ocpp16json.ErrErrorDescriptionAbsent,
		ocpp16json.ErrErrorDetailsInvalid,
		ocpp16json.ErrActionAlreadyRegistered,
		ocpp16json.ErrUnknownAction,
	}

	for outer := range sentinels {
		for inner := outer + nextIndex; inner < len(sentinels); inner++ {
			if errors.Is(sentinels[outer], sentinels[inner]) {
				t.Errorf(
					"sentinel %d and %d are not distinct",
					outer,
					inner,
				)
			}
		}
	}
}
