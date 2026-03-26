package ocpp16json_test

import (
	"encoding/json"
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	testRegistryAction = "Authorize"
	errFmtRegNilGot    = "expected nil error, got %v"
	errFmtRegExpGot    = "expected %v, got %v"
	errFmtRegNilResult = "expected non-nil result, got nil"
)

// testPayload is a simple struct used to test JSONDecoder.
type testPayload struct {
	Name string `json:"name"`
}

// errNameRequired is a test-only sentinel error.
var errNameRequired = errors.New("name is required")

func testConstructor(
	input testPayload,
) (testPayload, error) {
	if input.Name == "" {
		return testPayload{}, errNameRequired
	}

	return input, nil
}

// --- Registry: Register ---

func Test_Registry_Register_Success(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	decoder := ocpp16json.JSONDecoder(testConstructor)

	err := registry.Register(testRegistryAction, decoder)
	if err != nil {
		t.Fatalf(errFmtRegNilGot, err)
	}
}

func Test_Registry_Register_Duplicate(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	decoder := ocpp16json.JSONDecoder(testConstructor)

	_ = registry.Register(testRegistryAction, decoder)

	err := registry.Register(testRegistryAction, decoder)
	if !errors.Is(
		err, ocpp16json.ErrActionAlreadyRegistered,
	) {
		t.Fatalf(
			errFmtRegExpGot,
			ocpp16json.ErrActionAlreadyRegistered, err,
		)
	}
}

// --- Registry: Decode ---

func Test_Registry_Decode_Success(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	decoder := ocpp16json.JSONDecoder(testConstructor)

	_ = registry.Register(testRegistryAction, decoder)

	payload := json.RawMessage(`{"name": "test"}`)

	result, err := registry.Decode(
		testRegistryAction, payload,
	)
	if err != nil {
		t.Fatalf(errFmtRegNilGot, err)
	}

	if result == nil {
		t.Fatal(errFmtRegNilResult)
	}

	output, isTestPayload := result.(testPayload)
	if !isTestPayload {
		t.Fatalf("expected testPayload, got %T", result)
	}

	expectedName := "test"
	if output.Name != expectedName {
		t.Fatalf(
			errFmtStrExpGot, expectedName, output.Name,
		)
	}
}

func Test_Registry_Decode_UnknownAction(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	payload := json.RawMessage(`{}`)

	_, err := registry.Decode("Unknown", payload)
	if !errors.Is(err, ocpp16json.ErrUnknownAction) {
		t.Fatalf(
			errFmtRegExpGot,
			ocpp16json.ErrUnknownAction, err,
		)
	}
}

// --- JSONDecoder: error paths ---

func Test_JSONDecoder_InvalidJSON(t *testing.T) {
	t.Parallel()

	decoder := ocpp16json.JSONDecoder(testConstructor)
	payload := json.RawMessage(`not json`)

	_, err := decoder(payload)
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf(
			errFmtRegExpGot,
			ocpp16json.ErrPayloadDecode, err,
		)
	}
}

func Test_JSONDecoder_ConstructorError(t *testing.T) {
	t.Parallel()

	decoder := ocpp16json.JSONDecoder(testConstructor)
	payload := json.RawMessage(`{"name": ""}`)

	_, err := decoder(payload)
	if !errors.Is(err, ocpp16json.ErrPayloadDecode) {
		t.Fatalf(
			errFmtRegExpGot,
			ocpp16json.ErrPayloadDecode, err,
		)
	}
}
