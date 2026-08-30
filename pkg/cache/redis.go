package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func InitRedis() *redis.Client {
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", os.Getenv("REDIS_HOST"), port),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	err := client.Ping(Ctx).Err()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return client
}
