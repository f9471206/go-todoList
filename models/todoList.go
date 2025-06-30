package models

import (
	"todolist/models/base"
)

type TodoList struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	TypeID int    `gorm:"column:type_id;not null" json:"type_id"`
	Name   string `gorm:"type:varchar(255);not null" json:"name"`

	Type    TodoTypes         `gorm:"foreignKey:TypeID;constraint:OnDelete:CASCADE;" json:"type"`
	Details []TodoListDetails `gorm:"foreignKey:TodoListID;references:ID" json:"details"`

	base.TimeModel
	base.OperatorModel
}

func (TodoList) TableName() string {
	return "to_do_list"
}
