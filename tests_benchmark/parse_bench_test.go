//go:build bench

package ocpp16json_bench

import (
	"encoding/json"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

var (
	callBytes = []byte(
		`[2,"19223201","BootNotification",` +
			`{"chargePointVendor":"VendorX",` +
			`"chargePointModel":"Model1"}]`,
	)
	callResultBytes = []byte(
		`[3,"19223201",{"status":"Accepted"}]`,
	)
	callErrorBytes = []byte(
		`[4,"19223201","NotImplemented",` +
			`"Unknown action",{}]`,
	)
)

func BenchmarkParse_Call(b *testing.B) {
	for range b.N {
		_, _ = ocpp16json.Parse(callBytes)
	}
}

func BenchmarkParse_CallResult(b *testing.B) {
	for range b.N {
		_, _ = ocpp16json.Parse(callResultBytes)
	}
}

func BenchmarkParse_CallError(b *testing.B) {
	for range b.N {
		_, _ = ocpp16json.Parse(callErrorBytes)
	}
}

func BenchmarkMarshalJSON_Call(b *testing.B) {
	message, _ := ocpp16json.Parse(callBytes)

	call, _ := ocpp16json.AsCall(message)

	b.ResetTimer()

	for range b.N {
		_, _ = json.Marshal(call)
	}
}

func BenchmarkMarshalJSON_CallResult(b *testing.B) {
	message, _ := ocpp16json.Parse(callResultBytes)

	callResult, _ := ocpp16json.AsCallResult(message)

	b.ResetTimer()

	for range b.N {
		_, _ = json.Marshal(callResult)
	}
}

func BenchmarkMarshalJSON_CallError(b *testing.B) {
	message, _ := ocpp16json.Parse(callErrorBytes)

	callError, _ := ocpp16json.AsCallError(message)

	b.ResetTimer()

	for range b.N {
		_, _ = json.Marshal(callError)
	}
}

func BenchmarkRoundTrip_Call(b *testing.B) {
	for range b.N {
		message, _ := ocpp16json.Parse(callBytes)

		_, _ = json.Marshal(message)
	}
}
