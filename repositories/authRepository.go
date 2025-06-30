package repositories

import (
	"context"
	"todolist/models"
	"todolist/repositories/base"

	"gorm.io/gorm"
)

type AuthRepository struct {
	*base.BaseRepository[*models.User]
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		BaseRepository: base.NewBaseRepository[*models.User](), // 回傳 *BaseRepository[T]
	}
}

// 檢查帳號是否存在
func IsAccountExist(ctx context.Context, db *gorm.DB, account string) (bool, error) {
	var count int64
	if err := db.WithContext(ctx).
		Model(&models.User{}).
		Where("account = ?", account).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
