package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sse/internal/events"

	"github.com/gin-gonic/gin"
)

type Broker struct {
	Notifier       chan string
	NewClients     chan chan string
	ClosingClients chan chan string
	Clients        map[chan string]bool
}

type Message struct {
	Message *string `json:"message"`
}

func New() *Broker {
	return &Broker{
		Notifier:       make(chan string, 1),
		NewClients:     make(chan chan string),
		ClosingClients: make(chan chan string),
		Clients:        make(map[chan string]bool),
	}
}

func (broker *Broker) StreamNew(event events.Events) func(c *gin.Context) {
	return func(c *gin.Context) {
		messageChan := make(chan string)

		defer func() {
			broker.ClosingClients <- messageChan
			close(messageChan)
		}()

		broker.NewClients <- messageChan

		c.Stream(func(w io.Writer) bool {
			select {
			case msg := <-messageChan:
				c.SSEvent(string(event), msg)
			case <-c.Request.Context().Done():
				return false
			}
			return true
		})
	}
}

func (broker *Broker) BroadcastMessage(c *gin.Context) {
	var message Message
	err := c.BindJSON(&message)
	res, _ := json.Marshal(message)
	fmt.Println(string(res))

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	broker.Notifier <- string(res)
}

func (broker *Broker) Listen() {
	for {
		select {
		case s := <-broker.NewClients:
			// A new client has joined
			broker.Clients[s] = true
			log.Printf("🟢 Client added. %d registered clients", len(broker.Clients))
		case s := <-broker.ClosingClients:
			// A client has detached
			// remove them from our clients map
			delete(broker.Clients, s)
			log.Printf("🔴 Removed client. %d registered clients", len(broker.Clients))
		case event := <-broker.Notifier:
			// case for getting a new msg
			// Thus send it to all clients
			for clientMessageChan := range broker.Clients {
				fmt.Println("send to client")
				fmt.Println(event)
				clientMessageChan <- event
			}
		}
	}
}
