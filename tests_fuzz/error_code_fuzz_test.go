//go:build fuzz

package ocpp16json_fuzz

import (
	"strings"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

// validCodes returns the 10 valid ErrorCode values from
// spec Table 7.
func validCodes() map[string]bool {
	return map[string]bool{
		"NotImplemented":               true,
		"NotSupported":                 true,
		"InternalError":                true,
		"ProtocolError":                true,
		"SecurityError":                true,
		"FormationViolation":           true,
		"PropertyConstraintViolation":  true,
		"OccurenceConstraintViolation": true,
		"TypeConstraintViolation":      true,
		"GenericError":                 true,
	}
}

// FuzzNewErrorCode feeds random strings to NewErrorCode
// and verifies: never panics, accepts only Table 7 values,
// rejects everything else, round-trips through String().
func FuzzNewErrorCode(f *testing.F) {
	// All 10 valid codes.
	f.Add("NotImplemented")
	f.Add("NotSupported")
	f.Add("InternalError")
	f.Add("ProtocolError")
	f.Add("SecurityError")
	f.Add("FormationViolation")
	f.Add("PropertyConstraintViolation")
	f.Add("OccurenceConstraintViolation")
	f.Add("TypeConstraintViolation")
	f.Add("GenericError")

	// Case variations — must all be rejected.
	f.Add("notimplemented")
	f.Add("NOTIMPLEMENTED")
	f.Add("notImplemented")
	f.Add("genericerror")
	f.Add("GENERICERROR")

	// Near-misses.
	f.Add("NotImplemente")
	f.Add("NotImplementedd")
	f.Add("OccurrenceConstraintViolation")
	f.Add("Generic Error")
	f.Add("Generic_Error")

	// Empty and garbage.
	f.Add("")
	f.Add(" ")
	f.Add("MadeUpError")
	f.Add("null")
	f.Add("{}")
	f.Add(strings.Repeat("x", 1000))

	f.Fuzz(func(t *testing.T, input string) {
		code, err := ocpp16json.NewErrorCode(input)
		valid := validCodes()

		if err != nil {
			// Must not reject a valid code.
			if valid[input] {
				t.Fatalf(
					"rejected valid code %q: %v",
					input, err,
				)
			}

			return
		}

		// Must not accept an invalid code.
		if !valid[input] {
			t.Fatalf(
				"accepted invalid code %q", input,
			)
		}

		// Round-trip through String().
		if code.String() != input {
			t.Fatalf(
				"round-trip failed: %q vs %q",
				input, code.String(),
			)
		}

		// Creating again must produce the same result.
		duplicate, dupErr := ocpp16json.NewErrorCode(
			code.String(),
		)
		if dupErr != nil {
			t.Fatalf(
				"duplicate creation failed: %v",
				dupErr,
			)
		}

		if code != duplicate {
			t.Fatalf(
				"duplicate mismatch: %q vs %q",
				code, duplicate,
			)
		}
	})
}
