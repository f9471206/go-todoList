package services

import (
	"context"
	"fmt"
	"todolist/models"
	"todolist/repositories/base"
	"todolist/repositories/interfaces"
	"todolist/utils"

	"gorm.io/gorm"
)

type MemberService struct {
	ctx  context.Context
	repo interfaces.AuthRepository
}

func NewMemberService(ctx context.Context, repo interfaces.AuthRepository) *MemberService {
	return &MemberService{
		ctx:  ctx,
		repo: repo,
	}
}

func (s *MemberService) Index(db *gorm.DB, keyword string, page int, pageSize int, orderBy []string) (*utils.PaginatedResult[*models.User], error) {
	query := db.Model(&models.User{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	// 加入關聯查詢
	query = query.Preload("Roles")

	list, total, err := s.repo.FindAllWithQuery(s.ctx, query, page, pageSize, orderBy...)

	if err != nil {
		return nil, err
	}

	return utils.NewPaginatedResult(list, total, page, pageSize), nil
}

func (s *MemberService) Show(db *gorm.DB, id int) (*models.User, error) {
	opts := &base.FindOptions{
		Debug:         true,
		PreloadFields: []string{"Roles"},
	}
	return s.repo.FindByID(s.ctx, db, id, opts)
}

func (s *MemberService) Edit(db *gorm.DB, userID int, roleID int) (*models.User, error) {
	var user models.User
	var role models.Role

	err := db.Transaction(func(tx *gorm.DB) error {
		// 取角色，確保存在
		if err := tx.First(&role, roleID).Error; err != nil {
			return fmt.Errorf("找不到 role: %w", err)
		}

		// 取使用者
		if err := tx.Preload("Roles").First(&user, userID).Error; err != nil {
			return fmt.Errorf("找不到 user: %w", err)
		}

		// 替換使用者的角色，等於只保留這個角色
		if err := tx.Model(&user).Association("Roles").Replace(&role); err != nil {
			return fmt.Errorf("更新角色失敗: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 回傳完整資料（含角色）
	if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
