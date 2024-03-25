package voting_ws

import (
	shared_infra_broker "voting/shared/infra/broker"

	"github.com/centrifugal/centrifuge"
	"github.com/sirupsen/logrus"
)

func NewCentrifuge(broker shared_infra_broker.Broker) *centrifuge.Node {
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
