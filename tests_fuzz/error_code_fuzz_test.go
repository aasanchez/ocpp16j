//go:build fuzz

package ocpp16json_fuzz

import (
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
// and verifies it never panics. Valid codes must be
// accepted, invalid codes must be rejected.
func FuzzNewErrorCode(f *testing.F) {
	f.Add("NotImplemented")
	f.Add("GenericError")
	f.Add("")
	f.Add("MadeUpError")
	f.Add("notimplemented")

	f.Fuzz(func(t *testing.T, input string) {
		code, err := ocpp16json.NewErrorCode(input)
		valid := validCodes()

		if err != nil {
			if valid[input] {
				t.Fatalf(
					"rejected valid code %q: %v",
					input, err,
				)
			}

			return
		}

		if !valid[input] {
			t.Fatalf(
				"accepted invalid code %q", input,
			)
		}

		if code.String() != input {
			t.Fatalf(
				"round-trip failed: %q vs %q",
				input, code.String(),
			)
		}
	})
}
