package main

import (
	"github.com/sirupsen/logrus"
	"github.com/redis/go-redis/v9"
	"context"
	"encoding/json"	
)

// 定义一个 User 结构体
type User struct {
	ID   uint `json:"id"`
	Name string `json:"name"`
	Age  int `json:"age"`
}

var log = logrus.New()


// Initialize Redis client
func initRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis address
		DB:   0,                // use default DB
	})
	return client
}

// Fetch user information from Redis
func getUserFromRedis(client *redis.Client, userID string) (*User, error) {
	var ctx = context.Background()
	val, err := client.Get(ctx, userID).Result()
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func main() {
	log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })

	client := initRedisClient()
	userID := "user:1"

	user, err := getUserFromRedis(client, userID)
	if err != nil {
		log.Error(err)
	} 

	log.Info("User information: ", *user)

}