package controllers

import (
	"context"
	"encoding/json"
	"sidewarslobby/app/models"
	"sidewarslobby/platform/cache"
	"sidewarslobby/platform/database"
	"time"

	"github.com/google/uuid"
	"github.com/teivah/broadcast"
)

var (
	RedisNewMatchListener *broadcast.Relay[*NewMatch]
)

type MatchmakingUser struct {
	UserID       uint
	CreatedAt    time.Time
	Elo          int
	EloPlusMinus int
}

type NewMatch struct {
	MatchmakingID string
	UserIDs       []int
	Teams         map[int]uint8
	Match         *models.Match
}

func InitRedisController() {
	RedisNewMatchListener = broadcast.NewRelay[*NewMatch]()
	listenQueueNewMatch()
}

func RedisSendJoinQueue(mUser MatchmakingUser) {
	data, _ := json.Marshal(mUser)
	err := cache.RedisClient.Publish(context.Background(), "queue-add-user", data).Err()
	if err != nil {
		panic(err)
	}
}

func RedisSendLeaveQueue(userID uint) {
	data, _ := json.Marshal(MatchmakingUser{UserID: userID})
	err := cache.RedisClient.Publish(context.Background(), "queue-remove-user", data).Err()
	if err != nil {
		panic(err)
	}
}

func listenQueueNewMatch() {
	go redisListener("queue-new-match", func(data []byte, dict map[string]interface{}) {
		// Load NewMatch from Redis pubsub
		payload := &NewMatch{}
		json.Unmarshal(data, payload)
		// Create a new Match
		match, err := database.DBQueries.FindOrCreateMatch(&models.Match{MatchmakingID: uuid.MustParse(payload.MatchmakingID)})
		if err != nil {
			panic(err)
		}
		payload.Match = match
		// Broadcast it
		RedisNewMatchListener.Notify(payload)
		QueueWebsocketNewMatch(payload)
	})
}

func redisListener(channel string, callback func([]byte, map[string]interface{})) {
	subscriber := cache.RedisClient.Subscribe(context.Background(), channel)

	for {
		var data map[string]interface{}
		msg, err := subscriber.ReceiveMessage(context.Background())

		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
			panic(err)
		}

		callback([]byte(msg.Payload), data)
	}
}
