package session

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type RedisSession struct {
	client *redis.Client
}

// NewRedisSession 创建一个新的 RedisSession 实例
func NewRedisSession(config *redis.Options) *RedisSession {
	// 配置 Redis 客户端
	client := redis.NewClient(config)

	return &RedisSession{
		client: client,
	}
}

// StoreInSession 将指定的 key/value 对存储到 Redis 会话中
func (r *RedisSession) StoreInSession(key string, value string, req *http.Request, res http.ResponseWriter) error {
	redisKey := r.GetSessionKey(key)

	// 将数据存储到 Redis 中，设置超时为 5 分钟
	err := r.client.Set(context.Background(), redisKey, value, 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetFromSession 从 Redis 会话中检索指定的 key 对应的值
func (r *RedisSession) GetFromSession(key string, req *http.Request) (string, error) {
	redisKey := r.GetSessionKey(key)
	// 从 Redis 中获取数据
	value, err := r.client.Get(context.Background(), redisKey).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found in session")
	}
	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *RedisSession) ClearSession(key string, req *http.Request, res http.ResponseWriter) error {
	redisKey := fmt.Sprintf("session:%s:%s", SessionName, key)
	s.client.Del(context.Background(), redisKey)

	return nil
}

func (r *RedisSession) GetSessionKey(key string) string {
	// 使用 sessionID 和 key 作为符合键
	redisKey := fmt.Sprintf("session:%s:%s", SessionName, key)
	return redisKey
}
