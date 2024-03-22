package notification

import (
	"context"
	"encoding/json"
	"time"
	"voting/internal/events"
	shared_infra_broker "voting/shared/infra/broker"
	"voting/shared/shared_models"

	"github.com/centrifugal/centrifuge"
	"github.com/sirupsen/logrus"
)

func handleLog(e centrifuge.LogEntry) {
	logrus.Infof("%s: %v", e.Message, e.Fields)
}

var node *centrifuge.Node

func initHandlers(internalBroker shared_infra_broker.Broker) {
	node.OnConnecting(func(ctx context.Context, event centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		userContext, err := shared_models.GetUserContextFromToken(event.Token)

		if err != nil {
			return centrifuge.ConnectReply{}, centrifuge.ErrorUnauthorized
		}

		userContextByteArray, errMarshal := json.Marshal(userContext)

		if errMarshal != nil {
			logrus.Error(errMarshal)
			return centrifuge.ConnectReply{}, centrifuge.ErrorUnauthorized
		}

		_, ok := centrifuge.GetCredentials(ctx)
		if !ok {
			return centrifuge.ConnectReply{}, centrifuge.DisconnectServerError
		}
		return centrifuge.ConnectReply{
			Credentials: &centrifuge.Credentials{
				UserID: userContext.Email,
				Info:   userContextByteArray,
			},
		}, nil
	})

	node.OnConnect(func(client *centrifuge.Client) {
		transport := client.Transport()

		var userContext shared_models.UserContext
		json.Unmarshal(client.Info(), &userContext)

		logrus.Infof("🟩 connected via %s.", transport.Name())

		userBoundChannel := shared_infra_broker.UserBoundChannel{
			Channel: make(chan events.Event),
			User:    userContext,
		}

		internalBroker.AddClient(userBoundChannel)

		go func() {
			for {
				event := <-userBoundChannel.Channel
				evArr, err := json.Marshal(event)

				if err != nil {
					logrus.Warn(err)
				} else {
					client.Send(evArr)
				}
			}
		}()

		client.OnRefresh(func(e centrifuge.RefreshEvent, cb centrifuge.RefreshCallback) {
			logrus.Infof("user %s connection is going to expire, refreshing", client.UserID())
			cb(centrifuge.RefreshReply{
				ExpireAt: time.Now().Unix() + 10,
			}, nil)
		})

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			logrus.Infof("🟨 subscribes on %s", e.Channel)
			cb(centrifuge.SubscribeReply{
				Options: centrifuge.SubscribeOptions{
					EmitPresence: true,
				},
			}, nil)
		})

		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			logrus.Infof("🟦 unsubscribed from %s", e.Channel)
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			internalBroker.RemoveClient(userBoundChannel)
			logrus.Infof("🟥 disconnected, disconnect: %s", e.Disconnect)
		})
	})
}
