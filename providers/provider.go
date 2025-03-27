package providers

import (
	"bytes"
	"encoding/json"
	"github.com/ArtisanCloud/PowerSocialite/v4/auth"
	"github.com/ArtisanCloud/PowerSocialite/v4/config"
	"github.com/ArtisanCloud/PowerSocialite/v4/contracts"
	"github.com/ArtisanCloud/PowerSocialite/v4/models"
	"github.com/ArtisanCloud/PowerSocialite/v4/utils/object"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strings"
)

type BaseProvider struct {
	contracts.IProvider

	OAuthConfig *oauth2.Config

	providerName    string
	state           string
	forcePopup      bool
	parameters      *object.StringMap
	scopeSeparator  string
	HTTPClient      *http.Client
	guzzleOptions   *object.HashMap
	encodingType    int
	expiresInKey    string
	accessTokenKey  string
	refreshTokenKey string

	TokenFromCode func(code string) (contracts.ISession, error)
	GetAuthURL    func() (string, error)
	GetTokenURL   func() string

	GetAccessToken       func(token string) (contracts.ISession, error)
	BuildAuthURLFromBase func(url string) string
}

func NewBaseProvider(config *config.BaseConfig) (*BaseProvider, error) {

	base := &BaseProvider{
		expiresInKey:    "expires_in",
		accessTokenKey:  "access_token",
		refreshTokenKey: "refresh_token",
		OAuthConfig: &oauth2.Config{
			ClientID:     config.ClientId,
			ClientSecret: config.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  config.AuthUrl,
				TokenURL: config.TokenUrl,
			},
			RedirectURL: config.RedirectUrl,
			Scopes:      config.Scopes,
		},
	}

	return base, nil
}

// Name is the name used to retrieve this provider later.
func (p *BaseProvider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *BaseProvider) SetName(name string) {
	p.providerName = name
}
func (p *BaseProvider) StartAuth(state string) (contracts.ISession, error) {
	return nil, nil
}

func (p *BaseProvider) GetSession(string) (contracts.ISession, error) {
	return nil, nil
}

func (p *BaseProvider) UserFromSession(session contracts.ISession) (*models.User, error) {
	return nil, nil
}

func (p *BaseProvider) Debug(debug bool) {}

func (p *BaseProvider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return nil, nil
}
func (p *BaseProvider) RefreshTokenAvailable() bool {
	return true
}
func (p *BaseProvider) Redirect(redirectURL string) (string, error) {
	if redirectURL != "" {
		p.SetRedirectURL(redirectURL)
	}

	return p.GetAuthURL()
}

func (p *BaseProvider) UserFromCode(code string) (*models.User, error) {
	_, err := p.TokenFromCode(code)
	if err != nil {
		return nil, err
	}
	// 通过获取的code，获取用户信息

	user := models.User{}
	//err = object.HashMapToStructure(tokenResponse, &user)

	return &user, err
}

func (p *BaseProvider) SetRedirectURL(redirectURL string) *BaseProvider {
	p.OAuthConfig.RedirectURL = redirectURL

	return p
}

func (p *BaseProvider) SetState(state string) *BaseProvider {
	p.state = state

	return p
}

func (p *BaseProvider) SetForcePopup(forcePopup bool) *BaseProvider {
	p.forcePopup = forcePopup

	return p
}

func (p *BaseProvider) SetScopes(scopes []string) *BaseProvider {
	p.OAuthConfig.Scopes = scopes

	return p
}

func (p *BaseProvider) SetParameters(parameters *object.StringMap) *BaseProvider {
	p.parameters = parameters

	return p
}

func (p *BaseProvider) SetScopeSeparator(scopeSeparator string) *BaseProvider {
	p.scopeSeparator = scopeSeparator

	return p
}

func (p *BaseProvider) Client() *http.Client {
	return auth.HTTPClientWithFallBack(p.HTTPClient)
}

func (p *BaseProvider) SetFormatScopes(scopes []string, scopeSeparator string) string {
	return strings.Join(scopes, scopeSeparator)
}

func (p *BaseProvider) GetTokenFields(code string) *object.StringMap {
	return &object.StringMap{
		"client_id":     p.OAuthConfig.ClientID,
		"client_secret": p.OAuthConfig.ClientSecret,
		"code":          code,
		"redirect_uri":  p.OAuthConfig.RedirectURL,
	}
}

func (p *BaseProvider) ParseBody(body io.ReadCloser) (*object.HashMap, error) {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(body)
	jsonHashMap := object.HashMap{}
	err := json.Unmarshal(buf.Bytes(), &jsonHashMap)

	return &jsonHashMap, err
}
