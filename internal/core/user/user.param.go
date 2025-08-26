package user

import (
	"github.com/moweilong/chunyu/pkg/orm"
	"github.com/moweilong/chunyu/pkg/web"
)

type FindUserInput struct {
	web.PagerFilter
	Key string `form:"key"` // 昵称/手机号码/邮箱 模糊搜索
}

type AddUserInput struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Nickname string      `json:"nickname"`
	Phone    string      `json:"phone"`
	Email    string      `json:"email"`
	Portrait string      `json:"portrait"`
	Roles    []string    `json:"roles"`
	Contacts orm.JSONObj `json:"contacts"`
}

type EditUserInput struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Nickname string      `json:"nickname"`
	Phone    string      `json:"phone"`
	Email    string      `json:"email"`
	Roles    []string    `json:"roles"`
	Contacts orm.JSONObj `json:"contacts"`
}
