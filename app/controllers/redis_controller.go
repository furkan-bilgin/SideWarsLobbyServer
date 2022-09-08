package controllers

import (
	"context"
	"encoding/json"
	"sidewarslobby/platform/cache"
	"time"

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
}

func InitRedisController() {
	RedisPairupListener = broadcast.NewRelay[NewPairup]()
	listenQueueNewPair()
}

func RedisSendJoinQueue(mUser MatchmakingUser) {
	data, _ := json.Marshal(mUser)
	cache.RedisClient.Publish(context.TODO(), "queue-add-user", data)
}

func RedisSendLeaveQueue(userID int) {
	data, _ := json.Marshal(MatchmakingUser{UserID: userID})
	cache.RedisClient.Publish(context.TODO(), "queue-remove-user", data)
}

func listenQueueNewPair() {
	go redisListener("queue-new-pair", func(data []byte, dict map[string]interface{}) {
		var pairUp NewPairup
		json.Unmarshal(data, &pairUp)

		RedisPairupListener.Broadcast(pairUp)
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
