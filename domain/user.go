package domain

import "context"

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	// 🆕 新增規格：儲存使用者
	Store(ctx context.Context, user *User) error
}

type UserUsecase interface {
	GetProfile(ctx context.Context, id int64) (*User, error)
	// 🆕 新增規格：註冊/創建使用者
	Store(ctx context.Context, user *User) error
}