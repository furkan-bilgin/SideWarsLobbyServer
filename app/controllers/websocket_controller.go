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
	goroutineDone := make(chan bool)

	// Init attributes
	kws.SetAttribute("user", *user)
	kws.SetAttribute("isClosed", false)
	kws.SetAttribute("cancelGoroutine", cancelGoroutine)
	userId := user.ID

	mUser := MatchmakingUser{
		UserID:    userId,
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
				// Listen matches
				l := RedisNewMatchListener.Listener(1)

				// Received a new match
				for match := range l.Ch() {
					for _, otherUserID := range match.UserIDs {
						if otherUserID == int(userId) {
							// Create UserMatch
							userMatch := models.UserMatch{
								Token:        uuid.New().String(),
								UserID:       userId,
								MatchID:      match.Match.ID,
								UserChampion: user.UserInfo.SelectedChampion,
								TeamID:       match.Teams[otherUserID],
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
							// Make the main-thread set attribute and close the connection, because it somehow causes a deadlock when we do it inside a goroutine.
							goroutineDone <- true
							return
						}
					}
				}
			}
		}
	})()

	<-goroutineDone
	kws.SetAttribute("isClosed", true)
	kws.Close()
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
	RedisSendLeaveQueue(user.ID)

	// Cancel Redis subscription goroutine
	cancelGoroutine := ep.Kws.GetAttribute("cancelGoroutine").(chan bool)
	cancelGoroutine <- true
}
