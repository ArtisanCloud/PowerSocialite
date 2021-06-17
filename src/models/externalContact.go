package models

import (
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/golang-module/carbon"
)

// https://open.work.weixin.qq.com/api/doc/90000/90135/92114


type ExternalContact struct {
	ExternalUserID   string           `json:"external_userid"` // woAJ2GCAAAXtWyujaWJHDDGi0mACHAAA",
	Name             string           `json:"name"`            // 李四",
	Position         string           `json:"position"`        // Manager",
	Avatar           string           `json:"avatar"`          // http://p.qlogo.cn/bizmail/IcsdgagqefergqerhewSdage/0",
	CorpName         string           `json:"corp_name"`       // 腾讯",
	CorpFullName     string           `json:"corp_full_name"`  // 腾讯科技有限公司",
	Type             int              `json:"type"`            // ,
	Gender           int              `json:"gender"`          // ,
	UnionID          string           `json:"unionid"`         // ozynqsulJFCZ2z1aYeS8h-nuasdAAA",
	ExternalProfiles *ExternalProfile `json:"external_profile"`
	FollowUsers      *FollowUser      `json:"follow_user"`
}

type ExternalProfile struct {
	ExternalAttr map[string]*object.Collection `json:"external_attr"`
}

type FollowUser struct {
	UserID         string        `json:"userid"`      // tommy",
	Remark         string        `json:"remark"`      // :"李总",
	Description    string        `json:"description"` // 采购问题咨询",
	CreateTime     carbon.Carbon `json:"createtime"`  // :1525881637,
	State          string        `json:"state"`       // 外联二维码1",
	Tags           []*Tag        `json:"tags"`
	RemarkCorpName string        `json:"remark_corp_name"` //"腾讯科技",
	RemarkMobiles  RemarkMobiles `json:"remark_mobiles"`
	OperUserID     string        `json:"oper_userid"` // :"woAJ2GCAAAXtWyujaWJHDDGi0mACHAAA",
	AddWay         int           `json:"add_way"`     // 3
}

type RemarkMobiles struct {
	Mobiles []string
}
