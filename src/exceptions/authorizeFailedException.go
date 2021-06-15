package exceptions

type AuthorizeFailedException struct {
	*Exception
}

func NewAuthorizeFailedException() *AuthorizeFailedException {
	return &AuthorizeFailedException{
		NewException(),
	}
}
