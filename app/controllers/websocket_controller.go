package controllers

import (
	"encoding/json"
	"sidewarslobby/app/models"
	"sidewarslobby/platform/database"
	"time"

	"github.com/antoniodipinto/ikisocket"
	"github.com/google/uuid"
)

func QueueWebsocketNew(kws *ikisocket.Websocket) {
	userToken := kws.Params("token")
	user := database.DBQueries.GetUserByToken(userToken)

	// Disconnect if the user failed to authenticate
	if user == nil {
		println("QUEUEWEBSOCKET - User failed to authenticate. (token " + userToken + ")")

		kws.Close()
		return
	}

	// Init cancel signal
	cancelGoroutine := make(chan bool)

	// Init attributes
	kws.SetAttribute("user", *user)
	kws.SetAttribute("isClosed", false)
	kws.SetAttribute("cancelGoroutine", cancelGoroutine)
	userId := user.ID

	mUser := MatchmakingUser{
		UserID:    int(userId),
		CreatedAt: time.Now(),
		Elo:       user.CachedElo,
	}

	RedisSendJoinQueue(mUser)

	go (func() {
		for {
			select {
			case <-cancelGoroutine:
				return
			default:
				// Listen pairups
				l := RedisPairupListener.Listener(10)
				for pairUp := range l.Ch() {
					// If this user paired up...
					if pairUp.UserID == int(userId) {
						// Create UserMatch
						userMatch := models.UserMatch{
							Token:         uuid.New().String(),
							UserID:        userId,
							MatchmakingID: uuid.MustParse(pairUp.MatchID),
							UserChampion:  user.UserInfo.SelectedChampion,
							TeamID:        pairUp.TeamID,
						}

						err := database.DBQueries.CreateUserMatch(&userMatch)
						if err != nil {
							panic(err)
						}

						// Send payload to WebSocket client
						payload := struct {
							ServerIP   string
							MatchToken string
						}{ServerIP: "1.game.sw.furkanbilgin.net:9876", MatchToken: userMatch.Token} // TODO: Change this

						payloadBytes, _ := json.Marshal(payload)
						kws.Emit(payloadBytes)

						// Close connection without post-close actions
						kws.SetAttribute("isClosed", true)
					}
				}
			}
		}
	})()
}

func QueueWebsocketHandleDisconnect(ep *ikisocket.EventPayload) {
	// Return if the user didn't authenticate
	userRaw := ep.Kws.GetAttribute("user")
	if userRaw == nil {
		return
	}

	// Return if we already did post-close actions
	if ep.Kws.GetAttribute("isClosed") == true {
		return
	}

	// Set this to avoid double-closing
	ep.Kws.SetAttribute("isClosed", true)

	println("QUEUEWEBSOCKET - User disconnecting.")

	// Send matchmaking server to remove this user from the queue
	user := ep.Kws.GetAttribute("user").(models.User)
	RedisSendLeaveQueue(int(user.ID))

	// Cancel Redis subscription goroutine
	cancelGoroutine := ep.Kws.GetAttribute("cancelGoroutine").(chan bool)
	cancelGoroutine <- true
}
