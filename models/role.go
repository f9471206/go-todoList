package models

type Role struct {
	ID          uint
	Name        string
	Description string
	Users       []User `gorm:"many2many:user_roles"`
}
