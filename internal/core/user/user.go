package user

import (
	"context"
	"log/slog"

	"github.com/jinzhu/copier"
	"github.com/moweilong/chunyu/pkg/orm"
	"github.com/moweilong/chunyu/pkg/reason"
)

// UserStore data persistence
type UserStore interface {
	Find(context.Context, *[]*User, orm.Pager, ...orm.QueryOption) (int64, error)
	Get(context.Context, *User, ...orm.QueryOption) error
	Add(context.Context, *User) error
	Edit(context.Context, *User, func(*User), ...orm.QueryOption) error
	Del(context.Context, *User, ...orm.QueryOption) error
}

// FindUser Paginated search
func (c *Core) FindUser(ctx context.Context, in *FindUserInput) ([]*User, int64, error) {
	items := make([]*User, 0)

	query := orm.NewQuery(1)
	query.OrderBy("id ASC")
	if in.Key != "" {
		query.Where("nickname LIKE ? OR phone LIKE ? OR email LIKE ?", "%"+in.Key+"%", "%"+in.Key+"%", "%"+in.Key+"%")
	}
	total, err := c.store.User().Find(ctx, &items, in, query.Encode()...)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetUser Query a single object
func (c *Core) GetUser(ctx context.Context, id string) (*User, error) {
	var out User
	if err := c.store.User().Get(ctx, &out, orm.Where("id = ?", id)); err != nil {
		if orm.IsErrRecordNotFound(err) {
			return nil, reason.ErrNotFound.Withf(`Get err[%s]`, err.Error())
		}
		return nil, reason.ErrDB.Withf(`Get err[%s]`, err.Error())
	}
	return &out, nil
}

// AddUser  Insert into database
func (c *Core) AddUser(ctx context.Context, in *AddUserInput) (*User, error) {
	var out User
	if err := copier.Copy(&out, in); err != nil {
		slog.ErrorContext(ctx, "Copy", "err", err)
	}
	if err := c.store.User().Add(ctx, &out); err != nil {
		return nil, reason.ErrDB.Withf(`Add err[%s]`, err.Error())
	}
	return &out, nil
}

// EditUser Update object information
func (c *Core) EditUser(ctx context.Context, id int, in *EditUserInput) (*User, error) {
	var out User
	if err := c.store.User().Edit(ctx, &out, func(b *User) {
		if err := copier.Copy(b, in); err != nil {
			slog.ErrorContext(ctx, "Copy", "err", err)
		}
	}, orm.Where("id = ?", id)); err != nil {
		return nil, reason.ErrDB.Withf(`Edit err[%s]`, err.Error())
	}
	return &out, nil
}

// DelUser Delete object
func (c *Core) DelUser(ctx context.Context, id string) (*User, error) {
	var out User
	if err := c.store.User().Del(ctx, &out, orm.Where("id = ?", id)); err != nil {
		return nil, reason.ErrDB.Withf(`Del err[%s]`, err.Error())
	}
	return &out, nil
}

// GetUserByUsername 根据用户名查询用户
func (c *Core) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var out User
	if err := c.store.User().Get(ctx, &out, orm.Where("username = ?", username)); err != nil {
		if orm.IsErrRecordNotFound(err) {
			return nil, reason.ErrNotFound.Withf(`Get err[%s]`, err.Error())
		}
		return nil, reason.ErrDB.Withf(`Get err[%s]`, err.Error())
	}
	return &out, nil
}
