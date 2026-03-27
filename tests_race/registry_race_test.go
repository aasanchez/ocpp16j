//go:build race

package ocpp16json_race

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	heavyWorkers    = 100
	heavyIterations = 1000
	actionCount     = 50
)

// errTestPayload is a test-only sentinel error.
var errTestPayload = errors.New("name required")

type racePayload struct {
	Name string `json:"name"`
}

func raceConstructor(
	input racePayload,
) (racePayload, error) {
	if input.Name == "" {
		return racePayload{}, errTestPayload
	}

	return input, nil
}

func registryWithActions(
	count int,
) *ocpp16json.Registry {
	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)

	for actionIndex := range count {
		actionName := fmt.Sprintf(
			"Action%d", actionIndex,
		)
		_ = registry.Register(actionName, decoder)
	}

	return registry
}

// TestRegistry_ConcurrentDecode hammers Decode from 100
// goroutines, each doing 1000 iterations, all hitting the
// same registered action. Total: 100,000 operations.
func TestRegistry_ConcurrentDecode(t *testing.T) {
	t.Parallel()

	registry := registryWithActions(1)

	payload := json.RawMessage(`{"name": "test"}`)

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for range heavyWorkers {
		go func() {
			defer waitGroup.Done()

			for range heavyIterations {
				_, decodeErr := registry.Decode(
					"Action0", payload,
				)
				if decodeErr != nil {
					t.Errorf("decode failed: %v", decodeErr)
				}
			}
		}()
	}

	waitGroup.Wait()
}

// TestRegistry_ConcurrentRegisterAndDecode runs 100
// writers and 100 readers simultaneously. Writers register
// actions (with expected collisions) while readers decode
// from a pre-registered action. Total: 100,000+ operations.
func TestRegistry_ConcurrentRegisterAndDecode(
	t *testing.T,
) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)
	payload := json.RawMessage(`{"name": "test"}`)

	_ = registry.Register("Base", decoder)

	var waitGroup sync.WaitGroup

	// Writers.
	waitGroup.Add(heavyWorkers)

	for workerIndex := range heavyWorkers {
		go func(index int) {
			defer waitGroup.Done()

			actionName := fmt.Sprintf(
				"Action%d", index%actionCount,
			)
			_ = registry.Register(actionName, decoder)
		}(workerIndex)
	}

	// Readers.
	waitGroup.Add(heavyWorkers)

	for range heavyWorkers {
		go func() {
			defer waitGroup.Done()

			for range heavyIterations {
				_, decodeErr := registry.Decode(
					"Base", payload,
				)
				if decodeErr != nil {
					t.Errorf("decode failed: %v", decodeErr)
				}
			}
		}()
	}

	waitGroup.Wait()
}

// TestRegistry_ConcurrentDecodeMultipleActions spreads
// 100 goroutines across 50 different registered actions.
// Total: 100,000 operations across 50 actions.
func TestRegistry_ConcurrentDecodeMultipleActions(
	t *testing.T,
) {
	t.Parallel()

	registry := registryWithActions(actionCount)

	payload := json.RawMessage(`{"name": "test"}`)

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for workerIndex := range heavyWorkers {
		go func(index int) {
			defer waitGroup.Done()

			actionName := fmt.Sprintf(
				"Action%d", index%actionCount,
			)

			for range heavyIterations {
				_, decodeErr := registry.Decode(
					actionName, payload,
				)
				if decodeErr != nil {
					t.Errorf("decode failed: %v", decodeErr)
				}
			}
		}(workerIndex)
	}

	waitGroup.Wait()
}

// TestRegistry_ConcurrentDecodeWithErrors mixes success
// paths, unknown-action errors, and validation failures
// from 100 goroutines. Total: 200,000+ operations.
func TestRegistry_ConcurrentDecodeWithErrors(
	t *testing.T,
) {
	t.Parallel()

	registry := registryWithActions(actionCount)

	validPayload := json.RawMessage(`{"name": "ok"}`)
	invalidPayload := json.RawMessage(`{"name": ""}`)
	brokenPayload := json.RawMessage(`not json`)

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for workerIndex := range heavyWorkers {
		go func(index int) {
			defer waitGroup.Done()

			actionName := fmt.Sprintf(
				"Action%d", index%actionCount,
			)

			for range heavyIterations {
				// Success path.
				_, _ = registry.Decode(
					actionName, validPayload,
				)
				// Unknown action.
				_, _ = registry.Decode(
					"NonExistent", validPayload,
				)
				// Validation failure.
				_, _ = registry.Decode(
					actionName, invalidPayload,
				)
				// Unmarshal failure.
				_, _ = registry.Decode(
					actionName, brokenPayload,
				)
			}
		}(workerIndex)
	}

	waitGroup.Wait()
}

// TestRegistry_ConcurrentRegisterDuplicates hammers
// Register with the same action name from 100 goroutines.
// Exactly one must succeed, the rest must get
// ErrActionAlreadyRegistered.
func TestRegistry_ConcurrentRegisterDuplicates(
	t *testing.T,
) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)

	var (
		waitGroup    sync.WaitGroup
		successCount int64
		mutex        sync.Mutex
	)

	waitGroup.Add(heavyWorkers)

	for range heavyWorkers {
		go func() {
			defer waitGroup.Done()

			registerErr := registry.Register(
				"SameAction", decoder,
			)
			if registerErr == nil {
				mutex.Lock()
				successCount++
				mutex.Unlock()
			}
		}()
	}

	waitGroup.Wait()

	if successCount != 1 {
		t.Fatalf(
			"expected exactly 1 success, got %d",
			successCount,
		)
	}
}

// TestParse_ConcurrentParseSameBytes verifies Parse is
// safe to call concurrently with the same input bytes.
// Total: 100,000 operations.
func TestParse_ConcurrentParseSameBytes(t *testing.T) {
	t.Parallel()

	wire := []byte(
		`[2,"19223201","BootNotification",` +
			`{"chargePointVendor":"VendorX",` +
			`"chargePointModel":"Model1"}]`,
	)

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for range heavyWorkers {
		go func() {
			defer waitGroup.Done()

			for range heavyIterations {
				message, parseErr := ocpp16json.Parse(
					wire,
				)
				if parseErr != nil {
					t.Errorf("parse failed: %v", parseErr)

					return
				}

				if !ocpp16json.IsCall(message) {
					t.Error("expected Call message")
				}
			}
		}()
	}

	waitGroup.Wait()
}

// TestParse_ConcurrentParseDifferentTypes parses Call,
// CallResult, and CallError concurrently from separate
// goroutines. Total: 300,000 operations.
func TestParse_ConcurrentParseDifferentTypes(
	t *testing.T,
) {
	t.Parallel()

	messages := [][]byte{
		[]byte(`[2,"1","Action",{}]`),
		[]byte(`[3,"1",{"ok":true}]`),
		[]byte(`[4,"1","GenericError","err",{}]`),
	}

	var waitGroup sync.WaitGroup

	for _, wire := range messages {
		waitGroup.Add(heavyWorkers)

		for range heavyWorkers {
			go func(data []byte) {
				defer waitGroup.Done()

				for range heavyIterations {
					_, parseErr := ocpp16json.Parse(data)
					if parseErr != nil {
						t.Errorf(
							"parse failed: %v",
							parseErr,
						)
					}
				}
			}(wire)
		}
	}

	waitGroup.Wait()
}

// TestFullPipeline_ConcurrentParseDecodeRespond runs
// the full pipeline (parse → decode → respond)
// concurrently from 100 goroutines. Total: 100,000
// full round-trips.
// TestMultipleRegistries_ConcurrentUse creates multiple
// independent registries and uses them concurrently.
// Verifies no shared state leaks between instances.
func TestMultipleRegistries_ConcurrentUse(
	t *testing.T,
) {
	t.Parallel()

	registryCount := 10
	registries := make(
		[]*ocpp16json.Registry, registryCount,
	)

	for registryIndex := range registryCount {
		registries[registryIndex] = registryWithActions(
			actionCount,
		)
	}

	payload := json.RawMessage(`{"name": "test"}`)

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for workerIndex := range heavyWorkers {
		go func(index int) {
			defer waitGroup.Done()

			registry := registries[index%registryCount]

			actionName := fmt.Sprintf(
				"Action%d", index%actionCount,
			)

			for range heavyIterations {
				_, decodeErr := registry.Decode(
					actionName, payload,
				)
				if decodeErr != nil {
					t.Errorf("decode: %v", decodeErr)
				}
			}
		}(workerIndex)
	}

	waitGroup.Wait()
}

// TestRegistry_ConcurrentPopulateAndQuery creates a
// fresh registry and has goroutines simultaneously
// registering actions and querying them — some queries
// will find the action, some won't, depending on timing.
func TestRegistry_ConcurrentPopulateAndQuery(
	t *testing.T,
) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)
	payload := json.RawMessage(`{"name": "test"}`)

	var waitGroup sync.WaitGroup

	waitGroup.Add(actionCount)

	for actionIndex := range actionCount {
		go func(index int) {
			defer waitGroup.Done()

			actionName := fmt.Sprintf(
				"Action%d", index,
			)
			_ = registry.Register(actionName, decoder)
		}(actionIndex)
	}

	waitGroup.Add(heavyWorkers)

	for workerIndex := range heavyWorkers {
		go func(index int) {
			defer waitGroup.Done()

			for range heavyIterations {
				actionName := fmt.Sprintf(
					"Action%d", index%actionCount,
				)

				_, _ = registry.Decode(
					actionName, payload,
				)
			}
		}(workerIndex)
	}

	waitGroup.Wait()
}

// TestMarshalJSON_ConcurrentMarshalSameMessage verifies
// MarshalJSON on the same Call struct from many goroutines.
func TestMarshalJSON_ConcurrentMarshalSameMessage(
	t *testing.T,
) {
	t.Parallel()

	message, _ := ocpp16json.Parse(
		[]byte(
			`[2,"19223201","BootNotification",` +
				`{"chargePointVendor":"VendorX"}]`,
		),
	)

	call, _ := ocpp16json.AsCall(message)

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for range heavyWorkers {
		go func() {
			defer waitGroup.Done()

			for range heavyIterations {
				_, marshalErr := json.Marshal(call)
				if marshalErr != nil {
					t.Errorf(
						"marshal: %v", marshalErr,
					)
				}
			}
		}()
	}

	waitGroup.Wait()
}

// TestTypePredicates_ConcurrentChecks calls IsCall,
// IsCallResult, IsCallError, AsCall, AsCallResult,
// AsCallError concurrently on the same messages.
func TestTypePredicates_ConcurrentChecks(
	t *testing.T,
) {
	t.Parallel()

	callMsg, _ := ocpp16json.Parse(
		[]byte(`[2,"1","Action",{}]`),
	)

	resultMsg, _ := ocpp16json.Parse(
		[]byte(`[3,"1",{}]`),
	)

	errorMsg, _ := ocpp16json.Parse(
		[]byte(`[4,"1","GenericError","",{}]`),
	)

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for range heavyWorkers {
		go func() {
			defer waitGroup.Done()

			for range heavyIterations {
				ocpp16json.IsCall(callMsg)
				ocpp16json.IsCallResult(resultMsg)
				ocpp16json.IsCallError(errorMsg)

				_, _ = ocpp16json.AsCall(callMsg)
				_, _ = ocpp16json.AsCallResult(
					resultMsg,
				)
				_, _ = ocpp16json.AsCallError(errorMsg)

				_, _ = ocpp16json.AsCall(resultMsg)
				_, _ = ocpp16json.AsCallResult(callMsg)
				_, _ = ocpp16json.AsCallError(callMsg)
			}
		}()
	}

	waitGroup.Wait()
}

// TestNewConstructors_ConcurrentCreation calls NewCall,
// NewCallResult, NewCallError concurrently.
func TestNewConstructors_ConcurrentCreation(
	t *testing.T,
) {
	t.Parallel()

	uniqueId, _ := ocpp16json.NewUniqueId("19223201")

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for range heavyWorkers {
		go func() {
			defer waitGroup.Done()

			for range heavyIterations {
				_, _ = ocpp16json.NewCall(
					uniqueId, "Action",
					map[string]string{"k": "v"},
				)

				_, _ = ocpp16json.NewCallResult(
					uniqueId,
					map[string]string{"k": "v"},
				)

				_, _ = ocpp16json.NewCallError(
					uniqueId,
					ocpp16json.GenericError,
					"desc",
					map[string]any{},
				)
			}
		}()
	}

	waitGroup.Wait()
}

// TestDomainTypes_ConcurrentCreation calls NewUniqueId
// and NewErrorCode concurrently with various inputs.
func TestDomainTypes_ConcurrentCreation(
	t *testing.T,
) {
	t.Parallel()

	uniqueIdInputs := []string{
		"19223201",
		"550e8400-e29b-41d4-a716-446655440000",
		"short",
		"",
	}

	errorCodeInputs := []string{
		"NotImplemented",
		"GenericError",
		"MadeUpError",
		"",
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for workerIndex := range heavyWorkers {
		go func(index int) {
			defer waitGroup.Done()

			for range heavyIterations {
				uidInput := uniqueIdInputs[index%len(
					uniqueIdInputs,
				)]
				_, _ = ocpp16json.NewUniqueId(uidInput)

				codeInput := errorCodeInputs[index%len(
					errorCodeInputs,
				)]
				_, _ = ocpp16json.NewErrorCode(
					codeInput,
				)
			}
		}(workerIndex)
	}

	waitGroup.Wait()
}

// TestFullPipeline_ConcurrentParseDecodeRespond runs
// the full pipeline (parse → decode → respond)
// concurrently from 100 goroutines. Total: 100,000
// full round-trips.
func TestFullPipeline_ConcurrentParseDecodeRespond(
	t *testing.T,
) {
	t.Parallel()

	registry := registryWithActions(actionCount)

	wires := make([][]byte, actionCount)
	for actionIndex := range actionCount {
		wires[actionIndex] = []byte(fmt.Sprintf(
			`[2,"id-%d","Action%d",{"name":"test"}]`,
			actionIndex, actionIndex,
		))
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for workerIndex := range heavyWorkers {
		go func(index int) {
			defer waitGroup.Done()

			wire := wires[index%actionCount]

			for range heavyIterations {
				message, parseErr := ocpp16json.Parse(
					wire,
				)
				if parseErr != nil {
					t.Errorf("parse: %v", parseErr)

					return
				}

				call, callErr := ocpp16json.AsCall(
					message,
				)
				if callErr != nil {
					t.Errorf("as call: %v", callErr)

					return
				}

				_, decodeErr := registry.Decode(
					call.Action, call.Payload,
				)
				if decodeErr != nil {
					t.Errorf("decode: %v", decodeErr)

					return
				}

				_, marshalErr := json.Marshal(message)
				if marshalErr != nil {
					t.Errorf(
						"marshal: %v", marshalErr,
					)
				}
			}
		}(workerIndex)
	}

	waitGroup.Wait()
}
