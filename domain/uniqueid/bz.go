// uniqueid
// 的设计是用于生成全局唯一的 ID，避免重复。
// 此库不考虑分布式，仅通过数据库主键索引来实现。
// 当 id 重复时，由业务端抛出错误即可

package uniqueid

import (
	"context"
	"crypto/rand"
	"log/slog"
	"math/big"
	"time"

	"github.com/moweilong/chunyu/pkg/hook"
	"github.com/moweilong/chunyu/pkg/orm"
)

const (
	// 删除 o 和 i 的字符集，避免视觉混淆
	LetterBytes36NoOI = "abcdefghjklmnpqrstuvwxyz0123456789"

	LetterBytes72      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	LetterBytes36      = "abcdefghijklmnopqrstuvwxyz0123456789" // default
	LetterBytes36Upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type IDManager struct {
	store UniqueIDStorer
	// TODO: 可以初始化时读取数据库内的数量，判断重复因子，从而减少尝试或更换策略

	letterBytes string // 随机字符串字符集
}

func NewIDManager(store UniqueIDStorer) *IDManager {
	return &IDManager{
		store:       store,
		letterBytes: LetterBytes36,
	}
}

// SetLetterBytes 设置随机字符串字符集
func (m *IDManager) SetLetterBytes(letterBytes string) {
	m.letterBytes = letterBytes
}

// UniqueID 获取唯一 id
func (m *IDManager) UniqueID(prefix string, length int) string {
	cost := hook.UseTiming(time.Second)
	defer cost()

	// 如果在最低长度中，碰撞比较频繁，增加 1 位长度再试一次
	for i := range 10 {
		// 生成自定义长度随机数，通过数据库主键来防止碰撞，碰撞后再次尝试
		for range 36 {
			id := prefix + GenerateRandomString(m.letterBytes, length+i)
			if err := m.store.Add(context.Background(), &UniqueID{ID: id}); err != nil {
				slog.Error("UniqueID", "err", err)
				continue
			}
			return id
		}
	}
	slog.Error("UniqueID", "err", "超过最大循环次数，未获取到唯一 id")
	return "unknown"
}

// UndoUniqueID 删除唯一 id
func (m *IDManager) UndoUniqueID(id string) error {
	var uni UniqueID
	return m.store.Del(context.Background(), &uni, orm.Where("id=?", id))
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(letterBytes string, length int) string {
	lettersLength := big.NewInt(int64(len(letterBytes)))
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		idx, _ := rand.Int(rand.Reader, lettersLength)
		result[i] = letterBytes[idx.Int64()]
	}
	return string(result)
}
