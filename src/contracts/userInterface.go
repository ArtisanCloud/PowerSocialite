package contracts

type UserInterface interface {

	GetID() string

	GetMobile() string

	GetNickname() string

	GetName() string

	GetEmail() string

	GetAvatar() string

	GetAccessToken() string

	GetRefreshToken() string

	GetExpiresIn() int

}
