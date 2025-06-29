package interfaces

import (
	"context"
	"todolist/repositories/base"

	"todolist/models"

	"gorm.io/gorm"
)

type TodoListDetailsRepository interface {
	FindAllWithQuery(ctx context.Context, db *gorm.DB, page, pageSize int, orderBy ...string) ([]*models.TodoListDetails, int64, error)
	FindByID(ctx context.Context, db *gorm.DB, id int, opts ...*base.FindOptions) (*models.TodoListDetails, error)
	Create(ctx context.Context, db *gorm.DB, entity *models.TodoListDetails) error
	Update(ctx context.Context, db *gorm.DB, entity *models.TodoListDetails) error
	SoftDelete(ctx context.Context, db *gorm.DB, entity *models.TodoListDetails) error
}
