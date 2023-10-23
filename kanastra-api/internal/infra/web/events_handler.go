package web

import (
	"fmt"
	"net/http"
)

type EventHandler struct {
	SseChannel *chan string
}

func NewEventHandler(sseChanel *chan string) *EventHandler {
	return &EventHandler{
		SseChannel: sseChanel,
	}
}

func (h *EventHandler) SendEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New connection established")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	*h.SseChannel = make(chan string)

	defer func() {
		close(*h.SseChannel)
		fmt.Println("Closing connection")
	}()

	flusher, ok := w.(http.Flusher)

	if !ok {
		fmt.Println("Unable to initialize flusher")
	}

	for {
		select {
		case msg := <-*h.SseChannel:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()

		case <-r.Context().Done():
			fmt.Println("Connection closed")
			return
		}
	}
}
