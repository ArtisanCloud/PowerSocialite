package providers

import (
	"github.com/ArtisanCloud/go-libs/http"
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/ArtisanCloud/go-socialite/src"
	"github.com/ArtisanCloud/go-socialite/src/contracts"
)

type Base struct {
	contracts.ProviderInterface

	state           string
	config          *src.Config
	redirectUrl     string
	parameters      *object.HashMap
	scopes          *object.HashMap
	scopeSeparator  string
	httpClient      http.HttpRequest
	guzzleOptions   *object.HashMap
	encodingType    int
	expiresInKey    string
	accessTokenKey  string
	refreshTokenKey string
}


