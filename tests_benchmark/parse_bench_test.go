//go:build bench

package ocpp16json_bench

import (
	"encoding/json"
	"strings"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

// Sink variables prevent the compiler from optimizing
// away benchmark results.
var (
	sinkMessage ocpp16json.Message
	sinkBytes   []byte
	sinkErr     error
)

// --- Parse benchmarks ---

var (
	callMinimal = []byte(
		`[2,"1","A",{}]`,
	)
	callSmall = []byte(
		`[2,"19223201","BootNotification",` +
			`{"chargePointVendor":"VendorX",` +
			`"chargePointModel":"Model1"}]`,
	)
	callMedium = []byte(
		`[2,"19223201","StartTransaction",` +
			`{"connectorId":1,` +
			`"idTag":"RFID-ABC123",` +
			`"meterStart":0,` +
			`"timestamp":"2024-01-15T08:00:00Z"}]`,
	)
	callResultSmall = []byte(
		`[3,"19223201",{"status":"Accepted"}]`,
	)
	callErrorSmall = []byte(
		`[4,"19223201","NotImplemented",` +
			`"Unknown action",{}]`,
	)
)

func BenchmarkParse_Call_Minimal(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		sinkMessage, sinkErr = ocpp16json.Parse(
			callMinimal,
		)
	}
}

func BenchmarkParse_Call_Small(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		sinkMessage, sinkErr = ocpp16json.Parse(
			callSmall,
		)
	}
}

func BenchmarkParse_Call_Medium(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		sinkMessage, sinkErr = ocpp16json.Parse(
			callMedium,
		)
	}
}

func BenchmarkParse_Call_LargePayload(b *testing.B) {
	b.ReportAllocs()

	largeValue := strings.Repeat("x", 4096)
	largeCall := []byte(
		`[2,"19223201","DataTransfer",` +
			`{"vendorId":"` + largeValue + `"}]`,
	)

	b.ResetTimer()

	for range b.N {
		sinkMessage, sinkErr = ocpp16json.Parse(
			largeCall,
		)
	}
}

func BenchmarkParse_CallResult(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		sinkMessage, sinkErr = ocpp16json.Parse(
			callResultSmall,
		)
	}
}

func BenchmarkParse_CallError(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		sinkMessage, sinkErr = ocpp16json.Parse(
			callErrorSmall,
		)
	}
}

func BenchmarkParse_InvalidJSON(b *testing.B) {
	b.ReportAllocs()

	invalid := []byte(`not json at all`)

	for range b.N {
		sinkMessage, sinkErr = ocpp16json.Parse(invalid)
	}
}

func BenchmarkParse_Parallel_Call(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sinkMessage, sinkErr = ocpp16json.Parse(
				callSmall,
			)
		}
	})
}

// --- Marshal benchmarks ---

func BenchmarkMarshalJSON_Call(b *testing.B) {
	b.ReportAllocs()

	message, _ := ocpp16json.Parse(callSmall)

	call, _ := ocpp16json.AsCall(message)

	b.ResetTimer()

	for range b.N {
		sinkBytes, sinkErr = json.Marshal(call)
	}
}

func BenchmarkMarshalJSON_CallResult(b *testing.B) {
	b.ReportAllocs()

	message, _ := ocpp16json.Parse(callResultSmall)

	callResult, _ := ocpp16json.AsCallResult(message)

	b.ResetTimer()

	for range b.N {
		sinkBytes, sinkErr = json.Marshal(callResult)
	}
}

func BenchmarkMarshalJSON_CallError(b *testing.B) {
	b.ReportAllocs()

	message, _ := ocpp16json.Parse(callErrorSmall)

	callError, _ := ocpp16json.AsCallError(message)

	b.ResetTimer()

	for range b.N {
		sinkBytes, sinkErr = json.Marshal(callError)
	}
}

func BenchmarkMarshalJSON_Parallel_Call(b *testing.B) {
	b.ReportAllocs()

	message, _ := ocpp16json.Parse(callSmall)

	call, _ := ocpp16json.AsCall(message)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sinkBytes, sinkErr = json.Marshal(call)
		}
	})
}

// --- Round-trip benchmarks ---

func BenchmarkRoundTrip_Call(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		message, _ := ocpp16json.Parse(callSmall)
		sinkBytes, sinkErr = json.Marshal(message)
	}
}

func BenchmarkRoundTrip_CallError(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		message, _ := ocpp16json.Parse(callErrorSmall)
		sinkBytes, sinkErr = json.Marshal(message)
	}
}

func BenchmarkRoundTrip_Parallel_Call(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			message, _ := ocpp16json.Parse(callSmall)
			sinkBytes, sinkErr = json.Marshal(message)
		}
	})
}
