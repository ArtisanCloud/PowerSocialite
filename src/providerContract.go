package src

type ProviderInterface interface {
	Redirect(redirectURL string) (string, error)
	UserFromCode(code string, isExternal bool) (*User, error)
	UserFromToken(token string) (*User, error)
}
