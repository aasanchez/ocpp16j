//go:build fuzz

package ocpp16json_fuzz

import (
	"encoding/json"
	"strings"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

// FuzzParse feeds random bytes to Parse and verifies it
// never panics. Valid messages are round-tripped through
// marshal and re-parsed to verify full consistency.
func FuzzParse(f *testing.F) {
	// Valid messages from the spec.
	f.Add([]byte(
		`[2,"19223201","BootNotification",{}]`,
	))
	f.Add([]byte(
		`[2,"19223201","BootNotification",` +
			`{"chargePointVendor":"VendorX",` +
			`"chargePointModel":"Model1"}]`,
	))
	f.Add([]byte(
		`[2,"abc","Authorize",{"idTag":"RFID"}]`,
	))
	f.Add([]byte(
		`[2,"1","Heartbeat",{}]`,
	))
	f.Add([]byte(
		`[2,"tx","StartTransaction",` +
			`{"connectorId":1,"idTag":"A",` +
			`"meterStart":0,` +
			`"timestamp":"2024-01-01T00:00:00Z"}]`,
	))
	f.Add([]byte(
		`[3,"19223201",{"status":"Accepted"}]`,
	))
	f.Add([]byte(
		`[3,"19223201",{}]`,
	))
	f.Add([]byte(
		`[3,"19223201",null]`,
	))
	f.Add([]byte(
		`[4,"19223201","NotImplemented","",{}]`,
	))
	f.Add([]byte(
		`[4,"id","GenericError","desc",` +
			`{"detail":"value"}]`,
	))
	f.Add([]byte(
		`[4,"id","FormationViolation",` +
			`"bad payload",{}]`,
	))

	// Edge cases and invalid inputs.
	f.Add([]byte(`[]`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(`42`))
	f.Add([]byte(`true`))
	f.Add([]byte(``))
	f.Add([]byte(`[9,"id","Action",{}]`))
	f.Add([]byte(`[2]`))
	f.Add([]byte(`[2,""]`))
	f.Add([]byte(`[2,"id","",{}]`))
	f.Add([]byte(`[2,"id"]`))
	f.Add([]byte(`[3]`))
	f.Add([]byte(`[4,"id","BadCode","",{}]`))
	f.Add([]byte(`[4,"id","GenericError","",42]`))
	f.Add([]byte(`[2,123,"Action",{}]`))
	f.Add([]byte(`[2,"` +
		strings.Repeat("x", 37) +
		`","Action",{}]`))

	f.Fuzz(func(t *testing.T, data []byte) {
		message, parseErr := ocpp16json.Parse(data)
		if parseErr != nil {
			return
		}

		// Valid parse — round-trip through marshal.
		wireBytes, marshalErr := json.Marshal(message)
		if marshalErr != nil {
			t.Fatalf("marshal failed: %v", marshalErr)
		}

		// Re-parse must succeed.
		reparsed, reparseErr := ocpp16json.Parse(
			wireBytes,
		)
		if reparseErr != nil {
			t.Fatalf(
				"round-trip parse failed: %v\n"+
					"original: %s\n"+
					"marshaled: %s",
				reparseErr,
				string(data),
				string(wireBytes),
			)
		}

		// MessageType must match.
		if message.MessageType() != reparsed.MessageType() {
			t.Fatalf(
				"type mismatch: %d vs %d",
				message.MessageType(),
				reparsed.MessageType(),
			)
		}

		// UniqueId must match.
		if message.MessageId() != reparsed.MessageId() {
			t.Fatalf(
				"id mismatch: %q vs %q",
				message.MessageId(),
				reparsed.MessageId(),
			)
		}

		// For Call messages, Action must also match.
		if ocpp16json.IsCall(message) {
			original, _ := ocpp16json.AsCall(message)
			roundTripped, _ := ocpp16json.AsCall(reparsed)

			if original.Action != roundTripped.Action {
				t.Fatalf(
					"action mismatch: %q vs %q",
					original.Action,
					roundTripped.Action,
				)
			}
		}

		// For CallError, ErrorCode must match.
		if ocpp16json.IsCallError(message) {
			original, _ := ocpp16json.AsCallError(message)
			roundTripped, _ := ocpp16json.AsCallError(
				reparsed,
			)

			if original.ErrorCode != roundTripped.ErrorCode {
				t.Fatalf(
					"error code mismatch: %q vs %q",
					original.ErrorCode,
					roundTripped.ErrorCode,
				)
			}

			if original.ErrorDescription !=
				roundTripped.ErrorDescription {
				t.Fatalf(
					"description mismatch: %q vs %q",
					original.ErrorDescription,
					roundTripped.ErrorDescription,
				)
			}
		}
	})
}
