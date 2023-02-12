package providers

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/http/helper"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerSocialite/v3/src/response/wechat"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type WeChat struct {
	*Base

	baseURL         string
	scopes          []string
	withCountryCode bool
	component       *object.HashMap
}

func NewWeChat(config *object.HashMap) *WeChat {
	wechat := &WeChat{
		Base: NewBase(config),

		baseURL:         "https://api.weixin.qq.com/sns",
		scopes:          []string{"snsapi_login"},
		withCountryCode: false,
	}

	wechatConfig := wechat.GetConfig()
	if wechatConfig.Has("component") {
		wechat.PrepareForComponent(wechatConfig.Get("component", &object.HashMap{}).(*object.HashMap))
	}

	wechat.OverrideTokenFromCode()
	wechat.OverrideGetAuthURL()
	wechat.OverrideBuildAuthURLFromBase()
	wechat.OverrideGetCodeFields()
	wechat.OverrideGetTokenURL()
	wechat.OverrideGetUserByToken()
	wechat.OverrideMapUserToObject()
	wechat.OverrideGetTokenFields()

	return wechat
}

func (provider *WeChat) GetHttpClient() (*helper.RequestHelper, error) {
	if provider.httpHelper != nil {
		return provider.httpHelper, nil
	} else {
		h, err := helper.NewRequestHelper(&helper.Config{
			BaseUrl: provider.baseURL,
		})

		h.WithMiddleware(helper.HttpDebugMiddleware(provider.GetConfig().GetBool("http_debug", false)))

		return h, err
	}

}

func (provider *WeChat) GetName() string {
	return "wechat"
}
func (provider *WeChat) SetScopes(scopes []string) {
	provider.scopes = scopes
}

func (provider *WeChat) WithCountryCode() *WeChat {
	provider.withCountryCode = true

	return provider
}

func (provider *WeChat) OverrideTokenFromCode() {
	provider.TokenFromCode = func(code string) (*object.HashMap, error) {
		response, err := provider.GetTokenFromCode(code)
		if err != nil {
			return nil, err
		}
		return provider.normalizeAccessTokenResponse(response)
	}
}

func (provider *WeChat) WithComponent(component *object.HashMap) *WeChat {

	provider.PrepareForComponent(component)

	return provider
}

func (provider *WeChat) GetComponent() *object.HashMap {

	return provider.component
}

func (provider *WeChat) OverrideGetAuthURL() {
	provider.GetAuthURL = func() (string, error) {

		path := "oauth2/authorize"

		// 网页授权登录
		if len(provider.scopes) > 0 && object.ContainsString(provider.scopes, "snsapi_login") {
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

func (provider *WeChat) UserFromCode(code string) (*User, error) {
	if object.InArray("snsapi_base", provider.scopes) {
		tokenResponse, err := provider.GetTokenFromCode(code)
		if err != nil {
			return nil, err
		}
		bodyBuffer, err := ioutil.ReadAll(tokenResponse.Body)
		if err != nil {
			return nil, err
		}
		mapToken := &object.HashMap{}
		err = object.JsonDecode(bodyBuffer, mapToken)

		user := provider.MapUserToObject(mapToken)
		if user.GetString("id", "") == "" {
			return nil, errors.New((*mapToken)["errmsg"].(string))
		}

		return user, nil
	}

	tokenResponse, err := provider.TokenFromCode(code)
	if err != nil {
		return nil, err
	}

	// 检查is_snapshotuser是否返回
	isSnapShotUser := 0.0
	if (*tokenResponse)["is_snapshotuser"] != nil {
		isSnapShotUser = (*tokenResponse)["is_snapshotuser"].(float64)
	}

	token := (*tokenResponse)[provider.accessTokenKey].(string)
	openID := (*tokenResponse)["openid"].(string)
	user, err := provider.UserFromToken(token, openID)
	if err != nil {
		return nil, err
	}

	refreshToken := ""
	if (*tokenResponse)[provider.refreshTokenKey] != nil {
		refreshToken = (*tokenResponse)[provider.refreshTokenKey].(string)
	}

	expiresIn := 0.0
	if (*tokenResponse)[provider.expiresInKey] != nil {
		expiresIn = (*tokenResponse)[provider.expiresInKey].(float64)
	}

	return user.SetSnapShotUser(isSnapShotUser == 1).
		SetRefreshToken(refreshToken).
		SetExpiresIn(expiresIn).
		SetTokenResponse(tokenResponse), nil
}

//
func (provider *WeChat) OverrideGetUserByToken() {
	provider.GetUserByToken = func(token string, openID string) (*object.HashMap, error) {

		language := ""
		if provider.withCountryCode {
			if (*provider.parameters)["lang"] != "" {
				language = (*provider.parameters)["lang"]
			} else {
				language = "zh_CN"
			}
		}

		client, err := provider.GetHttpClient()
		if err != nil {
			return nil, err
		}
		response, err := client.Df().Url(provider.baseURL+"/userinfo").
			Method("GET").
			Query("access_token", token).
			Query("openid", openID).
			Query("lang", language).
			Request()

		if err != nil {
			return nil, err
		}

		return provider.ParseBody(response.Body)

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
			"openID":   collectionUser.Get("openid", ""),
			"email":    nil,
		}, provider)
	}
}

func (provider *WeChat) OverrideGetTokenFields() {
	provider.GetTokenFields = func(code string) *object.StringMap {

		if provider.component != nil {
			return &object.StringMap{
				"appid":                  provider.GetClientID(),
				"component_appid":        (*provider.component)["id"].(string),
				"component_access_token": (*provider.component)["token"].(string),
				"code":                   code,
				"grant_type":             "authorization_code",
			}
		}

		config := provider.GetConfig()
		return &object.StringMap{
			"appid":      config.GetString("client_id", ""),
			"secret":     config.GetString("client_secret", ""),
			"code":       code,
			"grant_type": "authorization_code",
		}
	}
}

func (provider *WeChat) GetTokenFromCode(code string) (*http.Response, error) {
	result := &wechat.ResponseAuthenticatedAccessToken{}
	client, err := provider.GetHttpClient()
	if err != nil {
		return nil, err
	}
	df := client.Df().Url(provider.GetTokenURL()).Method("GET").
		Header("Accept", "application/json")

	queries := provider.GetTokenFields(code)
	for k, v := range *queries {
		df.Query(k, v)
	}

	rs, err := df.Request()
	if err != nil {
		return nil, err
	}
	err = client.ParseResponseBodyContent(rs, result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, errors.New(result.ErrMsg)
	}

	return rs, err
}

func (provider *WeChat) PrepareForComponent(component *object.HashMap) error {
	config := object.HashMap{}
	for k, v := range *component {
		value := v
		if reflect.TypeOf(v).Kind() == reflect.Func {
			value = reflect.ValueOf(v)
		}
		switch k {
		case "id", "app_id", "component_app_id":
			config["id"] = value
			break
		case "token", "app_token", "access_token", "component_access_token":
			config["token"] = value
			break
		}
	}

	if len(config) != 2 {
		return errors.New("Please check your config arguments is available.")
	}

	if len(provider.scopes) == 1 && object.InArray("snsapi_login", provider.scopes) {
		provider.scopes = []string{"snsapi_base"}
	}

	provider.component = &config

	return nil
}
