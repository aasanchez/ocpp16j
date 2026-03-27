//go:build fuzz

package ocpp16json_fuzz

import (
	"strings"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const maxUniqueIdLength = 36

// FuzzNewUniqueId feeds random strings to NewUniqueId and
// verifies: never panics, accepts only valid inputs, rejects
// invalid inputs, and round-trips through String().
func FuzzNewUniqueId(f *testing.F) {
	// Valid seeds.
	f.Add("19223201")
	f.Add("550e8400-e29b-41d4-a716-446655440000")
	f.Add("a")
	f.Add("1")
	f.Add(strings.Repeat("x", maxUniqueIdLength))

	// Boundary seeds.
	f.Add("")
	f.Add(strings.Repeat("x", maxUniqueIdLength+1))
	f.Add(strings.Repeat("x", maxUniqueIdLength-1))

	// Special characters.
	f.Add("id with spaces")
	f.Add("id\twith\ttabs")
	f.Add("id\nwith\nnewlines")
	f.Add("ñoño")
	f.Add("日本語")
	f.Add("\x00\x01\x02")
	f.Add("null")
	f.Add("{}")

	f.Fuzz(func(t *testing.T, input string) {
		uniqueId, err := ocpp16json.NewUniqueId(input)

		shouldBeValid := input != "" &&
			len(input) <= maxUniqueIdLength

		if err != nil && shouldBeValid {
			t.Fatalf(
				"rejected valid input %q (len=%d): %v",
				input, len(input), err,
			)
		}

		if err == nil && !shouldBeValid {
			t.Fatalf(
				"accepted invalid input %q (len=%d)",
				input, len(input),
			)
		}

		if err != nil {
			return
		}

		// Round-trip through String().
		if uniqueId.String() != input {
			t.Fatalf(
				"round-trip failed: %q vs %q",
				input, uniqueId.String(),
			)
		}

		// Creating again with the same value must
		// produce the same result.
		duplicate, dupErr := ocpp16json.NewUniqueId(
			uniqueId.String(),
		)
		if dupErr != nil {
			t.Fatalf(
				"duplicate creation failed: %v", dupErr,
			)
		}

		if uniqueId != duplicate {
			t.Fatalf(
				"duplicate mismatch: %q vs %q",
				uniqueId, duplicate,
			)
		}
	})
}
