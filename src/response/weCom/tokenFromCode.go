package weCom

type ResponseTokenFromCode struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	*ResponseWeCom
}
