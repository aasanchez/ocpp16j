package ocpp16json_test

import (
	"encoding/json"
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

func ExampleNewCallResult() {
	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	payload := map[string]string{
		"status": "Accepted",
	}

	result, err := ocpp16json.NewCallResult(
		uniqueId, payload,
	)
	if err != nil {
		fmt.Println(errPrefix, err)

		return
	}

	wireBytes, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		fmt.Println(errPrefix, marshalErr)

		return
	}

	fmt.Println(string(wireBytes))

	// Output:
	// [3,"19223201",{"status":"Accepted"}]
}
