package user

import "github.com/moweilong/chunyu/pkg/orm"

type User struct {
	Id             int64           `json:"id" gorm:"primaryKey"`
	Username       string          `json:"username"`
	Nickname       string          `json:"nickname"`
	Password       string          `json:"-"`
	Phone          string          `json:"phone"`
	Email          string          `json:"email"`
	Portrait       string          `json:"portrait"`
	Roles          string          `json:"-"`              // 这个字段写入数据库
	RolesLst       []string        `json:"roles" gorm:"-"` // 这个字段和前端交互
	TeamsLst       []int64         `json:"-" gorm:"-"`     // 这个字段方便映射团队，前端和数据库都不用到
	Contacts       orm.JSONObj     `json:"contacts"`       // 内容为 map[string]string 结构
	Maintainer     int             `json:"maintainer"`     // 是否给管理员发消息 0:not send 1:send
	CreateAt       int64           `json:"create_at"`
	CreateBy       string          `json:"create_by"`
	UpdateAt       int64           `json:"update_at"`
	UpdateBy       string          `json:"update_by"`
	Belong         string          `json:"belong"`
	Admin          bool            `json:"admin" gorm:"-"` // 方便前端使用
	UserGroupsRes  []*UserGroupRes `json:"user_groups" gorm:"-"`
	BusiGroupsRes  []*BusiGroupRes `json:"busi_groups" gorm:"-"`
	LastActiveTime int64           `json:"last_active_time"`
}

type UserGroupRes struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type BusiGroupRes struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (u *User) TableName() string {
	return "users"
}
