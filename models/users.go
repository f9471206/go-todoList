package models

import "todolist/models/base"

type User struct {
	ID              int               `gorm:"primary_key" json:"id"`
	Account         string            `gorm:"type:varchar(255);NOT NULL;uniqueIndex" json:"account" binding:"required"`
	Password        string            `json:"-"`
	Roles           []Role            `gorm:"many2many:user_roles" json:"roles,omitempty"`
	TodoListDetails []TodoListDetails `gorm:"many2many:to_do_task_assignments" json:"todo_list_details,omitempty"`

	base.TimeModel
	base.OperatorModel
}

func (User) TableName() string {
	return "users"
}
