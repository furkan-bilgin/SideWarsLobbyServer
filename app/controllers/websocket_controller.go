package controllers

import (
	"encoding/json"
	"sidewarslobby/app/models"
	"sidewarslobby/platform/database"
	"sync"
	"time"

	"github.com/antoniodipinto/ikisocket"
)

var (
	connectedSockets sync.Map
)

func QueueWebsocketNew(kws *ikisocket.Websocket) {
	userToken := kws.Params("token")
	user, err := database.DBQueries.GetUserByToken(userToken)

	// Disconnect if the user failed to authenticate
	if user == nil || err != nil {
		println("QUEUEWEBSOCKET - User failed to authenticate. (token " + userToken + ")")

		kws.Close()
		return
	}

	// Init attributes
	kws.SetAttribute("user", user)
	kws.SetAttribute("isClosed", false)
	userId := user.ID

	mUser := MatchmakingUser{
		UserID:    userId,
		CreatedAt: time.Now(),
		Elo:       user.UserInfo.CachedElo,
	}

	RedisSendJoinQueue(mUser)
	connectedSockets.Store(user.ID, kws)
}

func QueueWebsocketHandleDisconnect(ep *ikisocket.EventPayload) {
	// Return if the user didn't authenticate
	userRaw := ep.Kws.GetAttribute("user")
	if userRaw == nil {
		return
	}

	user := ep.Kws.GetAttribute("user").(*models.User)

	// It's already deleted, that means we don't need to do post-close actions
	if _, contains := connectedSockets.Load(user.ID); !contains {
		return
	}

	connectedSockets.Delete(user.ID)

	// Return if we already did post-close actions
	if ep.Kws.GetAttribute("isClosed") == true {
		return
	}
	println("QUEUEWEBSOCKET - User disconnecting.")

	// Set this to avoid double-closing
	ep.Kws.SetAttribute("isClosed", true)

	// Send matchmaking server to remove this user from the queue
	RedisSendLeaveQueue(user.ID)
}

func QueueWebsocketNewMatch(match *NewMatch) {
	type WebsocketSendQueueItem struct {
		Payload   []byte
		Websocket *ikisocket.Websocket
	}

	sendQueue := make([]WebsocketSendQueueItem, len(match.UserIDs))

	// TODO: send blueteam redteam
	for idx, userId := range match.UserIDs {
		kw, ok := connectedSockets.Load(uint(userId))

		if !ok {
			continue
		}
		kws := kw.(*ikisocket.Websocket)
		user := kws.GetAttribute("user").(*models.User)

		// Create UserMatch
		userMatch := models.UserMatch{
			UserID:       user.ID,
			MatchID:      match.Match.ID,
			UserChampion: user.UserInfo.SelectedChampion,
			TeamID:       match.Teams[int(user.ID)],
		}

		// Insert it to database
		err := database.DBQueries.CreateUserMatch(&userMatch)
		if err != nil {
			panic(err)
		}

		// Remove from connectedSockets so we don't do post-close actions
		connectedSockets.Delete(user.ID)
		payload := struct {
			ServerIP   string
			MatchToken string
		}{ServerIP: "1.game.sw.furkanbilgin.net:9876", MatchToken: JWTCreateUserMatchToken(&userMatch)} // TODO: Change this

		// Add to send queue
		payloadBytes, _ := json.Marshal(payload)
		sendQueue[idx].Payload = payloadBytes
		sendQueue[idx].Websocket = kws
	}

	for _, item := range sendQueue {
		if item.Websocket == nil {
			continue
		}

		item.Websocket.Emit(item.Payload)
		item.Websocket.Close()
	}
}
