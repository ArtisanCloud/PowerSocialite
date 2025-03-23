package contracts

import "github.com/ArtisanCloud/PowerSocialite/v4/models"
import "golang.org/x/oauth2"

type IProvider interface {
	Name() string
	SetName(name string)
	StartAuth(state string) (ISession, error)
	GetSession(string) (ISession, error)
	UserFromSession(ISession) (*models.User, error)
	Debug(bool)
	RefreshToken(refreshToken string) (*oauth2.Token, error)
	RefreshTokenAvailable() bool
}
