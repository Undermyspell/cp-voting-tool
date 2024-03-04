package notification

import (
	"voting/internal/broker"

	"github.com/centrifugal/centrifuge"
	"github.com/sirupsen/logrus"
)

func NewCentrifuge(broker broker.Broker) *centrifuge.Node {
	node, _ = centrifuge.New(centrifuge.Config{
		LogLevel:   centrifuge.LogLevelDebug,
		LogHandler: handleLog,
	})

	initHandlers(broker)

	// We also start a separate goroutine for centrifuge itself, since we
	// still need to run gin web server.
	go func() {
		if err := node.Run(); err != nil {
			logrus.Fatal(err)
		}
	}()

	return node
}
