package bootstrap

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func NewRedisClient(env *Env) *redis.Client {
	redisHost := env.RedisHost
	redisPort := env.RedisPort
	redisPassword := env.RedisPassword
	redisDB := env.RedisDB

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	log.Infof("Connecting to Redis: %s DB=%d", addr, redisDB)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Errorf("Failed to connect to Redis: %v", err)
		log.Fatal(err)
	}

	log.Info("Successfully connected to Redis")
	return client
}

func CloseRedisConnection(client *redis.Client) {
	if client == nil {
		return
	}

	if err := client.Close(); err != nil {
		log.Error("Error closing Redis connection: ", err)
	}

	log.Info("Connection to Redis closed.")
}
