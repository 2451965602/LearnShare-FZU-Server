package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/dal/redis"
	"LearnShare/biz/model/module"
	"LearnShare/biz/model/user"
	"LearnShare/pkg/errno"
	"LearnShare/pkg/utils"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx: ctx, c: c}
}

func (s *UserService) Register(req *user.RegisterReq) error {

	if valid, err := utils.VerifyUsername(req.Username); !valid {
		return err
	}

	if valid, err := utils.VerifyPassword(req.Password); !valid {
		return err
	}

	if valid, err := utils.VerifyEmail(req.Email); !valid {
		return err
	}

	passwordHash, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = passwordHash

	err = db.CreateUser(s.ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) LoginIn(req *user.LoginInReq) (*module.User, error) {
	userInfo, err := db.GetUserByEmail(s.ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := utils.ComparePassword(userInfo.PasswordHash, req.Password); err != nil {
		return nil, err
	}

	userInfo.PasswordHash = ""
	return userInfo.ToUserModule(), nil
}

func (s *UserService) LoginOut() error {
	auth := string(s.c.GetHeader("Authorization"))
	if auth == "" {
		return nil
	}

	parts := strings.Fields(auth)
	if len(parts) == 0 {
		return nil
	}

	var token string
	if strings.EqualFold(parts[0], "Bearer") {
		if len(parts) < 2 {
			return nil
		}
		token = parts[1]
	} else {
		// 支持直接传 token 的情况
		token = parts[0]
	}

	if err := redis.SetBlacklistToken(s.ctx, token); err != nil {
		return err
	}

	return nil
}

func (s *UserService) SendVerifyEmail(req *user.SendVerifyEmailReq) error {
	code, err := utils.GenerateCode()
	if err != nil {
		return err
	}

	err = utils.MailSendCode(req.Email, code)
	if err != nil {
		return err
	}

	err = redis.PutCodeToCache(s.ctx, req.Email, code)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) VerifyEmail(req *user.VerifyEmailReq) error {

	userId := GetUidFormContext(s.c)

	storeCode, err := redis.GetCodeCache(s.ctx, req.Email)
	if err != nil {
		return err
	}

	if storeCode != req.Code {
		return errno.NewErrNo(errno.ServiceInvalidCode, "验证码不正确")
	}

	err = db.UpdateUserStatues(s.ctx, userId, "active")
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateEmail(req *user.UpdateEmailReq) error {

	userId := GetUidFormContext(s.c)

	storeCode, err := redis.GetCodeCache(s.ctx, req.NewEmail)
	if err != nil {
		return err
	}

	if storeCode != req.Code {
		return errno.NewErrNo(errno.ServiceInvalidCode, "验证码不正确")
	}

	err = db.UpdateUserEmail(s.ctx, int(userId), req.NewEmail)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdatePassword(req *user.UpdatePasswordReq) error {
	userId := GetUidFormContext(s.c)
	userInfo, err := db.GetUserByID(s.ctx, userId)
	if err != nil {
		return err
	}

	if err := utils.ComparePassword(userInfo.PasswordHash, req.OldPassword); err != nil {
		return errno.NewErrNo(errno.ServiceInvalidPassword, "旧密码不正确")
	}
	newPasswordHash, err := utils.EncryptPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = db.UpdateUserPassword(s.ctx, userInfo.UserID, newPasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateMajor(req *user.UpdateMajorReq) error {
	userId := GetUidFormContext(s.c)
	userInfo, err := db.GetUserByID(s.ctx, userId)
	if err != nil {
		return err
	}

	err = db.UpdateMajorID(s.ctx, int(userInfo.UserID), int(req.NewMajorId))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) uploadAvatarReq(req *user.UploadAvatarReq) (string, error) {
	return "", nil
}

func (s *UserService) ResetPassword(req *user.ResetPasswordReq) error {
	storeCode, err := redis.GetCodeCache(s.ctx, req.Email)
	if err != nil {
		return err
	}

	if storeCode != req.Code {
		return errno.NewErrNo(errno.ServiceInvalidCode, "验证码不正确")
	}

	newPasswordHash, err := utils.EncryptPassword(req.NewPassword)
	if err != nil {
		return err
	}

	userInfo, err := db.GetUserByEmail(s.ctx, req.Email)
	if err != nil {
		return err
	}

	err = db.UpdateUserPassword(s.ctx, userInfo.UserID, newPasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserInfo(req *user.GetUserInfoReq) (*module.User, error) {
	userInfo, err := db.GetUserByID(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return userInfo.ToUserModule(), nil
}
