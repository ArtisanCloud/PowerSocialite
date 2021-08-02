package weCom

type ResponseGetUserInfo struct {
	UserID         string `json:"UserId"`
	DeviceID       string `json:"DeviceId"`
	OpenID         string `json:"OpenId"`
	ExternalUserID string `json:"external_userid"`
	*ResponseWeCom
}
