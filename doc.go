// Package ocpp16json implements the OCPP-J 1.6 wire framing layer.
//
// The package is transport-focused: it models OCPP JSON frames, validates
// their envelope shape, and decodes payloads into message types from
// github.com/aasanchez/ocpp16messages. It intentionally excludes business
// logic, state machines, and WebSocket session handling.
package ocpp16json
