package configs

import "github.com/ArtisanCloud/PowerLibs/v3/object"

type Config struct {
	*object.Collection
}

func NewConfig(config *object.HashMap) *Config {
	return &Config{
		object.NewCollection(config),
	}
}
