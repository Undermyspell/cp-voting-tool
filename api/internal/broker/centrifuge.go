package broker

import (
	"context"
	"encoding/json"
	"time"
	"voting/internal/models"

	"github.com/centrifugal/centrifuge"
	"github.com/sirupsen/logrus"
)

func handleLog(e centrifuge.LogEntry) {
	logrus.Infof("%s: %v", e.Message, e.Fields)
}

var node *centrifuge.Node

func initHandlers() {
	node.OnConnecting(func(ctx context.Context, event centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		userContext, err := models.GetUserContextFromToken(event.Token)

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

		var userContext models.UserContext
		json.Unmarshal(client.Info(), &userContext)

		logrus.Infof("🟩 user %s with usercontext %s connected via %s.", client.UserID(), userContext, transport.Name())

		// userBoundChannel := UserBoundChannel{
		// 	Channel: make(chan events.Event),
		// 	User:    *userContext,
		// }

		client.OnRefresh(func(e centrifuge.RefreshEvent, cb centrifuge.RefreshCallback) {
			logrus.Infof("user %s connection is going to expire, refreshing", client.UserID())
			cb(centrifuge.RefreshReply{
				ExpireAt: time.Now().Unix() + 10,
			}, nil)
		})

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			logrus.Infof("🟨 user %s subscribes on %s", client.UserID(), e.Channel)
			cb(centrifuge.SubscribeReply{
				Options: centrifuge.SubscribeOptions{
					EmitPresence: true,
				},
			}, nil)
		})

		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			logrus.Infof("🟦 user %s unsubscribed from %s", client.UserID(), e.Channel)
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			logrus.Infof("🟥 user %s disconnected, disconnect: %s", client.UserID(), e.Disconnect)
		})
	})

}
