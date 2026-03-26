# OCPP-J 1.6 JSON Wrapper — Implementation Plan

## Context

Complete the OCPP-J 1.6 JSON message wrapper library. This package owns
**only** the wire-format layer: detect message type, validate the JSON
envelope, marshal/unmarshal OCPP-J arrays, and produce proper errors.
All payload logic is delegated to `github.com/aasanchez/ocpp16messages`.

## Terminology Audit (spec vs code)

Reviewed against OCPP-J 1.6 specification sections 4.1.3, 4.1.4, 4.2.

### Correct — matches spec exactly

| Spec Term (section)                     | Code Identifier                             | Notes                                 |
|-----------------------------------------|---------------------------------------------|---------------------------------------|
| MessageType (Table 2 header)            | `MessageType` type                          | Go PascalCase of spec term            |
| MessageTypeNumber (Table 2 header)      | used in error comments                      | Exact spec term                       |
| CALL / CALLRESULT / CALLERROR (Table 2) | `Call`, `CallResult`, `CallError` constants | Go-style casing, values 2/3/4 correct |
| UniqueId (Tables 4-6)                   | `UniqueId` struct field                     | Exact match                           |
| Action (Table 4)                        | `Action` struct field                       | Exact match                           |
| Payload (Tables 4-5)                    | `Payload json.RawMessage`                   | Exact match                           |
| ErrorCode (Table 6)                     | `ErrorCode` struct field                    | Exact match                           |
| ErrorDescription (Table 6)              | `ErrorDescription` struct field             | Exact match                           |
| ErrorDetails (Table 6)                  | `ErrorDetails map[string]any`               | Exact match                           |
| MessageTypeId (section 4.2.x)           | used in comments                            | Exact spec term                       |

### Gaps — missing from code

| Spec Requirement                                | Section    | What's Missing                                                                                                                                                                                                |
|-------------------------------------------------|------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| messageId max 36 characters                     | Table 3    | `validateUniqueId` only checks empty, no length limit                                                                                                                                                         |
| 10 valid ErrorCode values                       | Table 7    | No constants: NotImplemented, NotSupported, InternalError, ProtocolError, SecurityError, FormationViolation, PropertyConstraintViolation, OccurenceConstraintViolation, TypeConstraintViolation, GenericError |
| ErrorDescription allows empty `""`              | Table 6    | `ErrErrorDescriptionAbsent` treats empty as error — spec says "otherwise a clear empty string" is valid                                                                                                       |
| Payload allows `null` or `{}`                   | Tables 4-5 | Not yet validated (Parse not implemented)                                                                                                                                                                     |
| Direction (Client-to-Server / Server-to-Client) | Table 2    | Informational only — not enforced at this layer                                                                                                                                                               |

### Notes on ErrorDescription

Spec Table 6 says: "Should be filled in if possible, otherwise a clear
empty string `""`." This means empty string IS valid on the wire. The
sentinel `ErrErrorDescriptionAbsent` should only apply when constructing
a CallError programmatically (encouraging callers to provide one), NOT
when parsing an incoming message.

### Notes on OccurenceConstraintViolation

The spec spells it "Occurence" (missing the second 'r'). This is a typo
in the spec itself but we MUST match the spec spelling exactly since
this string appears on the wire.

## Design Principle: First-Class Domain Types

Every spec concept with constraints or semantics MUST be its own Go
type — never a bare primitive. This is the most important design rule
in the project:

- `MessageType uint8` — already exists, enforces valid values 2/3/4
- `UniqueId string` — new type, constructor validates non-empty +
  max 36 chars (Table 3). Replaces bare `string` in all structs.
- `ErrorCode string` — new type, constructor validates against the
  10 values in Table 7. Replaces bare `string` in RawCallError.

Each type gets a validating constructor (`NewUniqueId`, `NewErrorCode`)
that is the only way to produce valid instances. The zero value is
invalid by design. Structs use these types in their fields so the
compiler prevents accidentally passing a raw string where a spec
concept is expected.

This rule applies to all future spec concepts as well.

## Steps (one at a time, in order)

### Step 0: First-class UniqueId and ErrorCode types

Create `type UniqueId string` with `NewUniqueId(string) (UniqueId, error)`
enforcing non-empty + max 36 chars. Create `type ErrorCode string` with
`NewErrorCode(string) (ErrorCode, error)` enforcing Table 7 values.
Add constants for all 10 ErrorCode values. Update all structs and
constructors to use these types instead of bare strings.

- **New files:** `unique_id.go`, `error_code.go`
- **Modified:** `raw_message.go`, `decoded_message.go`, `message.go`
- **Tests:** `unique_id_test.go`, `tests/unique_id_test.go`,
  `tests/error_code_test.go`, update `tests/message_test.go`

### Step 1: MarshalJSON on Raw message types

Add `MarshalJSON()` to RawCall, RawCallResult, RawCallError so messages
can be serialized back to canonical OCPP-J arrays. Uses existing
`marshalJSONArray` helper.

- **Files:** `raw_message.go`, `tests/message_test.go`

### Step 2: Parse function

Implement `Parse([]byte) (Message, error)` — the main entry point.
Unmarshals a JSON array, reads MessageTypeId, dispatches to type-specific
parsing (Call/CallResult/CallError), validates envelope fields.
Payload accepts both `null` and `{}`. ErrorDescription accepts `""`.

- **Files:** `parse.go`, `parse_test.go`, `tests/parse_test.go`

### Step 3: Registry + JSONDecoder

Thread-safe registry mapping Action names to decoder functions.
`JSONDecoder[Input, Output any]` adapter bridges `json.RawMessage` to
`ocpp16messages` constructors.

- **Files:** `registry.go`, `decoder.go`, `tests/registry_test.go`
- **Dependency:** add `ocpp16messages v1.0.3` to go.mod

### Step 4: Example tests

Executable documentation for Parse, MarshalJSON, Registry usage.

- **Files:** `tests/example_parse_test.go`, `tests/example_marshal_test.go`

### Step 5: Fuzz tests

Fuzz `Parse` with random bytes to catch panics/crashes.

- **Files:** `tests_fuzz/parse_fuzz_test.go` (build tag `fuzz`)

### Step 6: Race tests

Concurrent Registry read/write to verify thread safety.

- **Files:** `tests_race/registry_race_test.go` (build tag `race`)
