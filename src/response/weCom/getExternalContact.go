package weCom

import "github.com/ArtisanCloud/go-socialite/src/models"

type ResponseGetExternalContact struct {
	*ResponseWeCom
	*models.ExternalContact `json:"external_contact"`
}


