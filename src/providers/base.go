package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/http/helper"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerSocialite/v3/src/configs"
	"github.com/ArtisanCloud/PowerSocialite/v3/src/contracts"
	"github.com/ArtisanCloud/PowerSocialite/v3/src/response/wechat"
	"io"
	"net/http"
	"strings"
)

type Base struct {
	ProviderInterface

	state           string
	forcePopup      bool
	config          *configs.Config
	redirectURL     string
	parameters      *object.StringMap
	scopes          []string
	scopeSeparator  string
	httpHelper      *helper.RequestHelper
	guzzleOptions   *object.HashMap
	encodingType    int
	expiresInKey    string
	accessTokenKey  string
	refreshTokenKey string

	TokenFromCode        func(code string) (*object.HashMap, error)
	GetAuthURL           func() (string, error)
	GetTokenURL          func() string
	GetUserByToken       func(token string, openID string) (*object.HashMap, error)
	MapUserToObject      func(user *object.HashMap) *User
	GetAccessToken       func(token string) (contracts.AccessTokenInterface, error)
	BuildAuthURLFromBase func(url string) string
	GetCodeFields        func() *object.StringMap
	GetTokenFields       func(code string) *object.StringMap
}

func NewBase(config *object.HashMap) *Base {

	base := &Base{
		config:          configs.NewConfig(config),
		scopes:          []string{},
		expiresInKey:    "expires_in",
		accessTokenKey:  "access_token",
		refreshTokenKey: "refresh_token",
	}

	// set scopes
	if (*config)["scopes"] != nil {
		base.scopes = (*config)["scopes"].([]string)
	}

	// normalize 'client_id'
	if base.config.Has("client_id") {
		id := base.config.Get("app_id", "").(string)
		if id != "" {
			base.config.Set("client_id", id)
		}
	}

	// normalize 'client_secret'
	if base.config.Has("client_secret") {
		secret := base.config.Get("app_secret", "").(string)
		if secret != "" {
			base.config.Set("client_secret", secret)
		}
	}

	// normalize 'redirect_url'
	if base.config.Has("redirect") {
		redirectURL := base.config.Get("redirect", "").(string)
		base.config.Set("redirect", redirectURL)
	}

	base.OverrideTokenFromCode()

	return base
}

func (base *Base) Redirect(redirectURL string) (string, error) {
	if redirectURL != "" {
		base.WithRedirectURL(redirectURL)
	}

	return base.GetAuthURL()
}

func (base *Base) UserFromCode(code string) (*User, error) {
	tokenResponse, err := base.TokenFromCode(code)
	if err != nil {
		return nil, err
	}

	user, err := base.UserFromToken((*tokenResponse)[base.accessTokenKey].(string), (*tokenResponse)["openid"].(string))
	if err != nil {
		return nil, err
	}

	refreshTokenKey := ""
	if (*tokenResponse)[base.refreshTokenKey] != nil {
		refreshTokenKey = (*tokenResponse)[base.refreshTokenKey].(string)
	}

	expiresInKey := 0.0
	if (*tokenResponse)[base.expiresInKey] != nil {
		expiresInKey = (*tokenResponse)[base.expiresInKey].(float64)
	}

	return user.SetRefreshToken(refreshTokenKey).
		SetExpiresIn(expiresInKey).
		SetTokenResponse(tokenResponse), nil
}

func (base *Base) UserFromToken(token string, openID string) (*User, error) {
	user, err := base.GetUserByToken(token, openID)
	if err != nil {
		return nil, err
	}

	return base.MapUserToObject(user).
		SetProvider(base).
		SetRaw(user).
		SetAccessToken(token), nil
}

func (base *Base) OverrideTokenFromCode() {
	base.TokenFromCode = func(code string) (*object.HashMap, error) {
		//outResponse := &weCom.ResponseTokenFromCode{}
		client, err := base.GetHttpClient()
		if err != nil {
			return nil, err
		}
		response, err := client.Df().
			Url(base.GetTokenURL()).
			Method("POST").
			Json(&object.HashMap{
				"form_params": base.GetTokenFields(code),
				"headers": &object.HashMap{
					"Accept": "application/json",
				},
			}).Request()

		if err != nil {
			return nil, err
		}

		return base.normalizeAccessTokenResponse(response)
	}
}

func (base *Base) refreshToken(refreshToken string) error {
	return errors.New("refreshToken does not support")

}

func (base *Base) WithRedirectURL(redirectURL string) ProviderInterface {
	base.redirectURL = redirectURL

	return base
}

func (base *Base) WithState(state string) ProviderInterface {
	base.state = state

	return base
}

func (base *Base) WithForcePopup(forcePopup bool) ProviderInterface {
	base.forcePopup = forcePopup

	return base
}

func (base *Base) Scopes(scopes []string) *Base {
	base.scopes = scopes

	return base
}

func (base *Base) With(parameters *object.StringMap) *Base {
	base.parameters = parameters

	return base
}

func (base *Base) GetConfig() *configs.Config {
	return base.config
}

func (base *Base) WithScopeSeparator(scopeSeparator string) *Base {
	base.scopeSeparator = scopeSeparator

	return base
}

func (base *Base) GetClientID() string {
	var result string
	if base.config.Get("client_id", "") != nil {
		result = base.config.Get("client_id", "").(string)
	}
	return result
}

func (base *Base) GetClientSecret() string {
	var result string
	if base.config.Get("client_secret", "") != nil {
		result = base.config.Get("client_secret", "").(string)
	}
	return result
}

func (base *Base) GetHttpClient() (*helper.RequestHelper, error) {
	if base.httpHelper != nil {
		return base.httpHelper, nil
	} else {
		return helper.NewRequestHelper(&helper.Config{
			BaseUrl: "",
		})
	}

}

func (base *Base) formatScopes(scopes []string, scopeSeparator string) string {
	return strings.Join(scopes, scopeSeparator)
}

func (base *Base) getTokenFields(code string) *object.StringMap {
	return &object.StringMap{
		"client_id":     base.GetClientID(),
		"client_secret": base.GetClientSecret(),
		"code":          code,
		"redirect_uri":  base.redirectURL,
	}
}

func (base *Base) ParseBody(body io.ReadCloser) (*object.HashMap, error) {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(body)
	jsonHashMap := object.HashMap{}
	err := json.Unmarshal(buf.Bytes(), &jsonHashMap)

	return &jsonHashMap, err
}

func (base *Base) parseAccessToken(body io.ReadCloser) (accessToken contracts.AccessTokenInterface, err error) {

	jsonHashMap, err := base.ParseBody(body)

	if err != nil {
		return nil, err
	}
	return NewAccessToken(jsonHashMap)
}

func (base *Base) buildAuthURLFromBase(url string) string {

	query := object.GetJoinedWithKSort(base.GetCodeFields())

	return url + "?" + query + string(base.encodingType)
}

func (base *Base) getCodeFields() *object.StringMap {
	fields := &object.StringMap{
		"client_id":     base.GetClientID(),
		"redirect_uri":  base.redirectURL,
		"forcePopup":    fmt.Sprintf("%t", base.forcePopup),
		"scope":         base.formatScopes(base.scopes, base.scopeSeparator),
		"response_type": "code",
	}
	fields = object.MergeStringMap(fields, base.parameters)
	if base.state != "" {
		(*fields)["state"] = base.state
	}

	return fields
}

func (base *Base) normalizeAccessTokenResponse(response *http.Response) (*object.HashMap, error) {

	token := wechat.ResponseAuthenticatedAccessToken{}

	body := response.Body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(body)
	err := json.Unmarshal(buf.Bytes(), &token)
	if err != nil {
		return nil, err
	}

	mapToken, err := object.StructToHashMap(token)

	return mapToken, err

}
