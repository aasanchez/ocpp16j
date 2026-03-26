package ocpp16json_test

import (
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

func ExampleNewUniqueId() {
	uniqueId, err := ocpp16json.NewUniqueId(
		"550e8400-e29b-41d4-a716-446655440000",
	)
	if err != nil {
		fmt.Println(errPrefix, err)

		return
	}

	fmt.Println(uniqueId)

	// Output:
	// 550e8400-e29b-41d4-a716-446655440000
}
