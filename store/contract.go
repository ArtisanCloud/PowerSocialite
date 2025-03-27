package session

import (
	"net/http"
)

type IStore interface {
	// StoreInSession 将指定的 key/value 对存储到会话中
	StoreInSession(key string, value string, req *http.Request, res http.ResponseWriter) error

	// GetFromSession 从会话中检索指定的 key 对应的值
	GetFromSession(key string, req *http.Request) (string, error)

	// ClearSession 清除会话
	ClearSession(key string, req *http.Request, res http.ResponseWriter) error
}
