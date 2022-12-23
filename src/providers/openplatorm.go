package providers

import (
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"net/url"
)

const (
	redirectOauthURL       = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	webAppRedirectOauthURL = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	accessTokenURL         = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	refreshAccessTokenURL  = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	userInfoURL            = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=%s"
	checkAccessTokenURL    = "https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s"
)

type OpenPlatform struct {
	*Base
}

func NewOpenPlatformform(oh *object.HashMap) *OpenPlatform {
	return &OpenPlatform{
		Base: NewBase(oh),
	}
}

func (provider *OpenPlatform) GetName() string {
	return "openplatform"
}

// GetRedirectURL 获取跳转的url地址
func (provider *OpenPlatform) GetRedirectURL(redirectURI, scope, state string) (string, error) {
	// url encode
	urlStr := url.QueryEscape(redirectURI)
	return fmt.Sprintf(redirectOauthURL, provider.config.Get("app_id", ""), urlStr, scope, state), nil
}

// GetWebAppRedirectURL 获取网页应用跳转的url地址
func (provider *OpenPlatform) GetWebAppRedirectURL(redirectURI, scope, state string) (string, error) {
	urlStr := url.QueryEscape(redirectURI)
	return fmt.Sprintf(webAppRedirectOauthURL, provider.config.Get("app_id", ""), urlStr, scope, state), nil
}

//// Redirect OpenPlatform
//func (provider *OpenPlatform) Redirect(writer http.ResponseWriter, req *http.Request, redirectURI, scope, state string) error {
//	location, err := provider.GetRedirectURL(redirectURI, scope, state)
//	if err != nil {
//		return err
//	}
//	http.Redirect(writer, req, location, http.StatusFound)
//	return nil
//}

type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// ResAccessToken 获取用户授权access_token的返回结果
type ResAccessToken struct {
	CommonError
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`

	// UnionID 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	UnionID string `json:"unionid"`
}

// GetUserAccessToken 通过网页授权的code 换取access_token(区别于context中的access_token)
func (provider *OpenPlatform) GetUserAccessToken(code string) (result ResAccessToken, err error) {
	urlStr := fmt.Sprintf(accessTokenURL, provider.config.Get("app_id", ""), provider.config.Get("app_secret", ""), code)
	client, err := provider.GetHttpClient()
	if err != nil {
		return result, err
	}
	err = client.Df().Url(urlStr).Method("GET").Result(&result)
	if err != nil {
		return result, err
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

// RefreshAccessToken 刷新access_token
func (provider *OpenPlatform) RefreshAccessToken(refreshToken string) (result ResAccessToken, err error) {
	urlStr := fmt.Sprintf(refreshAccessTokenURL, provider.config.Get("app_id", ""), refreshToken)
	client, err := provider.GetHttpClient()
	if err != nil {
		return result, err
	}
	err = client.Df().Url(urlStr).Method("GET").Result(&result)
	if err != nil {
		return result, err
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("RefreshAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

func (provider *OpenPlatform) CheckAccessToken(accessToken, openID string) (b bool, err error) {
	urlStr := fmt.Sprintf(checkAccessTokenURL, accessToken, openID)
	var result CommonError
	client, err := provider.GetHttpClient()
	if err != nil {
		return false, err
	}
	err = client.Df().Url(urlStr).Method("GET").Result(&result)
	if err != nil {
		return false, err
	}

	if result.ErrCode != 0 {
		b = false
		return
	}
	b = true
	return
}

// UserInfo 用户授权获取到用户信息
type UserInfo struct {
	CommonError

	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int32    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

// GetUserInfo 如果scope为 snsapi_userinfo 则可以通过此方法获取到用户基本信息
func (provider *OpenPlatform) GetUserInfo(accessToken, openID, lang string) (result UserInfo, err error) {
	if lang == "" {
		lang = "zh_CN"
	}
	urlStr := fmt.Sprintf(userInfoURL, accessToken, openID, lang)
	client, err := provider.GetHttpClient()
	if err != nil {
		return result, err
	}
	err = client.Df().Url(urlStr).Method("GET").Result(&result)
	if err != nil {
		return result, err
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return result, err
	}
	return result, err
}
