# ocpp16j

OCPP-J 1.6 RPC framework for Go — parse, validate, and marshal Call,
CallResult, and CallError messages per the OCPP-J specification.

[![Go Reference](https://pkg.go.dev/badge/github.com/aasanchez/ocpp16j.svg)](https://pkg.go.dev/github.com/aasanchez/ocpp16j)
[![Go Report Card](https://goreportcard.com/badge/github.com/aasanchez/ocpp16j)](https://goreportcard.com/report/github.com/aasanchez/ocpp16j)

## What is this?

This package implements the **message wrapper layer** defined in
[section 4 of the OCPP-J 1.6 specification](https://www.openchargealliance.org/).
It sits between raw WebSocket bytes and typed OCPP payloads:

```text
WebSocket bytes  -->  ocpp16j (this package)  -->  ocpp16messages
                      Parse / Marshal / Validate    Payload types
```

It owns the JSON array envelope — **not** the Payload contents. Payload
validation is delegated to
[ocpp16messages](https://github.com/aasanchez/ocpp16messages).

### What it does

- Parse raw bytes into typed `Call`, `CallResult`, or `CallError`
- Marshal messages back to canonical OCPP-J arrays
- Validate the envelope: MessageTypeId, UniqueId, Action, ErrorCode
- Decode Payloads through a thread-safe `Registry` + `JSONDecoder`
- Provide first-class Go types for every spec concept

### What it does not do

- Charge Point or Central System behavior
- WebSocket session management
- Action routing or protocol state machines
- Payload schema validation (that's `ocpp16messages`)

## Installation

```sh
go get github.com/aasanchez/ocpp16j
```

Requires Go 1.24 or later.

## Quick Start

### Parse an incoming message

```go
message, err := ocpp16json.Parse(rawBytes)
if err != nil {
    // handle error
}

if ocpp16json.IsCall(message) {
    call, _ := ocpp16json.AsCall(message)
    fmt.Println(call.Action)  // "BootNotification"
}
```

### Build and send a Call

```go
uniqueId, _ := ocpp16json.NewUniqueId("19223201")

call, _ := ocpp16json.NewCall(
    uniqueId, "Authorize",
    map[string]string{"idTag": "RFID-001"},
)

wireBytes, _ := json.Marshal(call)
// wireBytes: [2,"19223201","Authorize",{"idTag":"RFID-001"}]
```

### Decode Payloads with the Registry

```go
registry := ocpp16json.NewRegistry()

// Register a decoder using an ocpp16messages constructor.
decoder := ocpp16json.JSONDecoder(bootnotification.Req)
registry.Register("BootNotification", decoder)

// Parse and decode in two steps.
message, _ := ocpp16json.Parse(rawBytes)
call, _ := ocpp16json.AsCall(message)

result, err := registry.Decode(call.Action, call.Payload)
// result is a bootnotification.ReqMessage with validated fields
```

### Build a CallError response

```go
uniqueId, _ := ocpp16json.NewUniqueId("req-99")

callError, _ := ocpp16json.NewCallError(
    uniqueId,
    ocpp16json.NotImplemented,
    "Requested Action is not known by receiver",
    map[string]any{},
)

wireBytes, _ := json.Marshal(callError)
// wireBytes: [4,"req-99","NotImplemented","Requested Action...",{}]
```

## OCPP-J Message Structures

From section 4.2 of the specification:

```text
Call:       [<MessageTypeId>, "<UniqueId>", "<Action>", {<Payload>}]
CallResult: [<MessageTypeId>, "<UniqueId>", {<Payload>}]
CallError:  [<MessageTypeId>, "<UniqueId>", "<ErrorCode>",
             "<ErrorDescription>", {<ErrorDetails>}]
```

## Domain Types

Every spec concept with constraints is a first-class Go type:

| Type          | Underlying | Constraint        | Spec Reference |
|---------------|------------|-------------------|----------------|
| `MessageType` | `uint8`    | Values 2, 3, 4    | Table 2        |
| `UniqueId`    | `string`   | Max 36 characters | Table 3        |
| `ErrorCode`   | `string`   | 10 valid values   | Table 7        |

### ErrorCode Constants (Table 7)

```go
ocpp16json.NotImplemented
ocpp16json.NotSupported
ocpp16json.InternalError
ocpp16json.ProtocolError
ocpp16json.SecurityError
ocpp16json.FormationViolation
ocpp16json.PropertyConstraintViolation
ocpp16json.OccurenceConstraintViolation  // spec spelling
ocpp16json.TypeConstraintViolation
ocpp16json.GenericError
```

## Error Handling

All errors are sentinel values checked with `errors.Is()`:

```go
_, err := ocpp16json.Parse(badBytes)
if errors.Is(err, ocpp16json.ErrInvalidMessage) {
    // not valid JSON or wrong structure
}
if errors.Is(err, ocpp16json.ErrUnsupportedMessageType) {
    // MessageTypeId is not 2, 3, or 4
}
```

Available sentinels: `ErrInvalidMessage`, `ErrUnsupportedMessageType`,
`ErrInvalidMessageID`, `ErrInvalidAction`, `ErrPayloadRequired`,
`ErrPayloadDecode`, `ErrErrorCodeRequired`, `ErrErrorDescriptionAbsent`,
`ErrErrorDetailsInvalid`, `ErrActionAlreadyRegistered`, `ErrUnknownAction`.

## Terminology

This package follows the OCPP-J 1.6 specification terminology strictly.
Reading the code should feel like reading the spec itself. See the
[spec terminology mapping](CLAUDE.md#terminology) for the full glossary.

## Related Packages

- [ocpp16messages](https://github.com/aasanchez/ocpp16messages) —
  OCPP 1.6 message types with `Req()`/`Conf()` constructors and validation

## License

[MIT](LICENSE) - Copyright (c) 2026 Alexis Sanchez
