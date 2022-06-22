package models

type Tag struct {
	GroupName string `json:"group_name"`
	TagName   string `json:"tag_name"`
	TagID     string `json:"tag_id"`
	Type      int    `json:"type"`
}
