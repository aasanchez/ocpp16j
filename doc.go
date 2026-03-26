// Package ocpp16json implements the OCPP-J 1.6 RPC framework.
//
// It focuses on wire-format correctness for the three OCPP-J
// message types defined in section 4 of the OCPP-J 1.6
// specification:
//   - parsing OCPP-J messages (Call, CallResult, CallError)
//   - validating the message wrapper: MessageTypeId, UniqueId,
//     Action, Payload, ErrorCode, ErrorDescription, ErrorDetails
//   - marshaling messages back to canonical OCPP-J arrays
//   - decoding Payload through validating constructors from
//     github.com/aasanchez/ocpp16messages
//
// The package does not implement Charge Point behavior, Central
// System behavior, WebSocket session management, action routing,
// or protocol state machines.
//
// # Typical Flow
//
// A common integration flow looks like this:
//  1. Receive a raw JSON message from the WebSocket transport.
//  2. Parse it with Parse, or decode it with Registry.DecodeCall.
//  3. Apply application logic outside this package.
//  4. Validate the outgoing Payload with ocpp16messages.
//  5. Wrap the encoded Payload in a CallResult or CallError.
//
// # CallResult Context
//
// OCPP-J CallResult messages do not carry the Action on the wire:
//
//	[3, "<UniqueId>", {<Payload>}]
//
// Because of that protocol constraint, Registry.DecodeCallResult
// requires the caller to provide the related Action explicitly.
//
// # Response Encoding
//
// This package validates inbound Payloads by decoding them into
// typed message values, but outbound Payloads are still normal
// JSON. A typical pattern is: validate the outgoing response with
// ocpp16messages, encode the wire Payload, and wrap it in a
// CallResult.
package ocpp16json
