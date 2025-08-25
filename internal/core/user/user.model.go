package user

import "github.com/moweilong/chunyu/pkg/orm"

type User struct {
	orm.Model
	// Id           string      `json:"id" gorm:"primaryKey"`
	Username   string      `json:"username" gorm:"column:name;notNull;default:'';comment:用户名"`
	Nickname   string      `json:"nickname" gorm:"column:nickname;notNull;default:'';comment:用户昵称"`
	Password   string      `json:"-" gorm:"column:password;notNull;default:'';comment:密码"`
	Phone      string      `json:"phone" gorm:"column:phone;notNull;default:'';comment:手机号"`
	Email      string      `json:"email" gorm:"column:email;notNull;default:'';comment:邮箱"`
	Portrait   string      `json:"portrait" gorm:"column:portrait;notNull;default:'';comment:头像"`
	Roles      string      `json:"-" gorm:"column:roles;notNull;default:'';comment:角色"`                                       // 这个字段写入数据库
	Contacts   orm.JSONObj `json:"contacts" gorm:"column:contacts;notNull;default:'';comment:联系人"`                            // 内容为 map[string]string 结构
	Maintainer int         `json:"maintainer" gorm:"column:maintainer;notNull;default:0;comment:是否给管理员发消息 0:not send 1:send"` // 是否给管理员发消息 0:not send 1:send
	// CreateAt     orm.Time    `json:"create_at" gorm:"column:create_at;notNull;default:'';comment:创建时间"`
	CreateBy string `json:"create_by" gorm:"column:create_by;notNull;default:'';comment:创建人"`
	// UpdateAt     orm.Time    `json:"update_at" gorm:"column:update_at;notNull;default:'';comment:更新时间"`
	UpdateBy     string   `json:"update_by" gorm:"column:update_by;notNull;default:'';comment:更新人"`
	Belong       string   `json:"belong" gorm:"column:belong;notNull;default:'';comment:所属"`
	LastActiveAt orm.Time `json:"last_active_at" gorm:"column:last_active_at;notNull;default:'';comment:最后活跃时间"`
}

func (u *User) TableName() string {
	return "users"
}
