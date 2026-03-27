# CLAUDE.md

## Project Overview

OCPP-J 1.6 RPC framework for Go (`github.com/aasanchez/ocpp16j`).
Implements the message wrapper layer defined in section 4 of the
OCPP-J 1.6 specification: parsing, marshaling, and validating
OCPP-J messages (Call, CallResult, CallError). Delegates all
Payload validation to `github.com/aasanchez/ocpp16messages`.
This is a library with no binary build target.

## Terminology

This package follows OCPP-J 1.6 specification terminology strictly.
Reading the code should feel like reading the spec itself.

- **Message** — the JSON array (transport-agnostic, not "frame")
- **MessageType** / **MessageTypeId** — the integer 2, 3, or 4
- **UniqueId** — the string identifier matching request to response
- **Action** — the case-sensitive name of the remote procedure
- **Payload** — the JSON object with arguments or results
- **ErrorCode** — the string code in a CallError
- **ErrorDescription** — the human-readable description in CallError
- **ErrorDetails** — the JSON object with error details in CallError
- **Call** — message type 2 (request)
- **CallResult** — message type 3 (response)
- **CallError** — message type 4 (error response)
- **Charge Point** / **Central System** — the two OCPP actors
- **MessageTypeNumber** — the integer value (2, 3, or 4) in Table 2
- **messageId** — max 36 characters (Table 3), to allow for GUIDs
- **ErrorCode values** (Table 7, wire-format strings, exact spec spelling):
  NotImplemented, NotSupported, InternalError, ProtocolError,
  SecurityError, FormationViolation, PropertyConstraintViolation,
  OccurenceConstraintViolation (spec typo — missing second 'r'),
  TypeConstraintViolation, GenericError
- **ErrorDescription** — "should be filled in if possible, otherwise
  a clear empty string" (empty `""` is valid on the wire)
- **Payload** — allows both `null` and empty object `{}` on the wire

## Prerequisites

- Go >= 1.24
- Tools: golangci-lint, staticcheck, gci, gofumpt, golines

## Common Commands

```sh
go mod tidy                        # Resolve dependencies
go build ./...                     # Build all packages
go test -v ./...                   # Run all tests verbose
go test -race ./...                # Run with race detector
go test -cover ./...               # Show coverage percentage
go test -run TestSpecificName ./.. # Run a single test
go vet ./...                       # Static analysis
make lint                          # golangci-lint + go vet + staticcheck
make format                        # gci + gofumpt + golines + gofmt
make test                          # Unit and example tests with coverage
```

## Architecture

This package sits between raw JSON bytes and
`github.com/aasanchez/ocpp16messages`. It owns:

- **Message parsing**: `Parse([]byte) (Message, error)` — validates
  the OCPP-J message wrapper (MessageTypeId, UniqueId, Action,
  Payload) and returns a typed Call, CallResult, or CallError
- **Message marshaling**: `MarshalJSON()` on message types —
  serializes back to canonical OCPP-J arrays
- **Payload decoding**: `JSONDecoder` adapter bridging raw JSON to
  `ocpp16messages` constructors via a thread-safe `Registry`
- **Error vocabulary**: sentinel errors for every failure mode

### OCPP-J Message Structures (spec section 4.2)

```text
Call:       [<MessageTypeId>, "<UniqueId>", "<Action>", {<Payload>}]
CallResult: [<MessageTypeId>, "<UniqueId>", {<Payload>}]
CallError:  [<MessageTypeId>, "<UniqueId>", "<ErrorCode>",
             "<ErrorDescription>", {<ErrorDetails>}]
```

CallResult does not carry the Action on the wire — the caller must
provide it explicitly when decoding.

## Go Code Style

### General

- Line length: 80 characters max (enforced by revive)
- Cognitive complexity: max 7 per function (enforced by revive)
- Indentation: tabs for Go files (gofmt convention)
- Always run `make format` before committing

### Imports

- Managed by gci: stdlib first, then project modules, then third-party
- No unused imports
- Prefer full package names over aliases for readability
- Use short aliases only for name conflicts or to stay under 80 chars

### Naming

- Exported identifiers: PascalCase
- Acronyms stay uppercase (`ID` not `Id`) — except where revive
  var-naming allowlist permits (e.g., `Id` in OCPP field names)
- Constructors: `New` prefix (`NewRegistry`, `NewFoo`)
- Getters: no `Get` prefix — use `Value()` not `GetValue()`
- Variable names must be descriptive (enforced by varnamelen) — avoid
  single-letter names like `p`, `id`; use `purposeType`, `profileId`

### Error Handling

- Never panic in library code
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- Use sentinel errors from `errors.go` — check with `errors.Is()`
- Do not create helper functions that just wrap `fmt.Errorf`
- Accumulate multiple validation errors with `errors.Join()` when a
  constructor validates several fields

### Type Design

- **First-class domain types are mandatory.** Every spec concept that
  carries constraints or semantics MUST be its own Go type — not a
  bare primitive. This is the single most important design rule in
  the project. Developers using this library should work with spec
  concepts as types, not with raw strings or ints. Examples:
  - `MessageType uint8` — not a bare `uint8`
  - `UniqueId string` — not a bare `string` (max 36 chars, Table 3)
  - `ErrorCode string` — not a bare `string` (Table 7 vocabulary)
  - Future spec concepts follow the same rule
- Each domain type gets a validating constructor (`NewUniqueId`,
  `NewErrorCode`, etc.) that enforces the spec constraints. The
  zero value of a domain type is invalid by design — only the
  constructor produces valid instances.
- All constructors return `(T, error)` — no separate `Validate()` methods
- Value receivers and immutable fields for thread safety
- Use `json.RawMessage` for Payload fields decoded later
- Generics for decoder adapters: `JSONDecoder[Input, Output any]`

## Testing

### Organization

- **Same-package tests** (`package ocpp16json`): for testing unexported
  functions and internals
- **External tests** (`package ocpp16json_test` in `tests/`): black-box
  tests for the public API
- **Example tests** (`example_*_test.go`): executable documentation for
  public constructors and complex APIs — skip for simple getters

### Rules

- Write atomic, individual test functions — each tests ONE behavior
- Every test must call `t.Parallel()`
- Use named constants instead of magic numbers in tests
- Use descriptive variable names (not `p`, `id`)
- Fuzz tests go in `tests_fuzz/` with `//go:build fuzz` tag
- Race tests go in `tests_race/` with `//go:build race` tag

### Test Naming

- Unit tests: `Test_<Function>_<Case>` or `Test<Type>_<Method>_<Case>`
- Example tests: `Example<Function>` or `Example<Type>_<method>`
- Subtests: descriptive suffixes

## Linting

golangci-lint config is in `golangci.yml`. Key settings:

- All linters enabled except: `wsl`, `testpackage`, `godot`
- `wsl_v5` enabled instead of `wsl`
- `depguard`: only stdlib + `github.com/aasanchez/ocpp16j` +
  `github.com/aasanchez/ocpp16messages` allowed
- `exhaustruct`: all struct fields must be explicitly initialized
- Reports written to `reports/`

## Dependencies

- `github.com/aasanchez/ocpp16messages` (main branch) — OCPP 1.6
  message types with `Req()`/`Conf()` constructors and validation
- Standard library only beyond that — zero third-party dependencies
