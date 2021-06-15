package exceptions

import "github.com/ArtisanCloud/go-libs/exception"

type Exception struct {
	*exception.Exception
}

func NewException() *Exception {
	return &Exception{
		&exception.Exception{},
	}
}