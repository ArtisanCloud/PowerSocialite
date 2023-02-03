package main

import (
	"fmt"
	fmt2 "github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerSocialite/v3/src/providers"
)

func main() {

	fmt.Printf("hello Socialite! \n")
	provider := providers.NewWeCom(&object.HashMap{})
	fmt2.Dump(provider.GetConfig())

}
