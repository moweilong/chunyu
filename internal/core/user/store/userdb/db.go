package userdb

import (
	"github.com/moweilong/chunyu/internal/core/user"
	"gorm.io/gorm"
)

var _ user.Storer = DB{}

// DB Related business namespaces
type DB struct {
	db *gorm.DB
}

// NewDB instance object
func NewDB(db *gorm.DB) DB {
	return DB{db: db}
}

// User Get business instance
func (d DB) User() user.UserStore {
	return User(d)
}

// AutoMigrate sync database
func (d DB) AutoMigrate(ok bool) DB {
	if !ok {
		return d
	}
	if err := d.db.AutoMigrate(
		new(user.User),
	); err != nil {
		panic(err)
	}
	return d
}
