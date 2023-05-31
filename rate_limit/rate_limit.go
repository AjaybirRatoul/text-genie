package rate_limit

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter() *RateLimiter {
	// Initialize and configure the Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "",
		DB:       0, // Use default Redis database
	})

	// Test the connection to ensure Redis is running
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return &RateLimiter{
		client: client,
	}
}

func (rl *RateLimiter) CheckRateLimit(key string, limit int, duration time.Duration) bool {
	// Generate a unique key based on the provided key
	rateLimitKey := fmt.Sprintf("ratelimit:%s", key)

	// Get the current count of requests for the key
	count, err := rl.client.Incr(rateLimitKey).Result()
	if err != nil {
		fmt.Println("Redis error:", err)
		return false
	}

	// Set the expiration for the key if it doesn't exist
	if count == 1 {
		err = rl.client.Expire(rateLimitKey, duration).Err()
		if err != nil {
			fmt.Println("Redis error:", err)
			return false
		}
	}

	log.Println("count: ", count)

	// Check if the count exceeds the limit
	if count > int64(limit) {
		return false
	}

	return true
}
