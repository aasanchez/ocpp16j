package ocpp16json_test

import (
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

func ExampleNewErrorCode() {
	code, err := ocpp16json.NewErrorCode("GenericError")
	if err != nil {
		fmt.Println(errPrefix, err)

		return
	}

	fmt.Println(code)

	// Output:
	// GenericError
}
