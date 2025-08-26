package user

import "github.com/moweilong/chunyu/pkg/orm"

type User struct {
	orm.Model
	// Id           string      `json:"id" gorm:"primaryKey"`
	Username     string      `json:"username" gorm:"column:username;unique;notNull;size:64;default:'';comment:login name, cannot rename"`
	Nickname     string      `json:"nickname" gorm:"column:nickname;unique;notNull;size:64;default:'';comment:display name, chinese name"`
	Password     string      `json:"password" gorm:"column:password;notNull;size:128;default:'';comment:password"`
	Phone        string      `json:"phone" gorm:"column:phone;notNull;size:16;default:'';comment:phone number"`
	Email        string      `json:"email" gorm:"column:email;notNull;size:64;default:'';comment:email address"`
	Portrait     string      `json:"portrait" gorm:"column:portrait;notNull;size:255;default:'';comment:portrait image url"`
	Roles        string      `json:"roles" gorm:"column:roles;notNull;size:255;default:'';comment:Admin | Standard | Guest, split by space"`   // 这个字段写入数据库
	Contacts     orm.JSONObj `json:"contacts" gorm:"column:contacts;type:varchar(1024);comment:json e.g. {wecom:xx, dingtalk_robot_token:yy}"` // 内容为 map[string]string 结构
	Maintainer   int         `json:"maintainer" gorm:"column:maintainer;size:1;notNull;default:0;comment:是否给管理员发消息 0:not send 1:send"`         // 是否给管理员发消息 0:not send 1:send
	Belong       string      `json:"belong" gorm:"column:belong;size:191;notNull;default:'';"`
	LastActiveAt orm.Time    `json:"last_active_at" gorm:"column:last_active_at;comment:最后活跃时间"`
	CreateBy     string      `json:"create_by" gorm:"column:create_by;size:64;notNull;default:'';comment:创建人"`
	UpdateBy     string      `json:"update_by" gorm:"column:update_by;size:64;notNull;default:'';comment:更新人"`
}

func (u *User) TableName() string {
	return "users"
}
