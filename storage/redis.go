package storage

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client  *redis.Client
	timeout time.Duration
}

var RedisStorageClient *RedisStorage

// Setup the redis client
func SetupRedisStorageClient() error {
	// Get env variables and parse them
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	timeout := os.Getenv("REDIS_TIMEOUT")
	if timeout == "" {
		timeout = "5s"
	}
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		return err
	}

	// Create redis client
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Set redis storage client
	RedisStorageClient = &RedisStorage{
		client:  client,
		timeout: timeoutDuration,
	}
	return nil
}

// Function to set pack sizes in redis
func (r *RedisStorage) SetPackSizes(ctx context.Context, sizes []int) error {
	data, err := json.Marshal(sizes)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "default_packs", data, 0).Err()
}

// Function to get pack sizes from redis
func (r *RedisStorage) GetPackSizes() ([]int, error) {
	// Build context with redis timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get pack sizes
	val, err := r.client.Get(ctx, "default_packs").Result()
	if err != nil {
		// If key does not exist return empty array
		if err == redis.Nil {
			return []int{}, nil
		}
		return nil, err
	}

	// Unmarshal json
	var sizes []int
	if err := json.Unmarshal([]byte(val), &sizes); err != nil {
		return nil, err
	}

	return sizes, nil
}
