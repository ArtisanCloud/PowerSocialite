package providers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/object"
	"github.com/ArtisanCloud/PowerSocialite/src/contracts"
	"time"
)

type WeChat struct {
	*Base

	baseURL         string
	openID          bool
	scopes          []string
	stateless       bool
	withCountryCode bool
	component       contracts.WechatComponentInterface
}

func NewWeChat(config *object.HashMap) *WeChat {
	wechat := &WeChat{
		Base: NewBase(config),

		baseURL:         "https://api.weixin.qq.com/sns",
		scopes:          []string{"snsapi_login"},
		stateless:       true,
		withCountryCode: false,
	}

	wechat.OverrideGetAccessToken()
	wechat.OverrideGetAuthURL()
	wechat.OverrideBuildAuthURLFromBase()
	wechat.OverrideGetCodeFields()
	wechat.OverrideGetTokenURL()

	return wechat
}

func (provider *WeChat) WithCountryCode() *WeChat {
	provider.withCountryCode = true

	return provider
}

func (provider *WeChat) Component(component contracts.WechatComponentInterface) *WeChat {
	provider.scopes = []string{"snsapi_base"}
	provider.component = component

	return provider
}

func (provider *WeChat) OverrideGetAccessToken() {
	provider.GetAccessToken = func(code string) (contracts.AccessTokenInterface, error) {
		response, err := provider.GetHttpClient().PerformRequest(provider.GetTokenURL(""), "GET", &object.HashMap{
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
	provider.GetAuthURL = func(state string) (string, error) {

		path := "oauth2/authorize"

		// 网页授权登录
		if len(provider.scopes) > 0 {
			path = "qrconnect"
		}

		// 第三方网页应用登录（扫码登录）
		return provider.BuildAuthURLFromBase(fmt.Sprintf("https://open.weixin.qq.com/connect/%s", path), state), nil
	}
}
func (provider *WeChat) OverrideBuildAuthURLFromBase() {

	provider.BuildAuthURLFromBase = func(url string, state string) string {
		query := object.GetJoinedWithKSort(provider.GetCodeFields(state))

		return url + "?" + query + "#wechat_redirect"
	}
}

func (provider *WeChat) OverrideGetCodeFields() {

	provider.GetCodeFields = func(state string) *object.StringMap {

		if provider.component != nil {
			provider.With(object.MergeStringMap(provider.parameters, &object.StringMap{
				"component_appid": provider.component.GetAppID(),
			}))
		}

		if state == "" {
			data, _ := json.Marshal(time.Now())
			buffer := md5.Sum(data)
			state = string(buffer[:])
		}

		config := provider.GetConfig()
		fields := &object.StringMap{
			"appid":            config.GetString("client_id", ""),
			"redirect_uri":     provider.redirectURL,
			"response_type":    "code",
			"scope":            provider.formatScopes(provider.scopes, provider.scopeSeparator),
			"state":            state,
			"connect_redirect": "1",
		}
		fields = object.MergeStringMap(fields, provider.parameters)

		return fields
	}
}

func (provider *WeChat) OverrideGetTokenURL() {
	provider.GetTokenURL = func(state string) string {
		if provider.component != nil {
			return provider.baseURL + "/oauth2/component/access_token"
		}
		return provider.baseURL + "/oauth2/access_token"
	}
}
