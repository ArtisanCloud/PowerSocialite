package weCom

import (
	"github.com/ArtisanCloud/PowerSocialite/v3/src/models"
)

type ResponseGetUserDetail struct {
	*ResponseWeCom

	UserID  string `json:"userid"`
	Gender  string `json:"gender"`
	Avatar  string `json:"avatar"`
	QrCode  string `json:"qr_code"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	BizMail string `json:"biz_mail"`
	Address string `json:"address"`
}

type ResponseGetUserInfo struct {
	*ResponseWeCom
	DeviceID       string `json:"DeviceId"`
	UserID         string `json:"UserId"`
	ExternalUserID string `json:"external_userid"`
	OpenID         string `json:"OpenId"`
	UserTicket     string `json:"user_ticket"`
}

type ResponseGetUserByID struct {
	*ResponseWeCom
	*models.Employee
}
