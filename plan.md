# OCPP-J 1.6 JSON Wrapper — Implementation Plan

## Completed

- **Step 0**: First-class `UniqueId` and `ErrorCode` domain types
- **Step 1**: `MarshalJSON()` on `Call`, `CallResult`, `CallError`
- **Step 2**: `Parse([]byte) (Message, error)` function
- **Step 3**: `Registry` + `JSONDecoder` for payload decoding
- **Step 3b**: `NewCall`, `NewCallResult`, `NewCallError` constructors
- **Step 4**: Example tests for pkgsite (19 examples)
- **Refactor**: Split into `call.go`, `call_result.go`, `call_error.go`
- **Refactor**: Rename `RawCall` → `Call`, constants → `MessageTypeCall`

## Remaining

### Step 5: Fuzz tests

Directory: `tests_fuzz/`, build tag `//go:build fuzz`

- `FuzzParse` — random bytes to `Parse()`, never panics
- `FuzzNewUniqueId` — random strings, no panics
- `FuzzNewErrorCode` — random strings, no panics

### Step 6: Race tests

Directory: `tests_race/`, build tag `//go:build race`

- `TestRegistry_ConcurrentRegisterAndDecode`

## Terminology Mapping (spec → code)

| Spec Term              | Go Identifier            |
|------------------------|--------------------------|
| CALL (MessageType 2)   | `MessageTypeCall` const  |
| CALLRESULT (Type 3)    | `MessageTypeCallResult`  |
| CALLERROR (Type 4)     | `MessageTypeCallError`   |
| Call message struct     | `Call` struct             |
| CallResult struct       | `CallResult` struct      |
| CallError struct        | `CallError` struct       |
| UniqueId (Table 3)     | `UniqueId` type          |
| Action (Table 4)       | `Action string` field    |
| Payload (Tables 4-5)   | `Payload json.RawMessage`|
| ErrorCode (Table 7)    | `ErrorCode` type         |
| ErrorDescription        | `ErrorDescription string`|
| ErrorDetails            | `ErrorDetails map`       |

## Design Principle: First-Class Domain Types

Every spec concept with constraints MUST be its own Go type.
Structs use these types in fields so the compiler prevents
accidentally passing a raw string where a spec concept is
expected. Each type has a validating constructor; the zero
value is invalid by design.
