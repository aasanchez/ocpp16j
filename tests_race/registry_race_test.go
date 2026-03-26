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
	heavyWorkers    = 50
	heavyIterations = 500
	actionCount     = 20
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

// TestRegistry_ConcurrentDecode hammers Decode from 50
// goroutines, each doing 500 iterations, all hitting the
// same registered action.
func TestRegistry_ConcurrentDecode(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)
	_ = registry.Register("TestAction", decoder)

	payload := json.RawMessage(`{"name": "test"}`)

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for range heavyWorkers {
		go func() {
			defer waitGroup.Done()

			for range heavyIterations {
				_, decodeErr := registry.Decode(
					"TestAction", payload,
				)
				if decodeErr != nil {
					t.Errorf("decode failed: %v", decodeErr)
				}
			}
		}()
	}

	waitGroup.Wait()
}

// TestRegistry_ConcurrentRegisterAndDecode runs writers
// and readers simultaneously. 50 goroutines register
// different actions while another 50 goroutines decode
// from already-registered actions — all at the same time.
func TestRegistry_ConcurrentRegisterAndDecode(
	t *testing.T,
) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)
	payload := json.RawMessage(`{"name": "test"}`)

	// Pre-register a base action.
	_ = registry.Register("Base", decoder)

	var waitGroup sync.WaitGroup

	// Writers: register many actions concurrently.
	// Some will collide — that's expected.
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

	// Readers: decode the base action concurrently.
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

// TestRegistry_ConcurrentDecodeMultipleActions decodes
// from multiple different registered actions concurrently.
func TestRegistry_ConcurrentDecodeMultipleActions(
	t *testing.T,
) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)

	// Register multiple actions.
	for actionIndex := range actionCount {
		actionName := fmt.Sprintf(
			"Action%d", actionIndex,
		)
		_ = registry.Register(actionName, decoder)
	}

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

// TestRegistry_ConcurrentDecodeWithErrors mixes
// successful decodes and unknown-action errors from
// multiple goroutines.
func TestRegistry_ConcurrentDecodeWithErrors(
	t *testing.T,
) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)
	_ = registry.Register("Known", decoder)

	validPayload := json.RawMessage(`{"name": "ok"}`)
	invalidPayload := json.RawMessage(`{"name": ""}`)

	var waitGroup sync.WaitGroup

	waitGroup.Add(heavyWorkers)

	for workerIndex := range heavyWorkers {
		go func(index int) {
			defer waitGroup.Done()

			for range heavyIterations {
				if index%2 == 0 {
					// Success path.
					_, decodeErr := registry.Decode(
						"Known", validPayload,
					)
					if decodeErr != nil {
						t.Errorf(
							"decode failed: %v",
							decodeErr,
						)
					}
				} else {
					// Error paths: unknown action
					// and validation failure.
					_, _ = registry.Decode(
						"Unknown", validPayload,
					)
					_, _ = registry.Decode(
						"Known", invalidPayload,
					)
				}
			}
		}(workerIndex)
	}

	waitGroup.Wait()
}
