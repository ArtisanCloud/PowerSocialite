package google

import "github.com/ArtisanCloud/PowerSocialite/v4/config"

type Config struct {
	config.BaseConfig

	ResponseType          string   `json:"responseType" yaml:"response_type"`
	State                 string   `json:"state,omitempty" yaml:"state,omitempty"`
	AccessType            string   `json:"accessType,omitempty" yaml:"access_type,omitempty"`
	IncludeGrantedScopes  []string `json:"includeGrantedScopes,omitempty" yaml:"include_granted_scopes,omitempty"`
	EnableGranularConsent bool     `json:"enableGranularConsent,omitempty" yaml:"enable_granular_consent,omitempty"`
	LoginHint             string   `json:"loginHint,omitempty" yaml:"login_hint,omitempty"`
	Prompt                string   `json:"prompt,omitempty" yaml:"prompt,omitempty"`
	HostedDomain          string   `json:"hostedDomain,omitempty" yaml:"hosted_domain,omitempty"`
	SelectAccount         bool     `json:"selectAccount,omitempty" yaml:"select_account,omitempty"`
}
