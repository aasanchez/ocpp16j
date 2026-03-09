# AGENTS.md

## Project Overview

This repository is `github.com/aasanchez/ocpp16j`.

It is a transport-focused Go library for **OCPP-J 1.6**. The scope is limited
to JSON frame correctness:

- parse and validate OCPP-J frame envelopes
- marshal frames back to canonical JSON arrays
- decode payloads through `github.com/aasanchez/ocpp16messages`
- expose stable envelope-level sentinel errors

This repository does **not** implement:

- business logic or charger/backend behavior
- WebSocket session management
- correlation/state machines beyond what the wire format requires

## Current Layout

The codebase is currently a **single root package** plus dedicated black-box,
fuzz, and race test packages:

```text
.
├── decoder.go
├── errors.go
├── frame.go
├── registry.go
├── doc.go
├── decoder_test.go
├── frame_test.go
├── registry_test.go
├── tests/
│   └── public_api_test.go
├── tests_fuzz/
│   ├── doc.go
│   ├── frame_fuzz_test.go
│   └── registry_fuzz_test.go
├── tests_race/
│   ├── doc.go
│   ├── frame_race_test.go
│   └── registry_race_test.go
├── Makefile
├── golangci.yml
└── README.md
```

Key files:

- `frame.go`: raw OCPP-J frame model and parsing/marshaling
- `registry.go`: action registry for typed payload decoding
- `decoder.go`: adapter from `ocpp16messages` constructors to registry decoders
- `errors.go`: exported sentinel errors

## Go Version and Dependencies

- Go version: `1.24.6` (see `go.mod`)
- Main dependency: `github.com/aasanchez/ocpp16messages v1.0.3`

Do not introduce new dependencies unless they provide clear value to the
transport layer.

## Common Commands

### Dependencies

```sh
go mod tidy
```

### Tests

```sh
go test ./...
make test
make test-fuzz
make test-race
```

Notes:

- `make test` writes reports into `reports/`
- `tests_fuzz/` contains dedicated fuzz targets
- `tests_race/` contains concurrency tests enabled with the `race` build tag

### Lint and Formatting

```sh
make lint
make format
```

`make lint` is expected to pass with no ignored failures. It runs:

- `golangci-lint`
- `go vet`
- `staticcheck`

## API Design Rules

### Keep the public API transport-oriented

Public types and functions should remain focused on wire-format concerns.

Good fits:

- frame parsing
- frame serialization
- payload decoding adapters
- registry helpers for action-to-decoder mapping

Bad fits:

- reconnect policies
- message dispatch orchestration
- charge point / CSMS behavior
- protocol profile workflows

### `CALLRESULT` decoding requires context

OCPP-J `CALLRESULT` does not carry the action name on the wire. That is an
actual protocol constraint and should stay explicit in the API.

Keep this design principle:

- `Parse` returns raw frames
- `Registry.DecodeCall` uses the action embedded in `CALL`
- `Registry.DecodeCallResult` requires the caller to provide the related action

Do not hide that constraint behind guessed state.

### Validation layering

Validation is split across two levels:

- envelope validation belongs in this repo
- payload/domain validation belongs in `ocpp16messages`

When decoding payloads, prefer routing through upstream constructors using
`JSONDecoder(...)` rather than unmarshalling directly into validated message
types.

## Test Organization

Keep tests separated by visibility and responsibility.

### Root package tests

Files in the repo root use `package ocpp16json` and may test internal helpers.

Current pattern:

- `decoder_test.go`: tests `decoder.go`
- `frame_test.go`: tests `frame.go`
- `registry_test.go`: tests `registry.go`

Rules:

- keep tests atomic
- prefer one behavior per test
- keep tests close to the source file they exercise
- test private helpers from the same package only when needed

### Public API tests

Public black-box tests live in `tests/` and use `package ocpp16json_test`.

Current file:

- `tests/public_api_test.go`

Use that package to verify behavior as a consumer would see it:

- imports through `github.com/aasanchez/ocpp16j`
- interaction with `ocpp16messages`
- wrapped validation behavior

### Fuzz and race tests

Specialized transport stress tests live outside the root package:

- `tests_fuzz/`: fuzz targets for raw frame parsing and registry decoding
- `tests_race/`: race-detector-focused black-box tests

Keep those tests:

- small and targeted
- free of business-logic assumptions
- aligned with the existing transport-only API boundaries

## Coverage Expectations

This repository currently targets very high coverage for the root package.

When adding logic:

- add tests in the root package for internal branches
- add or update black-box tests in `tests/` for exported behavior
- add `tests_fuzz/` or `tests_race/` coverage when concurrency or parser
  robustness changes
- keep `go test ./...` green

## Lint Expectations

This repo uses a strict `golangci-lint` profile.

Practical consequences:

- add doc comments for exported API
- avoid swallowing lint failures in scripts
- keep line lengths and whitespace consistent
- prefer explicit stable errors over ad-hoc dynamic errors

Two lint rules are intentionally narrowed in `golangci.yml` because they
conflict with the intended package design:

- `ireturn`
- `revive:max-public-structs`

Do not broaden the lint disables casually.

## Editing Guidance

- preserve ASCII unless a file already requires otherwise
- prefer stable sentinel errors and wrapped errors
- keep exported API names aligned with OCPP-J terminology
- update `README.md` when API behavior or repo usage changes materially
- do not add generated files or editor artifacts

Ignore and remove workspace junk such as `.DS_Store` if it appears.
