package session

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strings"
)

const SessionName = "_media_x_session"

type CookieSession struct {
	Store  sessions.Store
	keySet bool
}

func NewCookieSession(sessionSecret string, Options *sessions.Options) *CookieSession {
	key := []byte(os.Getenv("SESSION_SECRET"))
	if sessionSecret != "" {
		key = []byte(sessionSecret)
	}
	store := sessions.NewCookieStore(key)
	store.Options = Options
	return &CookieSession{
		Store:  store,
		keySet: len(key) != 0,
	}
}

// StoreInSession stores a specified key/value pair in the session.
func (s *CookieSession) StoreInSession(key string, value string, req *http.Request, res http.ResponseWriter) error {
	session, _ := s.Store.New(req, SessionName)

	if err := s.UpdateSessionValue(session, key, value); err != nil {
		return err
	}

	return session.Save(req, res)
}

// GetFromSession retrieves a previously-stored value from the session.
// If no value has previously been stored at the specified key, it will return an error.
func (s *CookieSession) GetFromSession(key string, req *http.Request) (string, error) {
	session, err := s.Store.Get(req, SessionName)
	value, err := s.GetSessionValue(session, key)
	if err != nil {
		return "", errors.New("could not find a matching session for this request")
	}

	return value, nil
}

func (s *CookieSession) GetSessionValue(session *sessions.Session, key string) (string, error) {
	value := session.Values[key]
	if value == nil {
		return "", fmt.Errorf("could not find a matching session for this request")
	}

	rdata := strings.NewReader(value.(string))
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return "", err
	}
	res, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (s *CookieSession) UpdateSessionValue(session *sessions.Session, key, value string) error {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(value)); err != nil {
		return err
	}
	if err := gz.Flush(); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	session.Values[key] = b.String()
	return nil
}

func (s *CookieSession) ClearSession(key string, req *http.Request, res http.ResponseWriter) error {
	sess, err := s.Store.Get(req, SessionName)
	if err != nil {
		return err
	}
	sess.Values[key] = nil

	err = sess.Save(req, res)
	if err != nil {
		return errors.New("Could not delete user session ")
	}
	return nil
}
