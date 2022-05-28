package providers

type ProviderInterface interface {
	Redirect(redirectURL string) (string, error)

	UserFromCode(code string) (*User, error)

	// 多协程运作时，openID需要作为独立参数传入
	UserFromToken(token string, openID string) (*User, error)
}
