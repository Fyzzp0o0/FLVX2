package model

// User 对应 "user" 表(表名/列名均为 PG 保留字,SQL 必须带双引号)
type User struct {
	ID            int64  `json:"id"`
	User          string `json:"user"`
	Pwd           string `json:"-"`
	RoleID        int64  `json:"roleId"`
	ExpTime       int64  `json:"expTime"`
	Flow          int64  `json:"flow"` // GB
	InFlow        int64  `json:"inFlow"`
	OutFlow       int64  `json:"outFlow"`
	FlowResetTime int64  `json:"flowResetTime"`
	Num           int64  `json:"num"`
	CreatedTime   int64  `json:"createdTime"`
	UpdatedTime   int64  `json:"updatedTime"`
	Status        int64  `json:"status"`
}

// ViteConfig 对应 vite_config 表(键值配置)
type ViteConfig struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Time  int64  `json:"time"`
}
