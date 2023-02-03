package src

import (
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/ArtisanCloud/PowerSocialite/v2/src/models"
	"github.com/ArtisanCloud/PowerSocialite/v2/src/providers"
	"strings"
)

type SocialiteManager struct {
	Config         *models.Config
	Resolved       *object.HashMap
	CustomCreators *object.HashMap

	Providers []string
}

func NewSocialiteManager(config *object.HashMap) *SocialiteManager {

	manager := &SocialiteManager{
		Config:         models.NewConfig(config),
		Resolved:       &object.HashMap{},
		CustomCreators: &object.HashMap{},
		Providers: []string{
			"wechat",
			"wecom",
			"openplatform",
		},
	}

	return manager
}

func (manager *SocialiteManager) SetConfig(config *models.Config) *SocialiteManager {
	manager.Config = config
	return manager
}

func (manager *SocialiteManager) Create(name string) providers.ProviderInterface {
	name = strings.ToLower(name)

	if (*manager.Resolved)[name] == nil {
		(*manager.Resolved)[name] = manager.CreateProvider(name)
	}

	return (*manager.Resolved)[name].(providers.ProviderInterface)
}

func (manager *SocialiteManager) Extend(name string, callback func()) *SocialiteManager {
	(*manager.CustomCreators)[strings.ToLower(name)] = callback
	return manager
}

func (manager *SocialiteManager) GetResolvedProviders() *object.HashMap {
	return manager.Resolved
}

func (manager *SocialiteManager) BuildProvider(provider string, config *object.HashMap) providers.ProviderInterface {
	switch provider {
	case "wechat":
		return providers.NewWeChat(config)

		break

	case "wecom":
		return providers.NewWeCom(config)

		break

	case "openplatform":
		return providers.NewOpenPlatformform(config)

		break

	default:

	}
	return nil
}

func (manager *SocialiteManager) CreateProvider(name string) providers.ProviderInterface {
	config := manager.Config.Get(name, &object.HashMap{}).(*object.HashMap)
	provider := name
	if config != nil && (*config)["provider"] != nil {
		provider = (*config)["provider"].(string)
	}

	// not supported yet
	if (*manager.CustomCreators)[provider] != nil {
		return manager.CallCustomCreator(provider, config)
	}

	if !manager.IsValidProvider(provider) {
		panic(fmt.Sprintf("Provider [%s] not supported.", provider))
	}

	return manager.BuildProvider(provider, config)
}

func (manager *SocialiteManager) CallCustomCreator(driver string, config *object.HashMap) providers.ProviderInterface {

	return nil
}

func (manager *SocialiteManager) IsValidProvider(provider string) bool {
	return object.ContainsString(manager.Providers, provider)

}
