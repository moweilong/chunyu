package api

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/moweilong/chunyu/domain/uniqueid"
	"github.com/moweilong/chunyu/internal/core/user"
	"github.com/moweilong/chunyu/internal/core/user/store/userdb"
	"github.com/moweilong/chunyu/pkg/orm"
	"github.com/moweilong/chunyu/pkg/web"
	"gorm.io/gorm"
)

type UserAPI struct {
	userCore *user.Core
}

func NewUserAPI(db *gorm.DB, uni uniqueid.Core) UserAPI {
	core := user.NewCore(userdb.NewDB(db).AutoMigrate(orm.GetEnabledAutoMigrate()), uni)
	return UserAPI{
		userCore: core,
	}
}

func registerUser(g gin.IRouter, api UserAPI, handler ...gin.HandlerFunc) {
	{
		group := g.Group("/user", handler...)
		group.GET("", web.WarpH(api.findUser))
		group.GET("/:id", web.WarpH(api.getUser))
		group.POST("", web.WarpH(api.addUser))
		group.PUT("/:id", web.WarpH(api.editUser))
		group.DELETE("/:id", web.WarpH(api.delUser))
	}
}

// >>> user >>>>>>>>>>>>>>>>>>>>>>>>>>>
func (h *UserAPI) findUser(c *gin.Context, in *user.FindUserInput) (any, error) {
	items, total, err := h.userCore.FindUser(c.Request.Context(), in)
	return gin.H{"items": items, "total": total}, err
}

func (h *UserAPI) getUser(c *gin.Context, _ *struct{}) (any, error) {
	userID := c.Param("id")
	return h.userCore.GetUser(c.Request.Context(), userID)
}

func (h *UserAPI) editUser(c *gin.Context, in *user.EditUserInput) (any, error) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, err
	}
	return h.userCore.EditUser(c.Request.Context(), userID, in)
}

func (h *UserAPI) addUser(c *gin.Context, in *user.AddUserInput) (any, error) {
	return h.userCore.AddUser(c.Request.Context(), in)
}

func (h *UserAPI) delUser(c *gin.Context, _ *struct{}) (any, error) {
	userID := c.Param("id")
	return h.userCore.DelUser(c.Request.Context(), userID)
}

// initRoot 初始化管理员密码
func (h *UserAPI) InitRoot(c *gin.Context, _ *struct{}) bool {
	root, err := h.userCore.GetUserByUsername(c.Request.Context(), "root")
	if err != nil {
		os.Exit(1)
	}
	if root == nil {
		return false
	}
	if len(root.Password) > 31 {
		// already done before
		return false
	}
	// 加密密码
	encryptedPassword, err := web.Encrypt(root.Password)
	if err != nil {
		os.Exit(1)
	}
	root.Password = encryptedPassword
	// 更新用户密码
	in := user.EditUserInput{
		Password: encryptedPassword,
	}
	_, err = h.userCore.EditUser(c.Request.Context(), root.ID, &in)
	if err != nil {
		os.Exit(1)
	}
	return true
}
