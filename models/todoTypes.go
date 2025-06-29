package models

import (
	"todolist/models/base"
)

type TodoTypes struct {
	ID   int    `gorm:"primary_key" json:"id"`
	Name string `gorm:"type:varchar(255);NOT NULL;uniqueIndex" json:"name" binding:"required"`

	base.TimeModel
	base.OperatorModel
}

func (TodoTypes) TableName() string {
	return "to_do_types"
}
