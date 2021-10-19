package main

import (
	"fmt"
	fmt2 "github.com/ArtisanCloud/PowerLibs/fmt"
	"github.com/ArtisanCloud/PowerLibs/object"
	"github.com/ArtisanCloud/PowerSocialite/src/providers"
)

func main() {

	fmt.Printf("hello Socialite! \n")
	provider:=providers.NewWeCom(&object.HashMap{})
	fmt2.Dump(provider.GetConfig())

}