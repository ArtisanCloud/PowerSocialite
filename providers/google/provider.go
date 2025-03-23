package google

import (
	"encoding/json"
	"fmt"
	"github.com/ArtisanCloud/PowerSocialite/v4/contracts"
	"github.com/ArtisanCloud/PowerSocialite/v4/kernel"
	"github.com/ArtisanCloud/PowerSocialite/v4/models"
	"github.com/ArtisanCloud/PowerSocialite/v4/providers"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const endpointAuth string = "https://accounts.google.com/o/oauth2/auth"
const endpointProfile string = "https://www.googleapis.com/oauth2/v2/userinfo"

type Provider struct {
	*providers.BaseProvider
	config          *Config
	authCodeOptions []oauth2.AuthCodeOption
}

func NewProvider(config *Config) *Provider {
	baseProvider := providers.NewBaseProvider(&config.BaseConfig)
	provider := &Provider{
		BaseProvider: baseProvider,
		config:       config,
		authCodeOptions: []oauth2.AuthCodeOption{
			oauth2.AccessTypeOffline,
		},
	}

	if config.BaseConfig.Scope == nil || len(config.BaseConfig.Scope) == 0 {
		config.BaseConfig.Scope = []string{"email"}
	}

	return provider
}

// StartAuth asks Google for an authentication endpoint.
func (p *Provider) StartAuth(state string) (contracts.ISession, error) {
	url := p.OAuthConfig.AuthCodeURL(state, p.authCodeOptions...)
	session := &Session{
		BaseSession: providers.BaseSession{
			AuthUrl: url,
		},
	}
	return session, nil
}

func (p *Provider) UserFromSession(session contracts.ISession) (*models.User, error) {
	sess := session.(*Session)
	user := models.User{
		AccessToken:  sess.AccessToken,
		Provider:     p.Name(),
		RefreshToken: sess.RefreshToken,
		ExpiresAt:    sess.ExpiresAt,
		IDToken:      sess.IDToken,
	}

	if user.AccessToken == "" {
		// Data is not yet retrieved, since accessToken is still empty.
		return nil, fmt.Errorf("%s cannot get user information without accessToken", p.BaseProvider.Name())
	}

	response, err := p.Client().Get(endpointProfile + "?access_token=" + url.QueryEscape(sess.AccessToken))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with a %d trying to fetch user information", p.BaseProvider.Name(), response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var u User
	if err := json.Unmarshal(responseBytes, &u); err != nil {
		return nil, err
	}

	// Extract the user data we got from Google into our goth.User.
	user.Name = u.Name
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.NickName = u.Name
	user.Email = u.Email
	user.AvatarURL = u.Picture
	user.UserID = u.ID
	// Google provides other useful fields such as 'hd'; get them from RawData
	if err := json.Unmarshal(responseBytes, &user.Data); err != nil {
		return nil, err
	}

	return &user, nil
}

// RefreshToken get new access token based on the refresh token
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := p.OAuthConfig.TokenSource(kernel.ContextForClient(p.Client()), token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newToken, err
}

// SetPrompt sets the prompt values for the google OAuth call. Use this to
// force users to choose and account every time by passing "select_account",
// for example.
// See https://developers.google.com/identity/protocols/OpenIDConnect#authenticationuriparameters
func (p *Provider) SetPrompt(prompt ...string) {
	if len(prompt) == 0 {
		return
	}
	p.authCodeOptions = append(p.authCodeOptions, oauth2.SetAuthURLParam("prompt", strings.Join(prompt, " ")))
}

// SetHostedDomain sets the hd parameter for google OAuth call.
// Use this to force user to pick user from specific hosted domain.
// See https://developers.google.com/identity/protocols/oauth2/openid-connect#hd-param
func (p *Provider) SetHostedDomain(hd string) {
	if hd == "" {
		return
	}
	p.authCodeOptions = append(p.authCodeOptions, oauth2.SetAuthURLParam("hd", hd))
}

// SetLoginHint sets the login_hint parameter for the Google OAuth call.
// Use this to prompt the user to log in with a specific account.
// See https://developers.google.com/identity/protocols/oauth2/openid-connect#login-hint
func (p *Provider) SetLoginHint(loginHint string) {
	if loginHint == "" {
		return
	}
	p.authCodeOptions = append(p.authCodeOptions, oauth2.SetAuthURLParam("login_hint", loginHint))
}

// SetAccessType sets the access_type parameter for the Google OAuth call.
// If an access token is being requested, the client does not receive a refresh token unless a value of offline is specified.
// See https://developers.google.com/identity/protocols/oauth2/openid-connect#access-type-param
func (p *Provider) SetAccessType(at string) {
	if at == "" {
		return
	}
	p.authCodeOptions = append(p.authCodeOptions, oauth2.SetAuthURLParam("access_type", at))
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (contracts.ISession, error) {
	sess := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(sess)
	return sess, err
}
