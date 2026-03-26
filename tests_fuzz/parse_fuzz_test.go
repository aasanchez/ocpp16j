//go:build fuzz

package ocpp16json_fuzz

import (
	"encoding/json"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

// FuzzParse feeds random bytes to Parse and verifies it
// never panics. Valid messages are round-tripped through
// marshal and re-parsed to verify consistency.
func FuzzParse(f *testing.F) {
	// Seed corpus with valid messages from the spec.
	f.Add([]byte(
		`[2,"19223201","BootNotification",{}]`,
	))
	f.Add([]byte(
		`[3,"19223201",{"status":"Accepted"}]`,
	))
	f.Add([]byte(
		`[4,"19223201","NotImplemented","",{}]`,
	))
	f.Add([]byte(`[]`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(`[9,"id","Action",{}]`))

	f.Fuzz(func(t *testing.T, data []byte) {
		message, parseErr := ocpp16json.Parse(data)
		if parseErr != nil {
			return
		}

		// Valid parse — round-trip through marshal.
		wireBytes, marshalErr := json.Marshal(message)
		if marshalErr != nil {
			return
		}

		// Re-parse must also succeed.
		reparsed, reparseErr := ocpp16json.Parse(
			wireBytes,
		)
		if reparseErr != nil {
			t.Fatalf(
				"round-trip failed: %v", reparseErr,
			)
		}

		// Message type must match.
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
	})
}
