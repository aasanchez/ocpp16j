# OCPP 1.6 JSON for Go

A strict Go implementation of the OCPP-J 1.6 wire format. This module focuses
on JSON frame correctness and typed payload decoding on top of
[`github.com/aasanchez/ocpp16messages`](https://github.com/aasanchez/ocpp16messages).
It intentionally does not implement charge point behavior, central system
behavior, WebSocket reconnection, or any protocol state machine.

## Scope

This repository is the transport layer for OCPP 1.6 JSON:

- Parse and validate OCPP-J frames (`CALL`, `CALLRESULT`, `CALLERROR`)
- Marshal frames back to the canonical JSON array form
- Decode payloads into validated `ocpp16messages` request/confirmation types
- Expose stable sentinel errors for envelope-level failures

Out of scope:

- Business logic and profile behavior
- Session correlation and command routing policy
- WebSocket client/server implementation

The design follows the same principles as `ocpp16messages`: strict validation,
small focused packages, explicit errors, and tests centered on wire
correctness.

## Why `CALLRESULT` Decoding Needs Context

In OCPP-J 1.6, a `CALLRESULT` frame contains only:

```json
[3, "<uniqueId>", { ...payload... }]
```

The action name is not present on the wire, so typed response decoding is not
fully stateless. This package makes that explicit:

- `Parse` returns a raw frame for any valid OCPP-J message
- `Registry.DecodeCall` decodes `CALL` payloads using the action in the frame
- `Registry.DecodeCallResult` requires the caller to provide the related action

That keeps the transport layer correct without inventing session behavior.

## Installation

```bash
go get github.com/aasanchez/ocpp16json
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/aasanchez/ocpp16json"
	"github.com/aasanchez/ocpp16messages/authorize"
)

func main() {
	registry := ocpp16json.NewRegistry()
	_ = registry.RegisterRequest("Authorize", ocpp16json.JSONDecoder(authorize.Req))

	frame, err := registry.DecodeCall(
		[]byte(`[2,"uid-1","Authorize",{"idTag":"RFID-123"}]`),
	)
	if err != nil {
		panic(err)
	}

	req := frame.Payload.(authorize.ReqMessage)
	fmt.Println(frame.Action, req.IdTag.String())
}
```

## Status

Initial transport skeleton. The current implementation covers:

- Envelope parsing and validation
- JSON marshaling for all three OCPP-J frame types
- Registry-based typed decoding for request and confirmation payloads
- Tests against real payload types from `ocpp16messages`

Next steps:

- Add registry helpers for the full OCPP 1.6 action set
- Add fuzz tests for malformed frame envelopes
- Add CI, coverage, and compatibility checks matching the upstream repository
