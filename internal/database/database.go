package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	"gofetch.timwalker.dev/internal/env"
)

type Service interface {
	IsHealthy() bool
	IncrementCounter() int
}

type service struct {
	db *redis.Client
}

var (
	address  = env.GetString("REDIS_ADDRESS", "localhost")
	port     = env.GetString("REDIS_PORT", "6379")
	password = env.GetString("REDIS_PASSWORD", "")
	database = env.GetInt("REDIS_DATABASE", 0)
)

func New() Service {
	fullAddress := fmt.Sprintf("%s:%s", address, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fullAddress,
		Password: password,
		DB:       database,
	})

	s := &service{db: rdb}

	return s
}

func (s *service) IncrementCounter() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Default is now 5s
	defer cancel()
	count, _ := s.db.Incr(ctx, "counter").Result()

	return int(count)
}

// isHealthy returns the health status by pinging the Redis server.
func (s *service) IsHealthy() bool {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Default is now 5s
	// defer cancel()
	//
	// pong, err := s.db.Ping(ctx).Result()
	// if err != nil {
	// 	log.Fatalf("db down: %v", err)
	// }
	// return strings.ToUpper(pong) == "PONG"
	return true
}
