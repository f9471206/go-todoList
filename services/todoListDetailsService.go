package services

import (
	"context"
	"errors"
	"todolist/models"
	"todolist/repositories/interfaces"

	"gorm.io/gorm"
)

type TodoListDetailsService struct {
	ctx  context.Context
	repo interfaces.TodoListDetailsRepository
}

func NewTodoListDetailsService(ctx context.Context, repo interfaces.TodoListDetailsRepository) *TodoListDetailsService {
	return &TodoListDetailsService{
		ctx:  ctx,
		repo: repo,
	}
}

func (s *TodoListDetailsService) Create(db *gorm.DB, listID int, name string, detail string, ids []int) (*models.TodoListDetails, error) {
	data := &models.TodoListDetails{
		TodoListID: listID,
		Name:       name,
		Detail:     detail,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// 檢查 listID 是否存在
		var count int64
		if err := tx.Model(&models.TodoList{}).Where("id = ?", listID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("to_do_list_type_id 不存在")
		}

		// 建立 TodoListDetails
		if err := s.repo.Create(s.ctx, tx, data); err != nil {
			return err
		}

		// 查出 User 對象並建立關聯
		var users []models.User
		if err := tx.Where("id IN ?", ids).Find(&users).Error; err != nil {
			return err
		}

		// 驗證所有 ids 都存在
		if len(users) != len(ids) {
			return errors.New("部分 User ID 不存在")
		}

		// 加入關聯（many2many）
		if err := tx.Model(data).Association("Users").Replace(&users); err != nil {
			return err
		}

		return nil
	})

	return data, err
}

func (s *TodoListDetailsService) Edit(db *gorm.DB, id int, name string, detail string) (*models.TodoListDetails, error) {
	updated := &models.TodoListDetails{}
	err := db.Transaction(func(tx *gorm.DB) error {
		// 取的原本的資料
		item, err := s.repo.FindByID(s.ctx, tx, id)
		if err != nil {
			return err
		}

		//更新內容
		item.Name = name
		item.Detail = detail
		if err := s.repo.Update(s.ctx, tx, item); err != nil {
			return err
		}
		*updated = *item
		return nil
	})

	return updated, err
}

func (s *TodoListDetailsService) Delete(db *gorm.DB, id int) (*models.TodoListDetails, error) {
	var deleted *models.TodoListDetails

	err := db.Transaction(func(tx *gorm.DB) error {
		// 先取的資料
		item, err := s.repo.FindByID(s.ctx, tx, id)
		if err != nil {
			return err
		}

		if err := s.repo.SoftDelete(s.ctx, tx, item); err != nil {
			return err
		}

		deleted = item
		return nil
	})

	return deleted, err
}
