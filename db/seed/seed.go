package seed

import "gorm.io/gorm"

func Seed(db *gorm.DB) {
	SeedRoles(db)
	SeedUsers(db)
}
