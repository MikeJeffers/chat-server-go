package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       0,
		Protocol: 3,
	})
}

type UserTokenData struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

func getToken(userId int64, redisClient *redis.Client) *UserTokenData {
	value, err := redisClient.Get(context.TODO(), fmt.Sprintf("jwt:%v", userId)).Bytes()
	if err != nil {
		return nil
	}
	data := UserTokenData{}
	if json.Unmarshal(value, &data) != nil {
		return nil
	}
	return &data
}
