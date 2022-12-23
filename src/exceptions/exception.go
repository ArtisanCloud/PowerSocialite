package exceptions

import "github.com/ArtisanCloud/PowerLibs/v3/exception"

type Exception struct {
	*exception.Exception
}

func NewException() *Exception {
	return &Exception{
		&exception.Exception{},
	}
}
