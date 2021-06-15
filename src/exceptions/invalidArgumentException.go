package exceptions

type InvalidArgumentException struct {
	*Exception
}

func NewInvalidArgumentException() *InvalidArgumentException {
	return &InvalidArgumentException{
		NewException(),
	}
}