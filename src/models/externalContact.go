package models

// https://open.work.weixin.qq.com/api/doc/90000/90135/92114

type Text struct {
	Value string `json:"value"`
}

type MiniProgram struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
	Title    string `json:"title"`
}

type Web struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type ExternalAttr struct {
	Type        int          `json:"type"`
	Name        string       `json:"name"`
	Text        *Text        `json:"text,omitempty"`
	Web         *Web         `json:"web,omitempty"`
	MiniProgram *MiniProgram `json:"miniprogram,omitempty"`
}

type ExternalProfile struct {
	ExternalAttr []*ExternalAttr `json:"external_attr"`
}

type ExternalContact struct {
	ExternalUserID  string           `json:"external_userid"`
	Name            string           `json:"name"`
	Position        string           `json:"position"`
	Avatar          string           `json:"avatar"`
	CorpName        string           `json:"corp_name"`
	CorpFullName    string           `json:"corp_full_name"`
	Type            int              `json:"type"`
	Gender          int              `json:"gender"`
	UnionID         string           `json:"unionid"`
	ExternalProfile *ExternalProfile `json:"external_profile"`
}

type WechatChannel struct {
	NickName string `json:"nickname"`
	Source   int    `json:"source"`
}

type FollowUser struct {
	UserID         string         `json:"userid"`
	Remark         string         `json:"remark"`
	Description    string         `json:"description"`
	CreateTime     int            `json:"createtime"`
	TagIDs         []string       `json:"tag_id,omitempty"`
	Tags           []Tag          `json:"tags,omitempty"`
	RemarkCorpName string         `json:"remark_corp_name,omitempty"`
	RemarkMobiles  []string       `json:"remark_mobiles,omitempty"`
	OperUserID     string         `json:"oper_userid"`
	AddWay         int            `json:"add_way"`
	WechatChannels *WechatChannel `json:"wechat_channels,omitempty"`
	State          string         `json:"state,omitempty"`
}
