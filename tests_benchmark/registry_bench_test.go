//go:build bench

package ocpp16json_bench

import (
	"encoding/json"
	"errors"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

type benchPayload struct {
	Name string `json:"name"`
}

var errBenchName = errors.New("name required")

func benchConstructor(
	input benchPayload,
) (benchPayload, error) {
	if input.Name == "" {
		return benchPayload{}, errBenchName
	}

	return input, nil
}

func BenchmarkRegistry_Decode(b *testing.B) {
	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(benchConstructor)
	_ = registry.Register("TestAction", decoder)

	payload := json.RawMessage(`{"name":"test"}`)

	b.ResetTimer()

	for range b.N {
		_, _ = registry.Decode("TestAction", payload)
	}
}

func BenchmarkNewUniqueId(b *testing.B) {
	for range b.N {
		_, _ = ocpp16json.NewUniqueId("19223201")
	}
}

func BenchmarkNewErrorCode(b *testing.B) {
	for range b.N {
		_, _ = ocpp16json.NewErrorCode("GenericError")
	}
}

func BenchmarkNewCall(b *testing.B) {
	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	payload := map[string]string{"key": "value"}

	b.ResetTimer()

	for range b.N {
		_, _ = ocpp16json.NewCall(
			uniqueId, "Authorize", payload,
		)
	}
}

func BenchmarkNewCallError(b *testing.B) {
	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	b.ResetTimer()

	for range b.N {
		_, _ = ocpp16json.NewCallError(
			uniqueId,
			ocpp16json.NotImplemented,
			"Unknown action",
			map[string]any{},
		)
	}
}
