//go:build race

package ocpp16json_test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

type raceAuthorizeInput struct {
	IDTag string `json:"idTag"`
}

type raceAuthorizePayload struct {
	IDTag string
}

type raceHeartbeatInput struct {
	CurrentTime string `json:"currentTime"`
}

type raceHeartbeatPayload struct {
	CurrentTime string
}

func TestRegistryConcurrentRequestRegistrationAndDecode(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	decodeCall := []byte(`[2,"uid-1","Authorize",{"idTag":"RFID-123"}]`)

	err := registry.RegisterRequest(
		"Authorize",
		ocpp16json.JSONDecoder(raceAuthorize),
	)
	if err != nil {
		t.Fatalf("RegisterRequest baseline: %v", err)
	}

	var waitGroup sync.WaitGroup

	for worker := 0; worker < 6; worker++ {
		waitGroup.Add(1)

		go func(workerIndex int) {
			defer waitGroup.Done()

			for iteration := 0; iteration < 100; iteration++ {
				action := fmt.Sprintf("Authorize-%d-%d", workerIndex, iteration)

				err := registry.RegisterRequest(
					action,
					ocpp16json.JSONDecoder(raceAuthorize),
				)
				if err != nil {
					t.Errorf("RegisterRequest(%s): %v", action, err)
					return
				}

				decoded, err := registry.DecodeCall(decodeCall)
				if err != nil {
					t.Errorf("DecodeCall: %v", err)
					return
				}

				payload, ok := decoded.Payload.(raceAuthorizePayload)
				if !ok {
					t.Errorf("unexpected payload type: %T", decoded.Payload)
					return
				}

				if payload.IDTag != "RFID-123" {
					t.Errorf("unexpected idTag: %q", payload.IDTag)
					return
				}
			}
		}(worker)
	}

	waitGroup.Wait()
}

func TestRegistryConcurrentConfirmationRegistrationAndDecode(t *testing.T) {
	t.Parallel()

	registry := ocpp16json.NewRegistry()
	decodeResult := []byte(
		`[3,"uid-1",{"currentTime":"2025-01-02T15:04:05Z"}]`,
	)

	err := registry.RegisterConfirmation(
		"Heartbeat",
		ocpp16json.JSONDecoder(raceHeartbeat),
	)
	if err != nil {
		t.Fatalf("RegisterConfirmation baseline: %v", err)
	}

	var waitGroup sync.WaitGroup

	for worker := 0; worker < 6; worker++ {
		waitGroup.Add(1)

		go func(workerIndex int) {
			defer waitGroup.Done()

			for iteration := 0; iteration < 100; iteration++ {
				action := fmt.Sprintf("Heartbeat-%d-%d", workerIndex, iteration)

				err := registry.RegisterConfirmation(
					action,
					ocpp16json.JSONDecoder(raceHeartbeat),
				)
				if err != nil {
					t.Errorf("RegisterConfirmation(%s): %v", action, err)
					return
				}

				decoded, err := registry.DecodeCallResult(
					"Heartbeat",
					decodeResult,
				)
				if err != nil {
					t.Errorf("DecodeCallResult: %v", err)
					return
				}

				payload, ok := decoded.Payload.(raceHeartbeatPayload)
				if !ok {
					t.Errorf("unexpected payload type: %T", decoded.Payload)
					return
				}

				if payload.CurrentTime != "2025-01-02T15:04:05Z" {
					t.Errorf(
						"unexpected currentTime: %q",
						payload.CurrentTime,
					)
					return
				}
			}
		}(worker)
	}

	waitGroup.Wait()
}

func raceAuthorize(input raceAuthorizeInput) (raceAuthorizePayload, error) {
	return raceAuthorizePayload{
		IDTag: input.IDTag,
	}, nil
}

func raceHeartbeat(input raceHeartbeatInput) (raceHeartbeatPayload, error) {
	return raceHeartbeatPayload{
		CurrentTime: input.CurrentTime,
	}, nil
}

func TestJSONDecoderConcurrentUse(t *testing.T) {
	t.Parallel()

	decoder := ocpp16json.JSONDecoder(raceAuthorize)
	raw := json.RawMessage(`{"idTag":"RFID-123"}`)

	var waitGroup sync.WaitGroup

	for worker := 0; worker < 8; worker++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for iteration := 0; iteration < 200; iteration++ {
				payload, err := decoder(raw)
				if err != nil {
					t.Errorf("decoder: %v", err)
					return
				}

				request, ok := payload.(raceAuthorizePayload)
				if !ok {
					t.Errorf("unexpected payload type: %T", payload)
					return
				}

				if request.IDTag != "RFID-123" {
					t.Errorf("unexpected idTag: %q", request.IDTag)
					return
				}
			}
		}()
	}

	waitGroup.Wait()
}
