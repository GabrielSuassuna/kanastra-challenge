package handler

import (
	"encoding/json"
	"kanastra-api/pkg/events"
	"sync"
)

type FileEventHandler struct {
	SseChannel *chan string
}

func NewFileEventHandler(sseChannel *chan string) *FileEventHandler {
	return &FileEventHandler{
		SseChannel: sseChannel,
	}
}

func (h *FileEventHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	var payload string

	// check if channel is not initialized
	if *h.SseChannel == nil {
		return
	}

	// check if payload is not a string
	if _, ok := event.GetPayload().(string); !ok {
		// transform payload to a json string
		data, err := json.Marshal(event.GetPayload())
		if err != nil {
			return
		}
		payload = string(data)
	} else {
		payload = event.GetPayload().(string)
	}

	*h.SseChannel <- payload
}
