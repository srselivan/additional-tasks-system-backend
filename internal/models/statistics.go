package models

type Statistics struct {
	UserName  string `json:"userName"`
	GroupName string `json:"groupName"`
	Score     int64  `json:"score"`
}
