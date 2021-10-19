package weCom

type ResponseGetUserInfo struct {
	*ResponseWeCom
	DeviceID       string `json:"DeviceId"`
	UserID         string `json:"UserId"`
	ExternalUserID string `json:"external_userid"`
	OpenID         string `json:"OpenId"`
}
