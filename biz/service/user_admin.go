package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/module"
	"LearnShare/biz/model/user"
	"LearnShare/pkg/errno"
	"LearnShare/pkg/utils"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserAdminService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserAdminService(ctx context.Context, c *app.RequestContext) *UserAdminService {
	return &UserAdminService{ctx: ctx, c: c}
}

func (s *UserAdminService) AdminAddUser(req *user.AdminAddUserReq) (int64, error) {
	// 验证用户名
	if valid, err := utils.VerifyUsername(req.Username); !valid {
		return 0, err
	}

	// 验证密码
	if valid, err := utils.VerifyPassword(req.Password); !valid {
		return 0, err
	}

	// 验证邮箱
	if valid, err := utils.VerifyEmail(req.Email); !valid {
		return 0, err
	}

	// 加密密码
	passwordHash, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return 0, err
	}

	// 创建用户
	userID, err := db.AdminCreateUser(s.ctx, req.Username, passwordHash, req.Email, req.RoleID, req.Status)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *UserAdminService) AdminUpdateUser(req *user.AdminUpdateUserReq) error {
	// 检查用户是否存在
	_, err := db.GetUserByID(s.ctx, req.UserID)
	if err != nil {
		return errno.NewErrNo(errno.ServiceUserNotExist, "用户不存在")
	}

	// 如果需要更新密码，先加密
	var passwordHash *string
	if req.Password != nil {
		if valid, err := utils.VerifyPassword(*req.Password); !valid {
			return err
		}
		hash, err := utils.EncryptPassword(*req.Password)
		if err != nil {
			return err
		}
		passwordHash = &hash
	}

	// 更新用户信息
	err = db.AdminUpdateUser(s.ctx, req.UserID, req.Username, passwordHash, req.Email,
		req.CollegeID, req.MajorID, req.RoleID, req.Status)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserAdminService) AdminGetUserList(req *user.AdminGetUserListReq) ([]*module.User, int64, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取用户列表
	userList, total, err := db.AdminGetUserList(s.ctx, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, err
	}

	var result []*module.User
	for _, u := range userList {
		result = append(result, u.ToUserModule())
	}

	return result, total, nil
}
