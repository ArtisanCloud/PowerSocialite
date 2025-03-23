package config

type BaseConfig struct {
	Name         string   `json:"name" yaml:"name"`
	Scope        []string `json:"scope" yaml:"scope"`
	AuthUrl      string   `json:"authUrl" yaml:"auth_url"`
	ClientId     string   `json:"clientId,omitempty" yaml:"client_id,omitempty"`
	ClientKey    string   `json:"clientKey,omitempty" yaml:"client_key,omitempty"`
	ClientSecret string   `json:"clientSecret,omitempty" yaml:"client_secret,omitempty"`
	RedirectUri  string   `json:"redirectUri,omitempty" yaml:"redirect_uri,omitempty"`
	CallbackUrl  string   `json:"callbackUrl,omitempty" yaml:"callback_url,omitempty"`
}
