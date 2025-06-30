package base

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// BaseRepository 為泛型資料存取層，提供基本 CRUD 操作，適用於任意模型 T。
type BaseRepository[T any] struct{}

func NewBaseRepository[T any]() *BaseRepository[T] {
	return &BaseRepository[T]{}
}

// Create 新增一筆資料到資料庫。
func (r *BaseRepository[T]) Create(ctx context.Context, db *gorm.DB, entity T) error {
	return db.WithContext(ctx).Create(entity).Error
}

// FindByID 根據主鍵 ID 查詢單一資料。
// 回傳 *T 型別資料，如果找不到或發生錯誤會回傳 error。
type FindOptions struct {
	PreloadFields  []string
	PreloadSelects map[string][]string
	Debug          bool
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, db *gorm.DB, id int, opts ...*FindOptions) (T, error) {
	var zero T // T 的零值

	if id == 0 {
		return zero, fmt.Errorf("FindByID requires a valid non-zero ID")
	}

	var model T
	query := db.WithContext(ctx)

	if len(opts) > 0 && opts[0] != nil {
		opt := opts[0]

		if opt.Debug {
			query = query.Debug()
		}

		for _, field := range opt.PreloadFields {
			if selects, ok := opt.PreloadSelects[field]; ok {
				query = query.Preload(field, func(tx *gorm.DB) *gorm.DB {
					return tx.Select(selects)
				})
			} else {
				query = query.Preload(field)
			}
		}
	}

	if err := query.First(&model, id).Error; err != nil {
		return zero, err
	}

	return model, nil
}

// Update 更新傳入的 entity 資料（只更新非零值欄位）。
// 通常用於已經查詢過的實體做修改後再儲存。
func (r *BaseRepository[T]) Update(ctx context.Context, db *gorm.DB, entity T) error {
	return db.WithContext(ctx).Model(entity).Updates(entity).Error
}

// UpdateByID 根據 ID 更新指定欄位（使用 map 格式傳入欲更新的欄位與值）。
// 範例: map[string]interface{}{"name": "新名稱"}
func (r *BaseRepository[T]) UpdateByID(ctx context.Context, db *gorm.DB, id int, updates map[string]interface{}) error {
	if id == 0 {
		return fmt.Errorf("UpdateByID requires a valid non-zero ID")
	}

	var model T // 建立空模型以指定操作對象的型別
	return db.WithContext(ctx).
		Model(&model).
		Where("id = ?", id).
		Updates(updates).
		Error
}

// SoftDelete 執行軟刪除（需要傳入 entity 實體）。
// 使用 GORM 的 Delete 方法，搭配模型的 DeletedAt 欄位進行軟刪。
func (r *BaseRepository[T]) SoftDelete(ctx context.Context, db *gorm.DB, entity T) error {
	return db.WithContext(ctx).Delete(entity).Error
}

// SoftDeleteByID 根據 ID 執行軟刪除。
// 不需先查詢實體，直接透過 ID 進行刪除操作。
func (r *BaseRepository[T]) SoftDeleteByID(ctx context.Context, db *gorm.DB, id int) error {
	if id == 0 {
		return fmt.Errorf("UpdateByID requires a valid non-zero ID")
	}

	var model T
	return db.WithContext(ctx).
		Model(&model).
		Where("id = ?", id).
		Delete(&model).
		Error
}

func (r *BaseRepository[T]) FindAllWithQuery(
	ctx context.Context,
	db *gorm.DB, // 這裡傳入已帶好條件的 *gorm.DB
	page int,
	pageSize int,
	orderBy ...string,
) ([]T, int64, error) {
	var results []T
	var model T
	var total int64

	query := db.WithContext(ctx).Model(&model)

	for _, o := range orderBy {
		query = query.Order(o)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
