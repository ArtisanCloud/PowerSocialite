package session

import (
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

const testKey = "test_key"

func TestCookieSession_StoreAndRetrieve(t *testing.T) {
	// 初始化 CookieSession
	cookieSession := NewCookieSession(testKey, &sessions.Options{})
	if !cookieSession.keySet {
		t.Skip("SESSION_SECRET environment variable is not set")
	}

	// 模拟 HTTP 请求和响应
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 存储值
	key := "state"
	value := "xyz123"
	err := cookieSession.StoreInSession(key, value, req, w)
	assert.NoError(t, err, "Failed to store value in session")

	// 模拟一次新的请求，获取存储的值
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", w.Result().Header.Get("Set-Cookie")) // 模拟之前的 Cookie

	// 从 session 中获取存储的值
	retrievedValue, err := cookieSession.GetFromSession(key, req)
	println(value, retrievedValue)
	assert.NoError(t, err, "Failed to retrieve value from session")
	assert.Equal(t, value, retrievedValue, "Retrieved value does not match expected value")
}

func TestCookieSession_GetFromSession_KeyNotFound(t *testing.T) {
	cookieSession := NewCookieSession(testKey, &sessions.Options{})
	if !cookieSession.keySet {
		t.Skip("SESSION_SECRET environment variable is not set")
	}

	// 模拟 HTTP 请求
	req := httptest.NewRequest("GET", "/test", nil)

	// 获取一个不存在的 key
	_, err := cookieSession.GetFromSession("nonexistent_key", req)
	assert.Error(t, err, "Expected error when key is not found in session")
}

func TestCookieSession_StoreWithCompression(t *testing.T) {
	cookieSession := NewCookieSession(testKey, &sessions.Options{})
	if !cookieSession.keySet {
		t.Skip("SESSION_SECRET environment variable is not set")
	}

	// 模拟 HTTP 请求和响应
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 存储一个带有压缩值的 key/value
	key := "session_key"
	value := "some_value_to_compress"
	err := cookieSession.StoreInSession(key, value, req, w)
	assert.NoError(t, err, "Failed to store compressed value in session")

	// 模拟新的请求，获取存储的压缩值
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", w.Result().Header.Get("Set-Cookie")) // 模拟之前的 Cookie

	// 获取压缩后的值
	retrievedValue, err := cookieSession.GetFromSession(key, req)
	assert.NoError(t, err, "Failed to retrieve compressed value from session")
	assert.Equal(t, value, retrievedValue, "Retrieved compressed value does not match expected value")
}

func TestCookieSession_StoreAndRetrieve_WithEmptyValue(t *testing.T) {
	cookieSession := NewCookieSession(testKey, &sessions.Options{})
	if !cookieSession.keySet {
		t.Skip("SESSION_SECRET environment variable is not set")
	}

	// 模拟 HTTP 请求和响应
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 存储一个空值
	key := "empty_key"
	value := ""
	err := cookieSession.StoreInSession(key, value, req, w)
	assert.NoError(t, err, "Failed to store empty value in session")

	// 模拟新的请求，获取存储的空值
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", w.Result().Header.Get("Set-Cookie")) // 模拟之前的 Cookie

	// 获取空值
	retrievedValue, err := cookieSession.GetFromSession(key, req)
	assert.NoError(t, err, "Failed to retrieve empty value from session")
	assert.Equal(t, value, retrievedValue, "Retrieved empty value does not match expected value")
}
