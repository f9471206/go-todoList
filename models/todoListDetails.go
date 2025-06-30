package models

import (
	"todolist/models/base"
)

type TodoListDetails struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	TodoListID int    `gorm:"column:to_do_list_id;not null" json:"to_do_list_id"`
	Name       string `gorm:"type:varchar(255);not null" json:"name"`
	Detail     string `gorm:"type:varchar(255);not null" json:"detail"`

	Users []User `gorm:"many2many:to_do_task_assignments;joinForeignKey:ToDoListDetailID;joinReferences:UserID"`

	base.TimeModel
	base.OperatorModel
}

func (TodoListDetails) TableName() string {
	return "to_do_list_details"
}
