package main

import (
	"fmt"
	fmt2 "github.com/ArtisanCloud/go-libs/fmt"
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/ArtisanCloud/go-socialite/src/providers"
)

func main() {

	fmt.Printf("hello Socialite! \n")
	provider:=providers.NewWeCom(&object.HashMap{})
	fmt2.Dump(provider.GetConfig())

}