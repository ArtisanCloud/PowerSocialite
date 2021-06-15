package src

import "github.com/ArtisanCloud/go-libs/object"

type SocialiteManager struct {
	Config         *Config
	Resolved       *object.HashMap
	CustomCreators *object.HashMap
}
