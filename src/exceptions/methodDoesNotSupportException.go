package exceptions

type MethodDoesNotSupportException struct {
	*Exception
}

func NewMethodDoesNotSupportException() *MethodDoesNotSupportException {
	return &MethodDoesNotSupportException{
		NewException(),
	}
}