package providers

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/http/contract"
	"github.com/ArtisanCloud/PowerLibs/object"
	"github.com/ArtisanCloud/PowerSocialite/src/contracts"
	"reflect"
	"time"
)

type WeChat struct {
	*Base

	baseURL         string
	scopes          []string
	withCountryCode bool
	component       *object.HashMap
	openID          string
}

func NewWeChat(config *object.HashMap) *WeChat {
	wechat := &WeChat{
		Base: NewBase(config),

		baseURL:         "https://api.weixin.qq.com/sns",
		scopes:          []string{"snsapi_login"},
		withCountryCode: false,
	}

	wechat.OverrideGetAccessToken()
	wechat.OverrideGetAuthURL()
	wechat.OverrideBuildAuthURLFromBase()
	wechat.OverrideGetCodeFields()
	wechat.OverrideGetTokenURL()
	wechat.OverrideGetUserByToken()
	wechat.OverrideMapUserToObject()
	wechat.OverrideGetTokenFields()

	return wechat
}

func (provider *WeChat) WithOpenID(openid string) *WeChat {

	provider.openID = openid
	return provider
}

func (provider *WeChat) WithCountryCode() *WeChat {
	provider.withCountryCode = true

	return provider
}

func (provider *WeChat) TokenFromCode(code string) (*object.HashMap, error) {
	response, err := provider.GetTokenFromCode(code)
	if err != nil {
		return nil, err
	}

	return provider.normalizeAccessTokenResponse(response)
}

func (provider *WeChat) GetTokenFromCode(code string) (contract.ResponseInterface, error) {
	return provider.GetHttpClient().PerformRequest(provider.GetTokenURL(), "GET", &object.HashMap{
		"headers": object.StringMap{
			"Accept": "application/json",
		},
		"query": provider.getTokenFields(code),
	}, false, nil, nil)
}

func (provider *WeChat) WithComponent(component *object.HashMap) *WeChat {

	provider.PrepareForComponent(component)

	return provider
}

func (provider *WeChat) PrepareForComponent(component *object.HashMap) error {
	config := object.HashMap{}
	for k, v := range *component {
		value := v
		if reflect.TypeOf(v).Kind() == reflect.Func {
			value = reflect.ValueOf(v)
		}
		switch k {
		case "id":
		case "app_id":
		case "component_app_id":
			config["id"] = value
			break
		case "token":
		case "app_token":
		case "access_token":
		case "component_access_token":
			config["token"] = value
			break
		}
	}

	if len(config) == 2 {
		return errors.New("Please check your config arguments is available.")
	}

	if len(provider.scopes) == 1 && object.InArray("snsapi_login", provider.scopes) {
		provider.scopes = []string{"snsapi_base"}
	}

	provider.component = &config

	return nil
}

func (provider *WeChat) OverrideGetAccessToken() {
	provider.GetAccessToken = func(code string) (contracts.AccessTokenInterface, error) {
		response, err := provider.GetHttpClient().PerformRequest(provider.GetTokenURL(), "GET", &object.HashMap{
			"headers": object.StringMap{"Accept": "application/json"},
			"query":   provider.GetTokenFields(code),
		}, false, nil, nil)

		if err != nil {
			return nil, err
		}
		return provider.parseAccessToken(response.GetBody())
	}
}

func (provider *WeChat) OverrideGetAuthURL() {
	provider.GetAuthURL = func() (string, error) {

		path := "oauth2/authorize"

		// 网页授权登录
		if len(provider.scopes) > 0 {
			path = "qrconnect"
		}

		// 第三方网页应用登录（扫码登录）
		return provider.BuildAuthURLFromBase(fmt.Sprintf("https://open.weixin.qq.com/connect/%s", path)), nil
	}
}
func (provider *WeChat) OverrideBuildAuthURLFromBase() {

	provider.BuildAuthURLFromBase = func(url string) string {
		query := object.GetJoinedWithKSort(provider.GetCodeFields())

		return url + "?" + query + "#wechat_redirect"
	}
}

func (provider *WeChat) OverrideGetCodeFields() {

	provider.GetCodeFields = func() *object.StringMap {

		if provider.component != nil {
			provider.With(object.MergeStringMap(provider.parameters, &object.StringMap{
				"component_appid": (*provider.component)["id"].(string),
			}))
		}

		if provider.state == "" {
			data, _ := json.Marshal(time.Now())
			buffer := md5.Sum(data)
			provider.state = string(buffer[:])
		}

		config := provider.GetConfig()
		fields := &object.StringMap{
			"appid":            config.GetString("client_id", ""),
			"redirect_uri":     provider.redirectURL,
			"response_type":    "code",
			"scope":            provider.formatScopes(provider.scopes, provider.scopeSeparator),
			"state":            provider.state,
			"connect_redirect": "1",
		}
		fields = object.MergeStringMap(fields, provider.parameters)

		return fields
	}
}

func (provider *WeChat) OverrideGetTokenURL() {
	provider.GetTokenURL = func() string {
		if provider.component != nil {
			return provider.baseURL + "/oauth2/component/access_token"
		}
		return provider.baseURL + "/oauth2/access_token"
	}
}

func (provider *WeChat) OverrideGetUserByToken() {
	provider.GetUserByToken = func(token string) (*object.HashMap, error) {

		return nil, errors.New("WeCom doesn't support access_token mode")
	}
}

func (provider *WeChat) OverrideMapUserToObject() {

	provider.MapUserToObject = func(user *object.HashMap) *User {

		collectionUser := object.NewCollection(user)

		// weCom.ResponseGetUserInfo is response from code to user
		return NewUser(&object.HashMap{
			"id":       collectionUser.Get("openid", ""),
			"name":     collectionUser.Get("nickname", ""),
			"nickname": collectionUser.Get("nickname", ""),
			"avatar":   collectionUser.Get("headimgurl", ""),
			"email":    nil,
		}, provider)
	}
}

func (provider *WeChat) OverrideGetTokenFields() {
	provider.GetTokenFields = func(code string) *object.HashMap {

		config := provider.GetConfig()
		return &object.HashMap{
			"appid":                  config.GetString("client_id", ""),
			"secret":                 config.GetString("client_secret", ""),
			"component_appid":        provider.GetClientID(),
			"component_access_token": (*provider.component)["token"],
			"code":                   code,
			"grant_type":             "authorization_code",
		}
	}
}
