//go:build bench

package ocpp16json_bench

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	registrySmallSize  = 5
	registryMediumSize = 20
	registryLargeSize  = 100
)

// Sink variables for registry benchmarks.
var (
	sinkAny       any
	sinkUniqueId  ocpp16json.UniqueId
	sinkErrorCode ocpp16json.ErrorCode
	sinkCall      ocpp16json.Call
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

func registryWithActions(
	count int,
) *ocpp16json.Registry {
	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(benchConstructor)

	for actionIndex := range count {
		actionName := fmt.Sprintf(
			"Action%d", actionIndex,
		)
		_ = registry.Register(actionName, decoder)
	}

	return registry
}

// --- Registry Decode: varying registry sizes ---

func BenchmarkRegistry_Decode_5Actions(b *testing.B) {
	b.ReportAllocs()

	registry := registryWithActions(registrySmallSize)

	payload := json.RawMessage(`{"name":"test"}`)

	b.ResetTimer()

	for range b.N {
		sinkAny, sinkErr = registry.Decode(
			"Action0", payload,
		)
	}
}

func BenchmarkRegistry_Decode_20Actions(b *testing.B) {
	b.ReportAllocs()

	registry := registryWithActions(registryMediumSize)

	payload := json.RawMessage(`{"name":"test"}`)

	b.ResetTimer()

	for range b.N {
		sinkAny, sinkErr = registry.Decode(
			"Action10", payload,
		)
	}
}

func BenchmarkRegistry_Decode_100Actions(b *testing.B) {
	b.ReportAllocs()

	registry := registryWithActions(registryLargeSize)

	payload := json.RawMessage(`{"name":"test"}`)

	b.ResetTimer()

	for range b.N {
		sinkAny, sinkErr = registry.Decode(
			"Action50", payload,
		)
	}
}

func BenchmarkRegistry_Decode_UnknownAction(b *testing.B) {
	b.ReportAllocs()

	registry := registryWithActions(registryMediumSize)

	payload := json.RawMessage(`{"name":"test"}`)

	b.ResetTimer()

	for range b.N {
		sinkAny, sinkErr = registry.Decode(
			"NonExistent", payload,
		)
	}
}

func BenchmarkRegistry_Decode_Parallel(b *testing.B) {
	b.ReportAllocs()

	registry := registryWithActions(registryMediumSize)

	payload := json.RawMessage(`{"name":"test"}`)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sinkAny, sinkErr = registry.Decode(
				"Action5", payload,
			)
		}
	})
}

// --- Domain type constructors ---

func BenchmarkNewUniqueId_Short(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		sinkUniqueId, sinkErr = ocpp16json.NewUniqueId(
			"19223201",
		)
	}
}

func BenchmarkNewUniqueId_MaxLength(b *testing.B) {
	b.ReportAllocs()

	maxId := "550e8400-e29b-41d4-a716-446655440000"

	for range b.N {
		sinkUniqueId, sinkErr = ocpp16json.NewUniqueId(
			maxId,
		)
	}
}

func BenchmarkNewUniqueId_Invalid(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		sinkUniqueId, sinkErr = ocpp16json.NewUniqueId("")
	}
}

func BenchmarkNewErrorCode_Valid(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		sinkErrorCode, sinkErr = ocpp16json.NewErrorCode(
			"GenericError",
		)
	}
}

func BenchmarkNewErrorCode_Invalid(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		sinkErrorCode, sinkErr = ocpp16json.NewErrorCode(
			"MadeUpError",
		)
	}
}

// --- Message constructors ---

func BenchmarkNewCall_SmallPayload(b *testing.B) {
	b.ReportAllocs()

	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	payload := map[string]string{"key": "value"}

	b.ResetTimer()

	for range b.N {
		sinkCall, sinkErr = ocpp16json.NewCall(
			uniqueId, "Authorize", payload,
		)
	}
}

func BenchmarkNewCall_LargePayload(b *testing.B) {
	b.ReportAllocs()

	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	payload := map[string]string{
		"chargePointVendor":       "VendorX",
		"chargePointModel":        "Model1",
		"chargePointSerialNumber": "SN-12345",
		"chargeBoxSerialNumber":   "CB-67890",
		"firmwareVersion":         "1.0.0",
		"iccid":                   "89012345678901",
		"imsi":                    "310260000000000",
		"meterType":               "AC",
		"meterSerialNumber":       "MT-99999",
	}

	b.ResetTimer()

	for range b.N {
		sinkCall, sinkErr = ocpp16json.NewCall(
			uniqueId, "BootNotification", payload,
		)
	}
}

func BenchmarkNewCallError(b *testing.B) {
	b.ReportAllocs()

	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	b.ResetTimer()

	for range b.N {
		_, sinkErr = ocpp16json.NewCallError(
			uniqueId,
			ocpp16json.NotImplemented,
			"Unknown action",
			map[string]any{},
		)
	}
}

// --- Full pipeline: parse → decode → respond ---

func BenchmarkFullPipeline_ParseDecodeRespond(
	b *testing.B,
) {
	b.ReportAllocs()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(benchConstructor)
	_ = registry.Register("TestAction", decoder)

	wire := []byte(
		`[2,"19223201","TestAction",{"name":"test"}]`,
	)

	b.ResetTimer()

	for range b.N {
		message, _ := ocpp16json.Parse(wire)

		call, _ := ocpp16json.AsCall(message)

		sinkAny, _ = registry.Decode(
			call.Action, call.Payload,
		)

		_, sinkErr = ocpp16json.NewCallResult(
			call.UniqueId, sinkAny,
		)
	}
}

func BenchmarkFullPipeline_Parallel(b *testing.B) {
	b.ReportAllocs()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(benchConstructor)
	_ = registry.Register("TestAction", decoder)

	wire := []byte(
		`[2,"19223201","TestAction",{"name":"test"}]`,
	)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			message, _ := ocpp16json.Parse(wire)

			call, _ := ocpp16json.AsCall(message)

			result, _ := registry.Decode(
				call.Action, call.Payload,
			)

			_, _ = ocpp16json.NewCallResult(
				call.UniqueId, result,
			)
		}
	})
}
