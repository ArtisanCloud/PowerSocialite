package config

type BaseConfig struct {
	Name         string   `json:"name" yaml:"name"`
	Scopes       []string `json:"scopes" yaml:"scopes"`
	AuthUrl      string   `json:"authUrl" yaml:"auth_url"`
	TokenUrl     string   `json:"tokenUrl,omitempty" yaml:"token_url,omitempty"`
	ClientId     string   `json:"clientId,omitempty" yaml:"client_id,omitempty"`
	ClientSecret string   `json:"clientSecret,omitempty" yaml:"client_secret,omitempty"`
	RedirectUrl  string   `json:"redirectUrl,omitempty" yaml:"redirect_url,omitempty"`
	CallbackUrl  string   `json:"callbackUrl,omitempty" yaml:"callback_url,omitempty"`
}
