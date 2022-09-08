package wechat

type WechatBase struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type ResponseAuthenticatedAccessToken struct {
	WechatBase

	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid,omitempty"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid,omitempty"`
	// 是否为快照页模式虚拟账号，只有当用户是快照页模式虚拟账号是返回，值为1
	IsSnapShotUser int `json:"is_snapshotuser,omitempty"`
}
