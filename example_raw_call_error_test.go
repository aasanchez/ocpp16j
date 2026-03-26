package ocpp16json_test

import (
	"encoding/json"
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

func ExampleNewRawCallError() {
	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	rawCallError, err := ocpp16json.NewRawCallError(
		uniqueId,
		ocpp16json.NotImplemented,
		"Unknown action",
		map[string]any{},
	)
	if err != nil {
		fmt.Println(errPrefix, err)

		return
	}

	wireBytes, marshalErr := json.Marshal(rawCallError)
	if marshalErr != nil {
		fmt.Println(errPrefix, marshalErr)

		return
	}

	fmt.Println(string(wireBytes))

	// Output:
	// [4,"19223201","NotImplemented","Unknown action",{}]
}
