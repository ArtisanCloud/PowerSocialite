package models

import "github.com/ArtisanCloud/PowerLibs/v3/object"

const EMPLOYEE_STATUS_ACTIVE = 1     // 已激活
const EMPLOYEE_STATUS_PROHIBITED = 2 // 已禁用
const EMPLOYEE_STATUS_INACTIVE = 4   // 未激活
const EMPLOYEE_STATUS_QUIT = 5       // 退出企业

type Employee struct {
	Alias           string          `json:"alias"`
	Avatar          string          `json:"avatar"`
	Department      []int           `json:"department"`
	Email           string          `json:"email"`
	Enable          int             `json:"enable"`
	EnglishName     string          `json:"english_name"`
	ExtAttr         *object.HashMap `json:"extattr"`
	ExternalProfile *object.HashMap `json:"external_profile"`
	Gender          string          `json:"gender"`
	HideMobile      int             `json:"hide_mobile"`
	IsLeaderInDept  []int           `json:"is_leader_in_dept"`
	IsLeader        int             `json:"isleader"`
	MainDepartment  int             `json:"main_department"`
	Mobile          string          `json:"mobile"`
	Name            string          `json:"name"`
	Order           []int           `json:"order"`
	Position        string          `json:"position"`
	QrCode          string          `json:"qr_code"`
	Status          int             `json:"status"`
	Telephone       string          `json:"telephone"`
	ThumbAvatar     string          `json:"thumb_avatar"`
	UserID          string          `json:"userid"`
	OpenUserID      string          `json:"open_userid"`
	OpenID          string          `json:"-"`
	CorpID          string          `json:"-"`
}
