package providers

import (
	"errors"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerSocialite/v3/src/contracts"
)

type AccessToken struct {
	contracts.AccessTokenInterface
	*object.Attribute
}

func NewAccessToken(attributes *object.HashMap) (*AccessToken, error) {

	if (*attributes)["access_token"] == nil {
		return nil, errors.New("The key 'access_token' could not be empty.")
	}

	accessToken := &AccessToken{
		Attribute: object.NewAttribute(attributes),
	}
	return accessToken, nil

}

func (accessToken *AccessToken) GetToken() string {
	return accessToken.GetString("access_token", "")
}

func (accessToken *AccessToken) GetRefreshToken() string {
	return accessToken.GetString("refresh_token", "")
}

func (accessToken *AccessToken) SetRefreshToken(token string) {
	accessToken.SetAttribute("refresh_token", token)
}

func (accessToken *AccessToken) JsonSerialize() string {
	return accessToken.GetToken()
}
