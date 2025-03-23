package providers

import (
	"encoding/json"
	"github.com/ArtisanCloud/PowerSocialite/v4/contracts"
	"github.com/pkg/errors"
	"time"
)

type BaseSession struct {
	contracts.ISession
	AuthUrl         string    `json:"authUrl" yaml:"auth_url"`
	AccessToken     string    `json:"accessToken" yaml:"access_token"`
	RefreshToken    string    `json:"refreshToken" yaml:"refresh_token"`
	ExpiresAt       time.Time `json:"expiresAt" yaml:"expires_at"`
	Scope           string    `json:"scope" yaml:"scope"`
	TokenType       string    `json:"tokenType" yaml:"token_type"`
	ExpiresInKey    string    `json:"expiresInKey" yaml:"expires_in_key"`
	AccessTokenKey  string    `json:"accessTokenKey" yaml:"access_token_key"`
	RefreshTokenKey string    `json:"refreshTokenKey" yaml:"refresh_token_key"`
}

const NoAuthUrlErrorMessage = "an AuthUrl has not been set"

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the Google provider.
func (s BaseSession) GetAuthURL() (string, error) {
	if s.AuthUrl == "" {
		return "", errors.New(NoAuthUrlErrorMessage)
	}
	return s.AuthUrl, nil
}

// Marshal the session into a string
func (s BaseSession) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}
