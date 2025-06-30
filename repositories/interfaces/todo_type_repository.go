package interfaces

import (
	"context"
	"todolist/repositories/base"

	"todolist/models"

	"gorm.io/gorm"
)

type TodoTypeRepository interface {
	FindAllWithQuery(ctx context.Context, db *gorm.DB, page, pageSize int, orderBy ...string) ([]*models.TodoTypes, int64, error)
	FindByID(ctx context.Context, db *gorm.DB, id int, opts ...*base.FindOptions) (*models.TodoTypes, error)
	Create(ctx context.Context, db *gorm.DB, entity *models.TodoTypes) error
	Update(ctx context.Context, db *gorm.DB, entity *models.TodoTypes) error
	SoftDelete(ctx context.Context, db *gorm.DB, entity *models.TodoTypes) error
	IsNameExist(ctx context.Context, db *gorm.DB, name string, excludeID int) (bool, error)
}
