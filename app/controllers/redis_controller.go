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
	RedisPairupListener *broadcast.Relay[NewPairup]
)

type MatchmakingUser struct {
	UserID       int
	CreatedAt    time.Time
	Elo          int
	EloPlusMinus int
}

type NewPairup struct {
	MatchID string
	UserID  int
	TeamID  uint8
}

func InitRedisController() {
	RedisPairupListener = broadcast.NewRelay[NewPairup]()
	listenQueueNewPair()
	listenQueueNewMatch()
}

func RedisSendJoinQueue(mUser MatchmakingUser) {
	data, _ := json.Marshal(mUser)
	err := cache.RedisClient.Publish(context.Background(), "queue-add-user", data).Err()
	if err != nil {
		panic(err)
	}
}

func RedisSendLeaveQueue(userID int) {
	data, _ := json.Marshal(MatchmakingUser{UserID: userID})
	err := cache.RedisClient.Publish(context.Background(), "queue-remove-user", data).Err()
	if err != nil {
		panic(err)
	}
}

func listenQueueNewPair() {
	go redisListener("queue-new-pair", func(data []byte, dict map[string]interface{}) {
		var pairUp NewPairup
		json.Unmarshal(data, &pairUp)
	})
}

func listenQueueNewMatch() {
	go redisListener("queue-new-match", func(data []byte, dict map[string]interface{}) {
		payload := struct {
			MatchID string
		}{}
		json.Unmarshal(data, &payload)

		err := database.DBQueries.CreateMatch(&models.Match{MatchmakingID: uuid.MustParse(payload.MatchID)})

		if err != nil {
			panic(err)
		}
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
