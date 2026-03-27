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
//  1. Receive a raw JSON message (from a WebSocket, message
//     broker, or any other source).
//  2. Parse it with [Parse] to detect the message type.
//  3. For a Call, use [AsCall] to extract the Call struct,
//     then [Registry.Decode] to validate the Payload.
//  4. Apply application logic outside this package.
//  5. Build the response with [NewCallResult] or [NewCallError]
//     and marshal it with [json.Marshal].
//
// # CallResult Context
//
// OCPP-J CallResult messages do not carry the Action on the
// wire:
//
//	[3, "<UniqueId>", {<Payload>}]
//
// The caller must track which Action a UniqueId corresponds
// to in order to decode the Payload correctly.
//
// # Response Encoding
//
// This package validates inbound Payloads by decoding them
// into typed message values, but outbound Payloads are still
// normal JSON. A typical pattern is: validate the outgoing
// response with ocpp16messages, then wrap it in a CallResult
// using [NewCallResult].
package ocpp16json
