package google

import (
	"encoding/json"
	"github.com/ArtisanCloud/PowerSocialite/v4/contracts"
	"github.com/ArtisanCloud/PowerSocialite/v4/kernel"
	"github.com/ArtisanCloud/PowerSocialite/v4/providers"
	"github.com/pkg/errors"
	"strings"
)

type Session struct {
	providers.BaseSession
	IDToken string
}

// Authorize the session with Google and return the access token to be stored for future use.
func (s *Session) Authorize(provider contracts.IProvider, params contracts.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.OAuthConfig.Exchange(kernel.ContextForClient(p.Client()), params.Get("code"))
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("Invalid token received from provider")
	}

	s.AccessToken = token.AccessToken
	s.RefreshToken = token.RefreshToken
	s.ExpiresAt = token.Expiry
	if idToken := token.Extra("id_token"); idToken != nil {
		s.IDToken = idToken.(string)
	}
	return token.AccessToken, err
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (s *Session) UnmarshalSession(data string) (contracts.ISession, error) {
	sess := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(sess)
	return sess, err
}
