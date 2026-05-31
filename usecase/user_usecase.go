package usecase

import (
	"context"
	"fmt"
	"golang-docker-postgres/domain"
)

// userUsecase 結構體內部持有的是「介面（規格）」，而不是特定的資料庫物件
type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase 初始化這個 Usecase 的函式
func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: ur}
}

// GetProfile 實作業務邏輯
func (u *userUsecase) GetProfile(ctx context.Context, id int64) (*domain.User, error) {
	// 這裡可以寫你的業務邏輯，例如：如果 id <= 0 直接回傳空資料，不用白費力氣連資料庫
	if id <= 0 {
		return nil, nil
	}
	
	// 邏輯檢查完畢，呼叫 repository 去拿資料
	return u.userRepo.GetByID(ctx, id)
}
// ... 原本的 GetProfile 保持不動 ...

// Store 處理商業邏輯（例如名字檢查）
func (u *userUsecase) Store(ctx context.Context, user *domain.User) error {
	// 商業邏輯檢查：如果名字是空的，直接拒絕，不浪費效能連資料庫
	if user.Name == "" {
		return fmt.Errorf("使用者姓名不能為空") 
	}
	// 如果需要，要在這裡引進 "fmt" 套件到 import 喔！

	// 檢查過關，交給 Repository 存進資料庫
	return u.userRepo.Store(ctx, user)
}