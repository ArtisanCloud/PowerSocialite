package weCom

import "github.com/ArtisanCloud/PowerSocialite/v2/src/models"

type ResponseGetExternalContact struct {
	*ResponseWeCom
	*models.ExternalContact `json:"external_contact"`
	FollowInfo              []*models.FollowUser `json:"follow_user"`
}
