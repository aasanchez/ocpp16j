package ocpp16json

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/aasanchez/ocpp16messages/authorize"
)

func TestJSONDecoderBuildsValidatedPayload(t *testing.T) {
	t.Parallel()

	decoder := JSONDecoder(authorize.Req)

	value, err := decoder(json.RawMessage(`{"idTag":"RFID-123"}`))
	if err != nil {
		t.Fatalf("decoder: %v", err)
	}

	payload, ok := value.(authorize.ReqMessage)
	if !ok {
		t.Fatalf("unexpected payload type: %T", value)
	}

	if payload.IdTag.String() != "RFID-123" {
		t.Fatalf("unexpected idTag: %q", payload.IdTag.String())
	}
}

func TestJSONDecoderRejectsInvalidJSONType(t *testing.T) {
	t.Parallel()

	decoder := JSONDecoder(authorize.Req)

	_, err := decoder(json.RawMessage(`{"idTag":1}`))
	if !errors.Is(err, ErrPayloadDecode) {
		t.Fatalf("expected ErrPayloadDecode, got %v", err)
	}
}
