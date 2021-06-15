package src

type ProviderInterface interface {
	Redirect(redirectURL string) string
	UserFromCode(code string, isExternal bool) *User
	UserFromToken(token string) *User
}
