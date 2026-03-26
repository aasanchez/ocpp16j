package ocpp16json_test

import (
	"encoding/json"
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

func ExampleNewRawCall() {
	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	payload := map[string]string{
		"idTag": "RFID-001",
	}

	rawCall, err := ocpp16json.NewRawCall(
		uniqueId, "Authorize", payload,
	)
	if err != nil {
		fmt.Println(errPrefix, err)

		return
	}

	wireBytes, marshalErr := json.Marshal(rawCall)
	if marshalErr != nil {
		fmt.Println(errPrefix, marshalErr)

		return
	}

	fmt.Println(string(wireBytes))

	// Output:
	// [2,"19223201","Authorize",{"idTag":"RFID-001"}]
}
