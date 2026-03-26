package ocpp16json_test

import (
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const errPrefix = "error:"

func ExampleParse() {
	data := []byte(
		`[2, "19223201", "BootNotification",` +
			` {"chargePointVendor": "VendorX"}]`,
	)

	message, err := ocpp16json.Parse(data)
	if err != nil {
		fmt.Println(errPrefix, err)

		return
	}

	fmt.Println(message.MessageType())
	fmt.Println(message.MessageId())
	fmt.Println(ocpp16json.IsCall(message))

	// Output:
	// 2
	// 19223201
	// true
}

func ExampleParse_callResult() {
	data := []byte(
		`[3, "19223201",` +
			` {"status": "Accepted"}]`,
	)

	message, err := ocpp16json.Parse(data)
	if err != nil {
		fmt.Println(errPrefix, err)

		return
	}

	fmt.Println(message.MessageType())
	fmt.Println(ocpp16json.IsCallResult(message))

	// Output:
	// 3
	// true
}

func ExampleParse_callError() {
	data := []byte(
		`[4, "19223201", "NotImplemented",` +
			` "Unknown action", {}]`,
	)

	message, err := ocpp16json.Parse(data)
	if err != nil {
		fmt.Println(errPrefix, err)

		return
	}

	fmt.Println(message.MessageType())
	fmt.Println(ocpp16json.IsCallError(message))

	// Output:
	// 4
	// true
}
