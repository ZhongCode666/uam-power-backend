package dbservice

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"uam-power-backend/utils"
)

// RedisDict 结构体表示一个 Redis 字典
type RedisDict struct {
	client *redis.Client   // Redis 客户端
	ctx    context.Context // 上下文
}

// NewRedisDict 初始化一个新的 RedisDict 实例
func NewRedisDict(host string, port int, db int) *RedisDict {
	// 创建一个新的 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:        host + ":" + strconv.Itoa(port), // Redis 地址
		DB:          db,                              // Redis 数据库编号
		IdleTimeout: 5 * time.Second,                 // 空闲超时时间
	})
	// 返回 RedisDict 实例
	return &RedisDict{
		client: rdb,                  // Redis 客户端
		ctx:    context.Background(), // 上下文
	}
}

// Get 从 Redis 中根据键检索并转换值
func (r *RedisDict) Get(key string) (interface{}, error) {
	value, err := r.client.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 如果键不存在，返回 nil
		return nil, nil
	} else if err != nil {
		utils.MsgError("        [RedisGet]get key failed: >" + err.Error())
		// 如果发生其他错误，返回错误
		return nil, err
	}

	// 将字符串值转换为适当的类型
	return ConvertStringToInterface(value)
}

// ConvertStringToInterface 将字符串值转换为适当的类型
func ConvertStringToInterface(value string) (interface{}, error) {
	// 尝试将值转换为适当的类型
	if value == "true" {
		return true, nil
	} else if value == "false" {
		return false, nil
	} else if value == "null" {
		return nil, nil
	}

	// 尝试转换为浮点数或整数
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue, nil
		}
		return floatValue, nil
	}

	// 尝试解析 JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(value), &jsonData); err == nil {
		return jsonData, nil
	}
	return value, nil
}

// Set 将键值对存储在 Redis 中，处理不同类型
func (r *RedisDict) Set(key string, value interface{}) error {
	var stringValue string

	// 根据值的类型进行处理
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
		// 如果是复杂类型（例如 map 或 list），则将其序列化为 JSON
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			utils.MsgError("        [RedisSet]marshal value failed: >" + err.Error())
			return err
		}
		stringValue = string(jsonBytes)
	}

	// 将键值对存储在 Redis 中
	return r.client.Set(r.ctx, key, stringValue, 0).Err()
}

// SetWithDuration 将键值对存储在 Redis 中，并设置过期时间
func (r *RedisDict) SetWithDuration(key string, value interface{}, timeout int) error {
	var stringValue string

	// 根据值的类型进行处理
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
		// 如果是复杂类型（例如 map 或 list），则将其序列化为 JSON
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			utils.MsgError("        [RedisSetWithDuration]marshal value failed: >" + err.Error())
			return err
		}
		stringValue = string(jsonBytes)
	}

	// 将键值对存储在 Redis 中，并设置过期时间
	return r.client.Set(r.ctx, key, stringValue, time.Duration(timeout)*time.Second).Err()
}

// Delete 从 Redis 中删除一个键
func (r *RedisDict) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists 检查 Redis 中是否存在一个键
func (r *RedisDict) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

// Keys 检索 Redis 中的所有键
func (r *RedisDict) Keys() ([]string, error) {
	return r.client.Keys(r.ctx, "*").Result()
}

// GetVals 从 Redis 中检索多个键的值，并将其转换为适当的类型
func (r *RedisDict) GetVals(keys []string) ([]interface{}, error) {
	// 使用 MGet 命令从 Redis 中获取多个键的值
	ReInterface, err := r.client.MGet(r.ctx, keys...).Result()
	if err != nil {
		utils.MsgError("        [RedisGetVals]get keys failed: >" + err.Error())
		return nil, err
	}
	// 遍历结果并将每个值转换为适当的类型
	for i, v := range ReInterface {
		if v == nil {
			continue
		}
		ReInterface[i], _ = ConvertStringToInterface(v.(string))
	}
	return ReInterface, nil
}

// Pop 从 Redis 中检索一个键的值并删除该键
func (r *RedisDict) Pop(key string) (interface{}, error) {
	// 获取键的值
	value, err := r.Get(key)
	if err != nil {
		utils.MsgError("        [RedisPop]get key failed: >" + err.Error())
		return nil, err
	}
	// 删除键
	if _, err := r.client.Del(r.ctx, key).Result(); err != nil {
		utils.MsgError("        [RedisPop]delete key failed: >" + err.Error())
		return nil, err
	}
	// 返回值
	return value, nil
}
