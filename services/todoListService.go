package services

import (
	"context"
	"errors"
	"todolist/models"
	"todolist/repositories/base"
	"todolist/repositories/interfaces"
	"todolist/utils"

	"gorm.io/gorm"
)

type TodoListService struct {
	ctx  context.Context
	repo interfaces.TodoListRepository
}

func NewTodoListService(ctx context.Context, repo interfaces.TodoListRepository) *TodoListService {
	return &TodoListService{
		ctx:  ctx,
		repo: repo,
	}
}

func (s *TodoListService) Create(db *gorm.DB, name string, typeID int) (*models.TodoList, error) {
	result := &models.TodoList{
		Name:   name,
		TypeID: typeID,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// 可以先檢查 type_id 是否存在，避免外鍵錯誤
		var count int64
		if err := tx.Model(&models.TodoTypes{}).Where("id = ?", typeID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("type_id 不存在")
		}

		exist, err := s.repo.IsNameExist(s.ctx, tx, name, 0)
		if err != nil {
			return err
		}
		if exist {
			return errors.New("名稱已存在")
		}

		return s.repo.Create(s.ctx, tx, result)
	})

	return result, err
}

func (s *TodoListService) Index(db *gorm.DB, keyword string, page, pageSize int, orderBy []string) (*utils.PaginatedResult[*models.TodoList], error) {
	query := db.Model(&models.TodoList{})
	list, total, err := s.repo.FindAllWithQuery(s.ctx, query, page, pageSize, orderBy...)
	if err != nil {
		return nil, err
	}

	return utils.NewPaginatedResult(list, total, page, pageSize), nil
}

func (s *TodoListService) Show(db *gorm.DB, id int) (*models.TodoList, error) {
	opts := &base.FindOptions{
		Debug:         true,
		PreloadFields: []string{"Details", "Type", "Details.Users"},
		PreloadSelects: map[string][]string{
			"Details.Users": {"id", "account"},
		},
	}
	return s.repo.FindByID(s.ctx, db, id, opts)
}

func (s *TodoListService) Edit(db *gorm.DB, id int, name string, typeID int) (*models.TodoList, error) {
	updated := &models.TodoList{}

	err := db.Transaction(func(tx *gorm.DB) error {
		// 可以先檢查 type_id 是否存在，避免外鍵錯誤
		var count int64
		if err := tx.Model(&models.TodoTypes{}).Where("id = ?", typeID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("type_id 不存在")
		}

		// 名稱是否已存在（排除自己）
		exist, err := s.repo.IsNameExist(s.ctx, tx, name, id)
		if err != nil {
			return err
		}
		if exist {
			return errors.New("名稱已存在")
		}

		// 取得原本資料
		item, err := s.repo.FindByID(s.ctx, tx, id)
		if err != nil {
			return err
		}

		// 更新內容
		item.Name = name
		item.TypeID = typeID
		if err := s.repo.Update(s.ctx, tx, item); err != nil {
			return err
		}

		*updated = *item
		return nil
	})

	return updated, err
}

func (s *TodoListService) Delete(db *gorm.DB, id int) (*models.TodoList, error) {
	var deleted *models.TodoList

	err := db.Transaction(func(tx *gorm.DB) error {
		// 先取得資料
		item, err := s.repo.FindByID(s.ctx, tx, id)
		if err != nil {
			return err
		}

		// 軟刪除
		if err := s.repo.SoftDelete(s.ctx, tx, item); err != nil {
			return err
		}

		deleted = item
		return nil
	})

	return deleted, err
}
