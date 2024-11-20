package dbservice

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisDict struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisDict initializes a new RedisDict instance
func NewRedisDict(host string, port int, db int) *RedisDict {
	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":" + strconv.Itoa(port),
		DB:   db,
	})
	return &RedisDict{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Get retrieves and converts a value from Redis by key
func (r *RedisDict) Get(key string) (interface{}, error) {
	value, err := r.client.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Attempt to convert the value to the appropriate type
	if value == "true" {
		return true, nil
	} else if value == "false" {
		return false, nil
	} else if value == "null" {
		return nil, nil
	}

	// Try to convert to float or integer
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue, nil
		}
		return floatValue, nil
	}

	// Try to parse JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(value), &jsonData); err == nil {
		return jsonData, nil
	}

	// If all conversions fail, return the original string
	return value, nil
}

// Set stores a key-value pair in Redis, handling different types
func (r *RedisDict) Set(key string, value interface{}) error {
	var stringValue string

	switch v := value.(type) {
	case bool:
		if v {
			stringValue = "true"
		} else {
			stringValue = "false"
		}
	case nil:
		stringValue = "null"
	case int:
		stringValue = strconv.Itoa(v)
	case int32:
		stringValue = strconv.Itoa(int(v))
	case int64:
		stringValue = strconv.FormatInt(v, 10)
	case float32:
		stringValue = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		stringValue = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		stringValue = v
	default:
		// Marshal to JSON if it's a complex type (e.g., map or list)
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return err
		}
		stringValue = string(jsonBytes)
	}

	return r.client.Set(r.ctx, key, stringValue, 0).Err()
}

// Delete removes a key from Redis
func (r *RedisDict) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisDict) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

// Keys retrieves all keys in Redis
func (r *RedisDict) Keys() ([]string, error) {
	return r.client.Keys(r.ctx, "*").Result()
}

// Pop retrieves a value by key and deletes the key from Redis
func (r *RedisDict) Pop(key string) (interface{}, error) {
	value, err := r.Get(key)
	if err != nil {
		return nil, err
	}
	if _, err := r.client.Del(r.ctx, key).Result(); err != nil {
		return nil, err
	}
	return value, nil
}
