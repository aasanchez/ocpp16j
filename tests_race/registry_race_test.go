//go:build race

package ocpp16json_race

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

const (
	concurrentWorkers = 10
	decodeIterations  = 100
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

// TestRegistry_ConcurrentDecode verifies that multiple
// goroutines can call Decode simultaneously on the same
// Registry without data races.
func TestRegistry_ConcurrentDecode(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)
	_ = registry.Register("TestAction", decoder)

	payload := json.RawMessage(`{"name": "test"}`)

	var waitGroup sync.WaitGroup

	waitGroup.Add(concurrentWorkers)

	for range concurrentWorkers {
		go func() {
			defer waitGroup.Done()

			for range decodeIterations {
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

// TestRegistry_ConcurrentRegisterAndDecode verifies that
// Register and Decode can run concurrently without data
// races. Some registers will fail with
// ErrActionAlreadyRegistered — that is expected.
func TestRegistry_ConcurrentRegisterAndDecode(
	t *testing.T,
) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()

	decoder := ocpp16json.JSONDecoder(raceConstructor)
	payload := json.RawMessage(`{"name": "test"}`)

	// Pre-register one action so Decode has something
	// to work with from the start.
	_ = registry.Register("Action0", decoder)

	var waitGroup sync.WaitGroup

	// Writers: register new actions concurrently.
	waitGroup.Add(concurrentWorkers)

	for workerIndex := range concurrentWorkers {
		go func(index int) {
			defer waitGroup.Done()

			actionName := "Action" + string(
				rune('A'+index),
			)
			_ = registry.Register(actionName, decoder)
		}(workerIndex)
	}

	// Readers: decode concurrently.
	waitGroup.Add(concurrentWorkers)

	for range concurrentWorkers {
		go func() {
			defer waitGroup.Done()

			for range decodeIterations {
				// Action0 is always registered.
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
