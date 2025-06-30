package seed

import (
	"todolist/models"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "Admin", Description: "管理員"},
		{Name: "Guest", Description: "訪客"},
		{Name: "Member", Description: "一般會員"},
	}

	for _, role := range roles {
		var existing models.Role
		// 如果資料已存在，就跳過（避免重複插入）
		if err := db.Where("name = ?", role.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			db.Create(&role)
		}
	}
}
