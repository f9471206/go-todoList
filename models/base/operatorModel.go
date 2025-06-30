package base

import (
	"fmt"
	"todolist/utils"

	"gorm.io/gorm"
)

type OperatorModel struct {
	CreatedBy *uint `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *uint `gorm:"column:updated_by" json:"updated_by"`
	DeletedBy *uint `gorm:"column:deleted_by" json:"deleted_by"`
}

func (m *OperatorModel) BeforeCreate(tx *gorm.DB) (err error) {
	val := tx.Statement.Context.Value(utils.UserIDKey)
	if valFloat, ok := val.(float64); ok {
		userID := uint(valFloat)
		m.CreatedBy = &userID
	}
	return
}

func (m *OperatorModel) BeforeUpdate(tx *gorm.DB) (err error) {
	val := tx.Statement.Context.Value(utils.UserIDKey)
	if valFloat, ok := val.(float64); ok {
		userID := uint(valFloat)
		m.UpdatedBy = &userID
	}
	return
}

func (m *OperatorModel) BeforeDelete(tx *gorm.DB) (err error) {
	val := tx.Statement.Context.Value(utils.UserIDKey)
	if valFloat, ok := val.(float64); ok {
		userID := uint(valFloat)

		if tx.Statement.Schema != nil {
			pkField := tx.Statement.Schema.PrioritizedPrimaryField
			if pkField != nil {
				pkValue, _ := pkField.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue)
				return tx.Model(tx.Statement.Model).
					Where(fmt.Sprintf("%s = ?", pkField.DBName), pkValue).
					Update("deleted_by", userID).Error
			}
		}
	}
	return nil
}
