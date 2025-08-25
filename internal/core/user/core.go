package user

import "github.com/moweilong/chunyu/domain/uniqueid"

// Storer data persistence
// 相当于简洁架构的 Store，定义了数据层的方法
type Storer interface {
	User() UserStore
}

// Core business domain
// 相当于简洁架构的 Biz，依赖 store，可扩展其他依赖
type Core struct {
	store    Storer
	uniqueID uniqueid.Core
}

// NewCore create business domain
func NewCore(store Storer, uni uniqueid.Core) *Core {
	return &Core{
		store:    store,
		uniqueID: uni,
	}
}
