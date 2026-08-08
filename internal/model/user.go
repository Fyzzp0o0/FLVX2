package model

// User 对应 "user" 表(表名/列名均为 PG 保留字,SQL 必须带双引号)
type User struct {
	ID            int64  `json:"id"`
	User          string `json:"user"`
	Pwd           string `json:"-"`
	RoleID        int64  `json:"role_id"`
	ExpTime       int64  `json:"exp_time"`
	Flow          int64  `json:"flow"` // GB
	InFlow        int64  `json:"in_flow"`
	OutFlow       int64  `json:"out_flow"`
	FlowResetTime int64  `json:"flow_reset_time"`
	Num           int64  `json:"num"`
	CreatedTime   int64  `json:"created_time"`
	UpdatedTime   int64  `json:"updated_time"`
	Status        int64  `json:"status"`
}

// ViteConfig 对应 vite_config 表(键值配置)
type ViteConfig struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Time  int64  `json:"time"`
}
