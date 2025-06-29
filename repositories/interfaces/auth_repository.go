package interfaces

import (
	"context"
	"todolist/repositories/base"

	"todolist/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindAllWithQuery(ctx context.Context, db *gorm.DB, page, pageSize int, orderBy ...string) ([]*models.User, int64, error)
	FindByID(ctx context.Context, db *gorm.DB, id int, opts ...*base.FindOptions) (*models.User, error)
	Create(ctx context.Context, db *gorm.DB, entity *models.User) error
	Update(ctx context.Context, db *gorm.DB, entity *models.User) error
	SoftDelete(ctx context.Context, db *gorm.DB, entity *models.User) error
}
