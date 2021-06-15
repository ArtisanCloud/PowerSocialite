package weCom

import "github.com/ArtisanCloud/go-socialite/src/models"

type ResponseGetUserByID struct {
	*ResponseWeCom
	*models.Employee

}
