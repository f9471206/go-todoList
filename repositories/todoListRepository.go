package repositories

import (
	"context"
	"todolist/models"
	"todolist/repositories/base"

	"gorm.io/gorm"
)

type TodoListRepository struct {
	*base.BaseRepository[*models.TodoList]
}

func NewTodoListRepository() *TodoListRepository {
	return &TodoListRepository{
		BaseRepository: base.NewBaseRepository[*models.TodoList](), // 回傳 *BaseRepository[T]
	}
}

// 檢查名稱是否存在，排除指定 ID（可為 0 代表不排除）
func (r *TodoListRepository) IsNameExist(ctx context.Context, db *gorm.DB, name string, excludeID int) (bool, error) {
	var count int64
	query := db.WithContext(ctx).Model(&models.TodoList{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
