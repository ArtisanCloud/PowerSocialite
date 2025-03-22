package models

import (
	"github.com/ArtisanCloud/PowerSocialite/v4/utils/object"
	"time"
)

type User struct {
	Data              object.HashMap
	Provider          string
	Email             string
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	Description       string
	UserID            string
	AvatarURL         string
	Location          string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         time.Time
	IDToken           string
}
