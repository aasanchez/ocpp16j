# OCPP 1.6 JSON for Go

`github.com/aasanchez/ocpp16j` is a strict transport-layer implementation of
the OCPP-J 1.6 wire format for Go.

It is intentionally narrow:

- parse and validate OCPP JSON frames
- marshal frames back to the canonical array shape
- decode payloads through validated constructors from
  [`github.com/aasanchez/ocpp16messages`](https://github.com/aasanchez/ocpp16messages)
- expose stable sentinel errors for envelope failures

It intentionally does not implement:

- charge point behavior
- CSMS behavior
- WebSocket clients or servers
- retry, routing, or session state machines

This separation keeps the package focused on transport correctness.

## Installation

```bash
go get github.com/aasanchez/ocpp16j
```

## Quick Start

```go
package main

import (
	"fmt"

	ocpp16json "github.com/aasanchez/ocpp16j"
	"github.com/aasanchez/ocpp16messages/authorize"
)

func main() {
	registry := ocpp16json.NewRegistry()

	err := registry.RegisterRequest(
		"Authorize",
		ocpp16json.JSONDecoder(authorize.Req),
	)
	if err != nil {
		panic(err)
	}

	decoded, err := registry.DecodeCall(
		[]byte(`[2,"uid-1","Authorize",{"idTag":"RFID-123"}]`),
	)
	if err != nil {
		panic(err)
	}

	req := decoded.Payload.(authorize.ReqMessage)

	fmt.Println(decoded.Action)
	fmt.Println(req.IdTag.String())
}
```

Expected output:

```text
Authorize
RFID-123
```

## What The Package Handles

### Inbound frames

- `Parse` validates raw OCPP-J arrays and returns `RawCall`,
  `RawCallResult`, or `CallError`
- `Registry.DecodeCall` decodes request payloads using the action name found
  in the frame
- `Registry.DecodeCallResult` decodes confirmation payloads using the action
  name provided by the caller

### Outbound frames

- `RawCall`, `RawCallResult`, and `CallError` marshal back to OCPP-J arrays
- payload validation stays in `ocpp16messages`
- transport wrapping stays in this package

That split lets you validate payload semantics without mixing in transport or
session policy.

## End-To-End Flow

The usual request/response flow is:

1. Receive a raw OCPP-J JSON frame.
2. Decode it with `Registry.DecodeCall(...)`.
3. Work with the validated payload returned by `ocpp16messages`.
4. Run your application logic.
5. Validate the outgoing payload with `ocpp16messages`.
6. Encode the response payload JSON and wrap it in `RawCallResult`.

The package examples on pkgsite show this flow end to end, including:

- request decode and validation
- response construction and wrapping
- raw frame parsing
- `CALLRESULT` decoding
- `CALLERROR` inspection

Pkgsite:

- [github.com/aasanchez/ocpp16j](https://pkg.go.dev/github.com/aasanchez/ocpp16j)

## Why `CALLRESULT` Needs Action Context

In OCPP-J 1.6, a `CALLRESULT` frame contains only:

```json
[3, "<uniqueId>", { ...payload... }]
```

The action name is not present on the wire. That is a protocol constraint, so
this package keeps it explicit instead of hiding it behind guessed state:

- `Parse` returns the raw frame
- `Registry.DecodeCall` uses the action embedded in `CALL`
- `Registry.DecodeCallResult` requires the related action from the caller

## Error Model

Envelope-level failures use stable sentinel errors such as:

- `ErrInvalidFrame`
- `ErrInvalidAction`
- `ErrInvalidMessageID`
- `ErrPayloadRequired`
- `ErrPayloadDecode`
- `ErrUnknownAction`

Payload validation errors from `ocpp16messages` are wrapped, so `errors.Is`
continues to work for both transport and domain validation checks.

## Design Rules

- transport correctness first
- no business logic in this module
- no hidden state for response decoding
- validate payloads with `ocpp16messages`, not ad-hoc transport structs
- keep examples and tests executable

## Status

Current coverage includes:

- parsing and marshaling for `CALL`, `CALLRESULT`, and `CALLERROR`
- typed request and confirmation decoding through a registry
- package examples ready for pkgsite
- fuzz coverage for malformed frames and payload decode paths
- race-detector coverage for concurrent registry and parser use

Planned next steps:

- prebuilt registries for the full OCPP 1.6 action set
- broader confirmation/request helpers for consumers
