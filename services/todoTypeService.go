package services

import (
	"context"
	"errors"
	"todolist/models"
	"todolist/repositories/interfaces"
	"todolist/utils"

	"gorm.io/gorm"
)

type TodoTypeService struct {
	ctx  context.Context
	repo interfaces.TodoTypeRepository
}

func NewTodoTypeService(ctx context.Context, repo interfaces.TodoTypeRepository) *TodoTypeService {
	return &TodoTypeService{
		ctx:  ctx,
		repo: repo,
	}
}

func (s *TodoTypeService) Create(db *gorm.DB, name string) (*models.TodoTypes, error) {
	result := &models.TodoTypes{Name: name}

	err := db.Transaction(func(tx *gorm.DB) error {
		// 檢查名稱是否重複
		exist, err := s.repo.IsNameExist(s.ctx, tx, name, 0)
		if err != nil {
			return err
		}
		if exist {
			return errors.New("名稱已存在")
		}
		// 建立
		return s.repo.Create(s.ctx, tx, result)
	})

	return result, err
}

func (s *TodoTypeService) Index(db *gorm.DB, keyword string, page, pageSize int, orderBy []string) (*utils.PaginatedResult[*models.TodoTypes], error) {
	query := db.Model(&models.TodoTypes{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	list, total, err := s.repo.FindAllWithQuery(s.ctx, query, page, pageSize, orderBy...)
	if err != nil {
		return nil, err
	}

	return utils.NewPaginatedResult(list, total, page, pageSize), nil
}

func (s *TodoTypeService) Show(db *gorm.DB, id int) (*models.TodoTypes, error) {
	return s.repo.FindByID(s.ctx, db, id)
}

func (s *TodoTypeService) Edit(db *gorm.DB, id int, name string) (*models.TodoTypes, error) {
	updated := &models.TodoTypes{}

	err := db.Transaction(func(tx *gorm.DB) error {
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
		if err := s.repo.Update(s.ctx, tx, item); err != nil {
			return err
		}

		*updated = *item
		return nil
	})

	return updated, err
}

func (s *TodoTypeService) Delete(db *gorm.DB, id int) (*models.TodoTypes, error) {
	var deleted *models.TodoTypes

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
