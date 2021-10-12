package providers

import (
	"errors"
	contract2 "github.com/ArtisanCloud/go-libs/http/contract"
	"github.com/ArtisanCloud/go-libs/http/request"
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/ArtisanCloud/go-socialite/src"
	"github.com/ArtisanCloud/go-socialite/src/response/weCom"
	"strings"
)

type Base struct {
	src.ProviderInterface

	state           string
	config          *src.Config
	redirectURL     string
	parameters      *object.HashMap
	scopes          []string
	scopeSeparator  string
	httpClient      *request.HttpRequest
	guzzleOptions   *object.HashMap
	encodingType    int
	expiresInKey    string
	accessTokenKey  string
	refreshTokenKey string

	GetAuthURL      func() (string, error)
	GetTokenURL     func() string
	GetUserByToken  func(token string) (*object.HashMap, error)
	MapUserToObject func(userData interface{}) *src.User
}

func NewBase(config *object.HashMap) *Base {

	base := &Base{
		config: src.NewConfig(config),
		scopes: []string{},
	}

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

	return base
}

func (base *Base) Redirect(redirectURL string) (string, error) {
	if redirectURL != "" {
		base.WithRedirectURL(redirectURL)
	}

	return base.GetAuthURL()
}

func (base *Base) UserFromCode(code string) (*src.User, error) {
	tokenResponse, err := base.tokenFromCode(code)
	if err != nil {
		return nil, err
	}

	user, err := base.UserFromToken((*tokenResponse)[base.accessTokenKey].(string))
	if err != nil {
		return nil, err
	}

	refreshTokenKey := ""
	if (*tokenResponse)[base.refreshTokenKey] != nil {
		refreshTokenKey = (*tokenResponse)[base.refreshTokenKey].(string)
	}

	expiresInKey := 0
	if (*tokenResponse)[base.expiresInKey] != nil {
		expiresInKey = (*tokenResponse)[base.expiresInKey].(int)
	}

	return user.SetRefreshToken(refreshTokenKey).
		SetExpiresIn(expiresInKey).
		SetTokenResponse(tokenResponse), nil
}

func (base *Base) UserFromToken(token string) (*src.User, error) {
	user, err := base.GetUserByToken(token)
	if err != nil {
		return nil, err
	}

	return base.MapUserToObject(user).
		SetProvider(base).
		SetRaw(*user).
		SetAccessToken(token), nil
}

func (base *Base) tokenFromCode(code string) (*object.HashMap, error) {

	outResponse := &weCom.ResponseTokenFromCode{}

	response, err := base.GetHttpClient().PerformRequest(
		base.GetTokenURL(),
		"POST",
		&object.HashMap{
			"form_params": base.GetTokenFields(code),
			"headers": &object.StringMap{
				"Accept": "application/json",
			},
		},
		false, nil,
		outResponse,
	)

	return base.normalizeAccessTokenResponse(response), err
}

func (base *Base) refreshToken(refreshToken string) error {
	return errors.New("refreshToken does not support")

}

func (base *Base) WithRedirectURL(redirectURL string) src.ProviderInterface {
	base.redirectURL = redirectURL

	return base
}

func (base *Base) WithState(state string) src.ProviderInterface {
	base.state = state

	return base
}

func (base *Base) Scopes(scopes []string) *Base {
	base.scopes = scopes

	return base
}

func (base *Base) With(parameters *object.HashMap) *Base {
	base.parameters = parameters

	return base
}

func (base *Base) GetConfig() *src.Config {
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

func (base *Base) GetHttpClient() *request.HttpRequest {
	if base.httpClient != nil {
		return base.httpClient
	} else {
		return request.NewHttpRequest(base.config.All())
	}

}

func (base *Base) formatScopes(scopes []string, scopeSeparator string) string {
	return strings.Join(scopes, scopeSeparator)
}

func (base *Base) GetTokenFields(code string) *object.HashMap {
	return &object.HashMap{
		"client_id":     base.GetClientID(),
		"client_secret": base.GetClientSecret(),
		"code":          code,
		"redirect_uri":  base.redirectURL,
	}
}

func (base *Base) buildAuthURLFromBase(url string) string {
	// tbd
	return ""
}

func (base *Base) GetCodeFields() *object.HashMap {
	fields := &object.HashMap{
		"client_id":     base.GetClientID(),
		"redirect_uri":  base.redirectURL,
		"scope":         base.formatScopes(base.scopes, base.scopeSeparator),
		"response_type": "code",
	}
	fields = object.MergeHashMap(fields, base.parameters)
	if base.state != "" {
		(*fields)["state"] = base.state
	}

	return fields
}

func (base *Base) normalizeAccessTokenResponse(response contract2.ResponseInterface) *object.HashMap {
	// tbd

	return nil
}
