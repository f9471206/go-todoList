package interfaces

import (
	"context"
	"todolist/repositories/base"

	"todolist/models"

	"gorm.io/gorm"
)

type TodoListRepository interface {
	FindAllWithQuery(ctx context.Context, db *gorm.DB, page, pageSize int, orderBy ...string) ([]*models.TodoList, int64, error)
	FindByID(ctx context.Context, db *gorm.DB, id int, opts ...*base.FindOptions) (*models.TodoList, error)
	Create(ctx context.Context, db *gorm.DB, entity *models.TodoList) error
	Update(ctx context.Context, db *gorm.DB, entity *models.TodoList) error
	SoftDelete(ctx context.Context, db *gorm.DB, entity *models.TodoList) error
	IsNameExist(ctx context.Context, db *gorm.DB, name string, excludeID int) (bool, error)
}
