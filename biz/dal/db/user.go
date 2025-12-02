package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
	"errors"

	"gorm.io/gorm"
)

// CreateUser 创建新用户
func CreateUser(ctx context.Context, username, passwordHash, email string) error {
	user := &User{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
		RoleID:       2, // 默认普通用户角色ID
		Status:       "inactive",
	}

	err := DB.WithContext(ctx).Table(constants.UserTableName).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errno.NewErrNo(errno.ServiceUserExist, "用户已存在")
		}
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建用户失败: "+err.Error())
	}
	return nil
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(ctx context.Context, userID int64, newPasswordHash string) error {
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", userID).Update("password_hash", newPasswordHash).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新用户密码失败: "+err.Error())
	}
	return nil
}

// UpdateUserPasswordAsync 异步更新用户密码
func UpdateUserPasswordAsync(ctx context.Context, userID int64, newPasswordHash string) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateUserPassword(ctx, userID, newPasswordHash)
	})
}

// UpdateMajorID 更新用户专业ID
func UpdateMajorID(ctx context.Context, userID int, majorID int) error {
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", userID).Update("major_id", majorID).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新用户专业失败: "+err.Error())
	}
	return nil
}

// UpdateMajorIDAsync 异步更新用户专业ID
func UpdateMajorIDAsync(ctx context.Context, userID int, majorID int) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateMajorID(ctx, userID, majorID)
	})
}

func UpdateAvatarURL(ctx context.Context, userID int64, avatarURL string) error {
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", userID).Update("avatar_url", avatarURL).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新用户头像失败: "+err.Error())
	}
	return nil
}

// UpdateAvatarURLAsync 异步更新用户头像
func UpdateAvatarURLAsync(ctx context.Context, userID int64, avatarURL string) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateAvatarURL(ctx, userID, avatarURL)
	})
}

// UpdateUserStatues 更新用户状态
func UpdateUserStatues(ctx context.Context, userID int64, newStatus string) error {
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", userID).Update("status", newStatus).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新用户状态失败: "+err.Error())
	}
	return nil
}

// UpdateUserStatuesAsync 异步更新用户状态
func UpdateUserStatuesAsync(ctx context.Context, userID int64, newStatus string) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateUserStatues(ctx, userID, newStatus)
	})
}

// UpdateUserEmail 更新用户邮箱
func UpdateUserEmail(ctx context.Context, userID int, newEmail string) error {
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", userID).Update("email", newEmail).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新用户邮箱失败: "+err.Error())
	}
	return nil
}

// UpdateUserEmailAsync 异步更新用户邮箱
func UpdateUserEmailAsync(ctx context.Context, userID int, newEmail string) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateUserEmail(ctx, userID, newEmail)
	})
}

// GetUserByEmail 根据邮箱查询用户
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询用户失败")
	}
	return &user, nil
}

// GetUserByID 根据用户ID查询用户
func GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", id).First(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询用户失败")
	}
	return &user, nil
}

// AdminCreateUser 管理员创建新用户
func AdminCreateUser(ctx context.Context, username, passwordHash, email string, roleID int64, status string) (int64, error) {

	user := &User{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
		RoleID:       roleID,
		Status:       status,
	}

	err := DB.WithContext(ctx).Table(constants.UserTableName).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, errno.NewErrNo(errno.ServiceUserExist, "用户已存在")
		}
		return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建用户失败: "+err.Error())
	}
	return user.UserID, nil
}

// AdminUpdateUser 管理员更新用户信息
func AdminUpdateUser(ctx context.Context, userID int64, username, passwordHash, email, collegeID, majorID *string, roleID *int64, status *string) error {
	updates := make(map[string]interface{})

	if username != nil {
		updates["username"] = *username
	}
	if passwordHash != nil {
		updates["password_hash"] = *passwordHash
	}
	if email != nil {
		updates["email"] = *email
	}
	if collegeID != nil {
		updates["college_id"] = *collegeID
	}
	if majorID != nil {
		updates["major_id"] = *majorID
	}
	if roleID != nil {
		updates["role_id"] = *roleID
	}
	if status != nil {
		updates["status"] = *status
	}

	if len(updates) == 0 {
		return nil
	}

	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ?", userID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新用户信息失败: "+err.Error())
	}
	return nil
}

// IncrementUserReputation 增加用户信誉分
func IncrementUserReputation(ctx context.Context, userID int64, delta int64) error {
	err := DB.WithContext(ctx).Table(constants.UserTableName).Where("user_id = ? AND reputation_score < 100", userID).Update("reputation_score", gorm.Expr("reputation_score + ?", delta)).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新用户信誉分失败: "+err.Error())
	}
	return nil
}

// IncrementUserReputationAsync 异步增加用户信誉分
func IncrementUserReputationAsync(ctx context.Context, userID int64, delta int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return IncrementUserReputation(ctx, userID, delta)
	})
}

func AdminGetUserList(ctx context.Context, pageNum, pageSize int) ([]*User, int64, error) {
	var users []*User
	var total int64

	err := DB.WithContext(ctx).Table(constants.UserTableName).Count(&total).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询用户总数失败: "+err.Error())
	}

	err = DB.WithContext(ctx).Table(constants.UserTableName).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询用户列表失败: "+err.Error())
	}

	return users, total, nil
}
