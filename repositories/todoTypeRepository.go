package repositories

import (
	"context"
	"todolist/models"
	"todolist/repositories/base"

	"gorm.io/gorm"
)

type TodoTypeRepository struct {
	*base.BaseRepository[*models.TodoTypes] // 指標型態欄位
}

func NewTodoTypeRepository() *TodoTypeRepository {
	return &TodoTypeRepository{
		BaseRepository: base.NewBaseRepository[*models.TodoTypes](), // 回傳 *BaseRepository[T]
	}
}

func (r *TodoTypeRepository) IsNameExist(ctx context.Context, db *gorm.DB, name string, excludeID int) (bool, error) {
	var count int64
	query := db.WithContext(ctx).Model(&models.TodoTypes{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
