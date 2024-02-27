package broker

import (
	"voting/internal/events"

	"github.com/centrifugal/centrifuge"
	"github.com/sirupsen/logrus"
)

func NewInternal() Broker {

	broker := &ChannelBroker{
		NotifierAll:        make(chan events.Event, 1),
		NotifierUser:       make(chan events.UserBoundEvent, 1),
		NotifierAllButUser: make(chan events.UserBoundEvent, 1),
		NewClients:         make(chan UserBoundChannel),
		ClosingClients:     make(chan UserBoundChannel),
		Clients:            make(map[UserBoundChannel]bool),
	}

	go broker.Listen()
	go broker.SendHeartBeat()

	return broker
}

func NewCentrifuge() *centrifuge.Node {
	node, _ = centrifuge.New(centrifuge.Config{
		LogLevel:   centrifuge.LogLevelDebug,
		LogHandler: handleLog,
	})

	initHandlers()

	// We also start a separate goroutine for centrifuge itself, since we
	// still need to run gin web server.
	go func() {
		if err := node.Run(); err != nil {
			logrus.Fatal(err)
		}
		// ticker := time.NewTicker(5 * time.Second)
		// defer ticker.Stop()
		// // for {
		// 	<-ticker.C
		// 	logrus.Info("Logging to stdout")

		// 	q := models.Question{
		// 		Id:          "123",
		// 		Text:        "Centrifugo!",
		// 		Votes:       23,
		// 		Answered:    false,
		// 		Voted:       false,
		// 		Anonymous:   true,
		// 		CreatorName: "John Doe",
		// 		CreatorHash: "ysd1123a",
		// 	}

		// 	data, err := json.Marshal(q)
		// 	if err != nil {
		// 		logrus.Error("Error marshaling:", err)
		// 		return
		// 	}

		// 	node.Publish("voting", data)
		// }
	}()

	return node
}
