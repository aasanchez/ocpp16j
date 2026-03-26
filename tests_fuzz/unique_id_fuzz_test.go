//go:build fuzz

package ocpp16json_fuzz

import (
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const maxUniqueIdLength = 36

// FuzzNewUniqueId feeds random strings to NewUniqueId and
// verifies it never panics. Valid IDs must round-trip
// through String().
func FuzzNewUniqueId(f *testing.F) {
	f.Add("19223201")
	f.Add("550e8400-e29b-41d4-a716-446655440000")
	f.Add("")
	f.Add("a")
	f.Add("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	f.Fuzz(func(t *testing.T, input string) {
		uniqueId, err := ocpp16json.NewUniqueId(input)
		if err != nil {
			// Must reject empty or >36 chars.
			if input != "" && len(input) <= maxUniqueIdLength {
				t.Fatalf(
					"rejected valid input %q: %v",
					input, err,
				)
			}

			return
		}

		// Must accept non-empty, <=36 chars.
		if input == "" || len(input) > maxUniqueIdLength {
			t.Fatalf(
				"accepted invalid input %q", input,
			)
		}

		// Round-trip.
		if uniqueId.String() != input {
			t.Fatalf(
				"round-trip failed: %q vs %q",
				input, uniqueId.String(),
			)
		}
	})
}
