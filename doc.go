// Package ocpp16json implements the transport layer for OCPP-J 1.6.
//
// It focuses on wire-format correctness:
//   - parsing OCPP JSON arrays into typed frame values
//   - validating frame shape and mandatory envelope fields
//   - marshaling frames back to canonical OCPP-J arrays
//   - decoding payloads through validating constructors from
//     github.com/aasanchez/ocpp16messages
//
// The package does not implement charge point behavior, CSMS behavior,
// WebSocket session management, action routing, or protocol state machines.
//
// # Typical Flow
//
// A common integration flow looks like this:
//  1. Receive a raw JSON frame from the transport.
//  2. Parse it with Parse, or decode it with Registry.DecodeCall.
//  3. Apply application logic outside this package.
//  4. Validate the outgoing payload with ocpp16messages.
//  5. Wrap the encoded payload in RawCallResult or CallError.
//
// # CALLRESULT Context
//
// OCPP-J CALLRESULT frames do not carry the action name on the wire:
//
//	[3, "<uniqueId>", { ...payload... }]
//
// Because of that protocol constraint, Registry.DecodeCallResult requires the
// caller to provide the related action explicitly.
//
// # Response Encoding
//
// This package validates inbound payloads by decoding them into typed message
// values, but outbound payloads are still normal JSON. A typical pattern is:
// validate the outgoing response with ocpp16messages, encode the wire payload
// you want to send, and then wrap it in RawCallResult.
package ocpp16json
