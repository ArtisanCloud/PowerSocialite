package src

import "github.com/ArtisanCloud/go-libs/object"

type Config struct {
	*object.Collection

}

func NewConfig(config *object.HashMap) *Config {
	return &Config{
		object.NewCollection(config),
	}
}
