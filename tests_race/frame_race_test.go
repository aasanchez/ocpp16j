//go:build race

package ocpp16json_test

import (
	"encoding/json"
	"sync"
	"testing"

	ocpp16json "github.com/aasanchez/ocpp16j"
)

func TestParseAndMarshalConcurrently(t *testing.T) {
	t.Parallel()

	frames := []ocpp16json.Frame{
		ocpp16json.RawCall{
			UniqueID: "call-1",
			Action:   "Authorize",
			Payload:  json.RawMessage(`{"idTag":"RFID-123"}`),
		},
		ocpp16json.RawCallResult{
			UniqueID: "result-1",
			Payload:  json.RawMessage(`{"currentTime":"2025-01-02T15:04:05Z"}`),
		},
		ocpp16json.CallError{
			UniqueID:         "error-1",
			ErrorCode:        "ProtocolError",
			ErrorDescription: "bad payload",
			ErrorDetails: map[string]any{
				"field": "idTag",
			},
		},
	}

	var waitGroup sync.WaitGroup

	for worker := 0; worker < 8; worker++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for iteration := 0; iteration < 200; iteration++ {
				for _, frame := range frames {
					data, err := json.Marshal(frame)
					if err != nil {
						t.Errorf("Marshal(%T): %v", frame, err)
						return
					}

					parsed, err := ocpp16json.Parse(data)
					if err != nil {
						t.Errorf("Parse(%T): %v", frame, err)
						return
					}

					if parsed.MessageType() != frame.MessageType() {
						t.Errorf(
							"message type mismatch: got %v want %v",
							parsed.MessageType(),
							frame.MessageType(),
						)
						return
					}
				}
			}
		}()
	}

	waitGroup.Wait()
}

func TestParseConcurrentlyAcrossFrameShapes(t *testing.T) {
	t.Parallel()

	rawFrames := [][]byte{
		[]byte(`[2,"call-1","Authorize",{"idTag":"RFID-123"}]`),
		[]byte(`[3,"result-1",{"currentTime":"2025-01-02T15:04:05Z"}]`),
		[]byte(`[4,"error-1","ProtocolError","bad payload",{"field":"idTag"}]`),
	}

	var waitGroup sync.WaitGroup

	for worker := 0; worker < 8; worker++ {
		waitGroup.Add(1)

		go func(workerIndex int) {
			defer waitGroup.Done()

			for iteration := 0; iteration < 200; iteration++ {
				for frameIndex, data := range rawFrames {
					frame, err := ocpp16json.Parse(data)
					if err != nil {
						t.Errorf(
							"Parse(worker=%d, frame=%d): %v",
							workerIndex,
							frameIndex,
							err,
						)
						return
					}

					if frame.MessageID() == "" {
						t.Errorf(
							"empty message id for worker=%d frame=%d",
							workerIndex,
							frameIndex,
						)
						return
					}
				}
			}
		}(worker)
	}

	waitGroup.Wait()
}
