package session

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 测试 RedisSession 的 StoreInSession 和 GetFromSession
func TestRedisSession(t *testing.T) {
	// 配置 Redis 客户端
	redisOptions := &redis.Options{
		Addr: "localhost:6379", // Redis 服务地址
	}

	// 创建 Redis 会话管理器
	redisSession := NewRedisSession(redisOptions)

	// 创建一个模拟 HTTP 请求和响应
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// 模拟一个 session_id cookie
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test_session_id"})

	// 创建响应记录器
	res := httptest.NewRecorder()

	// 测试 StoreInSession
	err := redisSession.StoreInSession("username", "testUser", req, res)
	assert.NoError(t, err, "StoreInSession should not return an error")

	// 测试 GetFromSession
	value, err := redisSession.GetFromSession("username", req)
	assert.NoError(t, err, "GetFromSession should not return an error")
	assert.Equal(t, "testUser", value, "Stored and retrieved values should be the same")

	// 测试 GetFromSession 对不存在的 key 返回的错误
	_, err = redisSession.GetFromSession("nonexistent_key", req)
	assert.Error(t, err, "GetFromSession should return an error when the key doesn't exist")
}

// 测试当 session ID 为空时，StoreInSession 和 GetFromSession 的行为
func TestSessionIDEmpty(t *testing.T) {
	// 配置 Redis 客户端
	redisOptions := &redis.Options{
		Addr: "localhost:6379", // Redis 服务地址
	}

	// 创建 Redis 会话实例
	redisSession := NewRedisSession(redisOptions)

	// 创建没有 session_id cookie 的请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	res := httptest.NewRecorder()

	// 测试当 session ID 为空时，StoreInSession 应返回错误
	err := redisSession.StoreInSession("username", "testUser", req, res)
	assert.Error(t, err, "StoreInSession should return an error when session ID is empty")

	// 测试当 session ID 为空时，GetFromSession 应返回错误
	_, err = redisSession.GetFromSession("username", req)
	assert.Error(t, err, "GetFromSession should return an error when session ID is empty")
}

// 测试 RedisSession 是否处理 Redis 错误
func TestRedisSessionErrorHandling(t *testing.T) {
	// 使用无效的 Redis 地址来模拟 Redis 错误
	redisOptions := &redis.Options{
		Addr: "localhost:9999", // 使用无效的端口
	}

	redisSession := NewRedisSession(redisOptions)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test_session_id"})
	res := httptest.NewRecorder()

	// 测试 StoreInSession 在 Redis 错误时的行为
	err := redisSession.StoreInSession("username", "testUser", req, res)
	assert.Error(t, err, "StoreInSession should return an error when Redis is unreachable")

	// 测试 GetFromSession 在 Redis 错误时的行为
	_, err = redisSession.GetFromSession("username", req)
	assert.Error(t, err, "GetFromSession should return an error when Redis is unreachable")
}
