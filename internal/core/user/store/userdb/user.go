package userdb

import (
	"context"

	"github.com/moweilong/chunyu/internal/core/user"
	"github.com/moweilong/chunyu/pkg/orm"
	"gorm.io/gorm"
)

var _ user.UserStore = User{}

// User Related business namespaces
type User DB

// NewUser instance object
func NewUser(db *gorm.DB) User {
	return User{db: db}
}

// Find implements user.UserStore.
func (s User) Find(ctx context.Context, bs *[]*user.User, page orm.Pager, opts ...orm.QueryOption) (int64, error) {
	return orm.FindWithContext(ctx, s.db, bs, page, opts...)
}

// Get implements user.UserStore.
func (s User) Get(ctx context.Context, model *user.User, opts ...orm.QueryOption) error {
	return orm.FirstWithContext(ctx, s.db, model, opts...)
}

// Add implements user.UserStore.
func (s User) Add(ctx context.Context, model *user.User) error {
	return s.db.WithContext(ctx).Create(model).Error
}

// Edit implements user.UserStore.
func (s User) Edit(ctx context.Context, model *user.User, changeFn func(*user.User), opts ...orm.QueryOption) error {
	return orm.UpdateWithContext(ctx, s.db, model, changeFn, opts...)
}

// Del implements user.UserStore.
func (s User) Del(ctx context.Context, model *user.User, opts ...orm.QueryOption) error {
	return orm.DeleteWithContext(ctx, s.db, model, opts...)
}
