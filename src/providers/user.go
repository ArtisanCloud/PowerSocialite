package providers

import (
	"encoding/json"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerSocialite/v3/src/contracts"
)

type User struct {
	contracts.UserInterface

	*object.Attribute
	provider *ProviderInterface
}

func NewUser(attributes *object.HashMap, provider ProviderInterface) *User {
	return &User{
		Attribute: object.NewAttribute(attributes),
		provider:  &provider,
	}
}

func (user *User) GetID() string {
	if user.Attributes["id"] != nil {
		return user.Attributes["id"].(string)
	} else {
		return user.GetEmail()
	}
}

func (user *User) GetOpenID() string {
	if user.Attributes["openID"] != nil {
		return user.Attributes["openID"].(string)
	} else {
		return ""
	}
}

func (user *User) GetExternalUserID() string {
	if user.Attributes["externalUserID"] != nil {
		return user.Attributes["externalUserID"].(string)
	} else {
		return ""
	}
}

func (user *User) GetDeviceID() string {
	if user.Attributes["deviceID"] != nil {
		return user.Attributes["deviceID"].(string)
	} else {
		return ""
	}
}

func (user *User) GetNickname() string {
	if user.Attributes["nickname"] != nil {
		return user.Attributes["nickname"].(string)
	} else {
		return user.GetName()
	}
}

func (user *User) GetMobile() string {
	if user.Attributes["mobile"] != nil {
		return user.Attributes["mobile"].(string)
	}
	return ""
}

func (user *User) GetName() string {
	if user.Attributes["name"] != nil {
		return user.Attributes["name"].(string)
	}
	return ""
}

func (user *User) GetEmail() string {
	if user.Attributes["email"] != nil {
		return user.Attributes["email"].(string)
	}
	return ""
}

func (user *User) GetAvatar() string {
	if user.Attributes["avatar"] != nil {
		return user.Attributes["avatar"].(string)
	}
	return ""
}

func (user *User) SetAccessToken(token string) *User {
	user.SetAttribute("access_token", token)

	return user
}
func (user *User) GetAccessToken() string {
	return user.GetAttribute("access_token", "").(string)
}

func (user *User) SetRefreshToken(refreshToken string) *User {
	user.SetAttribute("refresh_token", refreshToken)

	return user
}
func (user *User) GetRefreshToken() string {
	return user.GetAttribute("refresh_token", "").(string)
}

func (user *User) SetExpiresIn(expiresIn float64) *User {
	user.SetAttribute("expires_in", expiresIn)

	return user
}
func (user *User) GetExpiresIn() int {
	return user.GetAttribute("expires_in", 0).(int)
}

func (user *User) SetRaw(u *object.HashMap) *User {
	// copy object from raw u
	raw, _ := object.StructToHashMap(u)

	user.SetAttribute("raw", raw)

	return user
}
func (user *User) GetRaw() (raw *object.HashMap, err error) {
	if user.GetAttribute("raw", nil) != nil {
		strRaw := user.GetAttribute("raw", nil).(string)
		raw = &object.HashMap{}
		err = object.JsonDecode([]byte(strRaw), raw)
		if err != nil {
			return nil, err
		}

		return raw, nil
	}
	return nil, nil

}

func (user *User) SetTokenResponse(response *object.HashMap) *User {
	user.SetAttribute("token_response", response)

	return user
}
func (user *User) GetTokenResponse() *object.HashMap {
	rsToken := user.GetAttribute("token_response", nil)
	switch rs := rsToken.(type) {
	case *object.HashMap:
		return rs
	case string:
		mapToken := &object.HashMap{}
		err := object.JsonDecode([]byte(rs), mapToken)
		if err != nil {
			println(err)
			return nil
		}
		return mapToken
	default:
		return nil
	}

}

func (user *User) JsonSerialize() *object.HashMap {
	return &user.Attributes
}

func (user *User) Serialize() string {
	buffer, err := json.Marshal(user.Attributes)
	if err != nil {
		return ""
	}
	return string(buffer)
}

func (user *User) UnSerialize(serialized string) *object.HashMap {

	return user.GetAttributes()
}

func (user *User) GetProvider() *ProviderInterface {
	return user.provider
}

func (user *User) SetProvider(provider ProviderInterface) *User {
	user.provider = &provider
	return user
}

func (user *User) IsSnapShotUser() bool {
	return user.GetAttribute("is_snapshotuser", false).(bool)
}

func (user *User) SetSnapShotUser(IsSnapShotUser bool) *User {
	user.SetAttribute("is_snapshotuser", IsSnapShotUser)
	return user
}
