package userdb

import (
	"context"
	"testing"

	"github.com/moweilong/chunyu/internal/core/user"
	"github.com/moweilong/chunyu/pkg/orm"
)

func TestUserGet(t *testing.T) {
	db, mock, err := generateMockDB()
	if err != nil {
		t.Fatal(err)
	}
	userDB := NewUser(db)
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id=\$1 (.+) LIMIT \$2`).WithArgs("jack", 1)
	var out user.User
	if err := userDB.Get(context.Background(), &out, orm.Where("id=?", "jack")); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("ExpectationsWereMet err:", err)
	}
}
