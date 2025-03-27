package google

import (
	"github.com/ArtisanCloud/PowerSocialite/v4/auth"
	"github.com/ArtisanCloud/PowerSocialite/v4/contracts"
	"github.com/ArtisanCloud/PowerSocialite/v4/providers"
	"github.com/pkg/errors"
)

type Session struct {
	providers.BaseSession
	IDToken string
}

// Authorize the session with Google and return the access token to be stored for future use.
func (s *Session) Authorize(provider contracts.IProvider, params contracts.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.OAuthConfig.Exchange(auth.ContextForClient(p.Client()), params.Get("code"))
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
