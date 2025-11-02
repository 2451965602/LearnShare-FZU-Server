package db

import (
	"LearnShare/biz/model/module"
	"time"
)

type User struct {
	UserID          int64     `json:"user_id" db:"user_id"`
	Username        string    `json:"username" db:"username"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	Email           string    `json:"email" db:"email"`
	CollegeID       *int64    `json:"college_id,omitempty" db:"college_id"`
	MajorID         *int64    `json:"major_id,omitempty" db:"major_id"`
	AvatarURL       *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	ReputationScore int64     `json:"reputation_score" db:"reputation_score"`
	RoleID          int64     `json:"role_id" db:"role_id"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func (u User) ToUserModule() *module.User {
	user := &module.User{
		UserId:          u.UserID,
		Username:        u.Username,
		Email:           u.Email,
		ReputationScore: u.ReputationScore,
		RoleId:          u.RoleID,
		Status:          u.Status,
	}

	if u.AvatarURL != nil {
		user.AvatarUrl = *u.AvatarURL
	}
	if u.CollegeID != nil {
		user.CollegeId = *u.CollegeID
	}
	if u.MajorID != nil {
		user.MajorId = *u.MajorID
	}
	return user
}
