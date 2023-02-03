package main

import (
	"fmt"
	fmt2 "github.com/ArtisanCloud/PowerLibs/v2/fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/ArtisanCloud/PowerSocialite/v2/src/providers"
)

func main() {

	fmt.Printf("hello Socialite! \n")
	provider := providers.NewWeCom(&object.HashMap{})
	fmt2.Dump(provider.GetConfig())

}
