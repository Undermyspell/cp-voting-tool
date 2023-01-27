package broker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Broker struct {
	Notifier       chan []byte
	NewClients     chan chan []byte
	ClosingClients chan chan []byte
	Clients        map[chan []byte]bool
}

type Message struct {
	Message *string `json:"message"`
}

func New() *Broker {
	return &Broker{
		Notifier:       make(chan []byte, 1),
		NewClients:     make(chan chan []byte),
		ClosingClients: make(chan chan []byte),
		Clients:        make(map[chan []byte]bool),
	}
}

func (broker *Broker) Stream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.NewClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.ClosingClients <- messageChan
	}()

	go func() {
		// Listen to connection close and un-register messageChan
		<-r.Context().Done()
		broker.ClosingClients <- messageChan
	}()

	for {
		// Write to the ResponseWriter
		// Server Sent Events compatible
		fmt.Fprintf(w, "data: %s\n\n", <-messageChan)

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}
}

func (broker *Broker) BroadcastMessage(w http.ResponseWriter, r *http.Request) {
	var message Message

	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cant parse request"))
		return
	}

	byteData, err := json.Marshal(message)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("paging struct to byte"))
		return
	}

	broker.Notifier <- byteData
}
