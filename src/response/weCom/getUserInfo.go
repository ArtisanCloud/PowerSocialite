package weCom

type ResponseGetUserInfo struct {
	UserID   string `json:"UserId"`
	DeviceID string `json:"DeviceId"`
	OpenID   string `json:"OpenId"`
	*ResponseWeCom
}
