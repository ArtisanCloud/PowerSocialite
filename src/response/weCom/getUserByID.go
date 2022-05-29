package weCom

import "github.com/ArtisanCloud/PowerSocialite/v2/src/models"

type ResponseGetUserByID struct {
	*ResponseWeCom
	*models.Employee
}
