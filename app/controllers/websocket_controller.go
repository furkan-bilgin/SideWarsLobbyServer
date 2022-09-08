package controllers

import (
	"encoding/json"
	"sidewarslobby/app/models"
	"sidewarslobby/platform/database"

	"github.com/antoniodipinto/ikisocket"
)

func QueueWebsocketNew(kws *ikisocket.Websocket) {
	userToken := kws.Params("token")
	user := database.DBQueries.GetUserByToken(userToken)
	if user == nil {
		println("QUEUEWEBSOCKET - User failed to authenticate. (token " + userToken + ")")

		kws.Close()
		return
	}

	kws.SetAttribute("user", *user)
	userId := user.ID

	go (func() {
		for {
			l := RedisPairupListener.Listener(1)

			for pairUp := range l.Ch() {
				if pairUp.UserID == int(userId) {
					payload := struct {
					}{}
					payloadBytes, _ := json.Marshal(payload)
					kws.Emit(payloadBytes)
				}
			}
		}
	})()
}

func QueueWebsocketHandleDisconnect(ep *ikisocket.EventPayload) {
	println("QUEUEWEBSOCKET - User disconnecting.")
	userRaw := ep.Kws.GetAttribute("user")
	if userRaw == nil {
		return
	}

	user := ep.Kws.GetAttribute("user").(models.User)
	RedisSendLeaveQueue(int(user.ID))
}
